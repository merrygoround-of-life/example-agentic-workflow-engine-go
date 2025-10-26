package agent

import (
	"context"
	"example-agentic-workflow-engine-go/domain"
	"log"
	"time"

	"github.com/openai/openai-go/v3"
)

type chatgptAgent struct {
	name   string
	client openai.Client
}

func NewChatGptAgent() domain.Agent {
	return &chatgptAgent{
		name:   "ChatGPT",
		client: openai.NewClient(),
	}
}

func (c *chatgptAgent) Run(input domain.Request) <-chan domain.ResponseChunk {
	out := make(chan domain.ResponseChunk)

	go func() {
		log.Printf("%s agent started\n", c.name)
		begin := time.Now()

		defer func() {
			close(out)

			elapsed := time.Since(begin)
			log.Printf("%s agent ended (elapsed: %v)\n", c.name, elapsed)
		}()

		stream := c.client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(input.Message)},
			Model:    openai.ChatModelGPT4_1Mini,
		})

		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) == 0 {
				continue
			}

			contentChunk := chunk.Choices[0].Delta.Content
			if contentChunk == "" {
				continue
			}

			out <- domain.ResponseChunk{
				Name:    c.name,
				Content: contentChunk,
			}
		}
		err := stream.Err()
		if err != nil {
			log.Printf("%s agent error: %v\n", c.name, err)
		}
	}()

	return out
}
