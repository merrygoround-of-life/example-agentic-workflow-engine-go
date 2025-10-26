package agent

import (
	"example-agentic-workflow-engine-go/domain"
	"math/rand"
	"time"
)

type mockAgent struct {
	name string
}

func NewMockAgent() domain.Agent {
	return &mockAgent{
		name: "Echo",
	}
}

func (m *mockAgent) Run(input domain.Request) <-chan domain.ResponseChunk {
	// 출력 채널 생성
	out := make(chan domain.ResponseChunk)

	// 고루틴을 이용한 비동기 작업 수행
	go func() {
		defer close(out)

		// 입력 문자열의 각 문자마다 0~1s 사이의 랜덤 대기 후 출력 채널로 전송 (비동기 처리 시뮬레이션)
		for _, v := range input.Message {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			out <- domain.ResponseChunk{Name: m.name, Content: string(v)}
		}
	}()

	return out
}
