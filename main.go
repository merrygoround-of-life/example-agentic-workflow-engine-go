package main

import (
	"example-agentic-workflow-engine-go/agent"
	"example-agentic-workflow-engine-go/domain"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// POST /echo 엔드포인트
	r.POST("/echo", func(c *gin.Context) {
		var req domain.Request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		a := agent.NewMockAgent()
		out := a.Run(req) // 입력 메시지를 처리하는 고루틴 + 채널 생성
		sse(c, out)       // 채널을 SSE 응답으로 스트리밍
		return
	})

	agents := []domain.Agent{
		agent.NewChatGptAgent(),
		agent.NewClaudeAgent(),
		agent.NewGeminiAgent(),
	}

	// POST /orchestrate 엔드포인트
	r.POST("/orchestrate", func(c *gin.Context) {
		var req domain.Request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chs := spawnWorkloads(req, agents) // 병렬 실행
		merged := fanIn(chs)               // 채널 병합
		out := orchestrate(merged)         // 순서 조율

		sse(c, out)
		return
	})

	r.Run()
}

func sse(c *gin.Context, ch <-chan domain.ResponseChunk) bool {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}

	return c.Stream(func(w io.Writer) bool {
		if chunk, ok := <-ch; ok {
			// 채널에서 읽은 데이터를 SSE 이벤트로 전송
			c.SSEvent("", chunk)
			return true
		} else {
			// 채널이 닫히면 스트리밍 종료
			return false
		}
	})
}

func spawnWorkloads(input domain.Request, agents []domain.Agent) []<-chan domain.ResponseChunk {
	var chs []<-chan domain.ResponseChunk

	for _, a := range agents {
		ch := a.Run(input) // 각 에이전트를 비동기로 시작
		chs = append(chs, ch)
	}

	return chs
}

func fanIn(chs []<-chan domain.ResponseChunk) <-chan domain.FanInChunk {
	var wg sync.WaitGroup
	wg.Add(len(chs))

	out := make(chan domain.FanInChunk)
	for i, v := range chs {
		// 각 채널을 동시에 읽기
		go func(order int, ch <-chan domain.ResponseChunk) {
			for chunk := range ch {
				// 읽은 데이터를 에이전트 인덱스와 함께 출력 채널로 전달
				out <- domain.FanInChunk{Order: order, Chunk: chunk, Done: false}
			}

			// 다 읽었으면 완료 신호 전달
			out <- domain.FanInChunk{Order: order, Done: true}
			wg.Done()
		}(i, v)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func orchestrate(in <-chan domain.FanInChunk) <-chan domain.ResponseChunk {
	log.Printf("orchestrator started\n")
	begin := time.Now()

	out := make(chan domain.ResponseChunk)

	go func() {
		defer func() {
			close(out)
			elapsed := time.Since(begin)
			log.Printf("orchestrator ended (elapsed: %v)\n", elapsed)
		}()

		buf := make([][]domain.ResponseChunk, 0) // 각 에이전트의 청크를 저장할 버퍼
		cur := 0                                 // 현재 출력 대상 에이전트 인덱스

		for chunk := range in {
			// 0. 버퍼 확장
			for len(buf) <= chunk.Order {
				buf = append(buf, []domain.ResponseChunk{})
			}

			// 1. 현재 순서가 아니면: 버퍼에 저장
			if chunk.Order != cur {
				if !chunk.Done {
					buf[chunk.Order] = append(buf[chunk.Order], chunk.Chunk)
				}

				continue
			}

			// 2. 현재 순서: 버퍼 flush 후 처리
			for _, buffered := range buf[chunk.Order] {
				out <- buffered
			}
			buf[chunk.Order] = nil

			if chunk.Done {
				cur++
			} else {
				out <- chunk.Chunk
			}
		}

		// 3. 남은 버퍼 flush
		for cur < len(buf) {
			for _, buffered := range buf[cur] {
				out <- buffered
			}
			cur++
		}
	}()

	return out
}
