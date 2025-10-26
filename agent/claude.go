package agent

import (
	"context"
	"example-agentic-workflow-engine-go/domain"
	"log"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

type claudeAgent struct {
	name   string
	client anthropic.Client
}

func NewClaudeAgent() domain.Agent {
	return &claudeAgent{
		name:   "Claude",
		client: anthropic.NewClient(),
	}
}

func (c *claudeAgent) Run(input domain.Request) <-chan domain.ResponseChunk {
	out := make(chan domain.ResponseChunk)

	go func() {
		log.Printf("%s agent started\n", c.name)
		begin := time.Now()

		defer func() {
			close(out)

			elapsed := time.Since(begin)
			log.Printf("%s agent ended (elapsed: %v)\n", c.name, elapsed)
		}()

		stream := c.client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
			MaxTokens: 1024,
			Model:     anthropic.ModelClaudeSonnet4_0,
			Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(input.Message))},
		})

		for stream.Next() {
			event := stream.Current()

			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					if deltaVariant.Text == "" {
						continue
					}
					out <- domain.ResponseChunk{
						Name:    c.name,
						Content: deltaVariant.Text,
					}
				}

			}
		}
		err := stream.Err()
		if err != nil {
			log.Printf("%s agent error: %v\n", c.name, err)
		}
	}()

	return out
}
