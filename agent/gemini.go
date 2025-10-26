package agent

import (
	"context"
	"example-agentic-workflow-engine-go/domain"
	"log"
	"time"

	"google.golang.org/genai"
)

type geminiAgent struct {
	name   string
	client *genai.Client
}

func NewGeminiAgent() domain.Agent {
	client, _ := genai.NewClient(context.TODO(), &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})

	return &geminiAgent{
		name:   "Gemini",
		client: client,
	}
}

func (g *geminiAgent) Run(input domain.Request) <-chan domain.ResponseChunk {
	out := make(chan domain.ResponseChunk)

	go func() {
		log.Printf("%s agent started\n", g.name)
		begin := time.Now()

		defer func() {
			close(out)

			elapsed := time.Since(begin)
			log.Printf("%s agent ended (elapsed: %v)\n", g.name, elapsed)
		}()

		if g.client == nil {
			log.Printf("%s agent error: Gemini client is nil\n", g.name)
			return
		}

		iter := g.client.Models.GenerateContentStream(
			context.TODO(),
			"gemini-2.5-flash-lite",
			genai.Text(input.Message),
			nil)

		for chunk, err := range iter {
			if err != nil {
				log.Printf("%s agent error: %v\n", g.name, err)
				return
			}

			if len(chunk.Candidates) == 0 || len(chunk.Candidates[0].Content.Parts) == 0 {
				continue
			}

			contentChunk := chunk.Candidates[0].Content.Parts[0]
			if contentChunk.Text == "" {
				continue
			}

			out <- domain.ResponseChunk{
				Name:    g.name,
				Content: contentChunk.Text,
			}
		}
	}()

	return out
}
