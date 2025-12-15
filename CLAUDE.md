# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
```bash
go run main.go
```
Server runs on `localhost:8080` by default.

### Environment Setup
Required API keys as environment variables:
```bash
export OPENAI_API_KEY='your-openai-api-key-here'
export ANTHROPIC_API_KEY='your-anthropic-api-key-here'
export GOOGLE_API_KEY='your-google-api-key-here'
```

### Testing the API
```bash
# Test the orchestrated multi-agent workflow
curl -N -H 'Content-Type: application/json' \
  -d '{"message": "What is the meaning of life?"}' \
  'localhost:8080/orchestrate'

# Test the echo agent
curl -N -H 'Content-Type: application/json' \
  -d '{"message": "Hello World"}' \
  'localhost:8080/echo'
```

## Architecture Overview

This is an agentic workflow engine that orchestrates multiple AI agents (ChatGPT, Claude, Gemini) to process requests in parallel and stream responses sequentially via Server-Sent Events (SSE).

### Core Components

**main.go**: HTTP server with two endpoints:
- `/echo` - Single mock agent for testing
- `/orchestrate` - Multi-agent parallel execution with ordered streaming

**domain/model.go**: Core data structures:
- `Request` - Input message structure
- `ResponseChunk` - Streaming response unit with agent name and content
- `FanInChunk` - Internal orchestration structure with ordering
- `Agent` interface - Contract for all agent implementations

**agent/ directory**: Agent implementations:
- `chatgpt.go` - OpenAI ChatGPT integration
- `claude.go` - Anthropic Claude integration  
- `gemini.go` - Google Gemini integration
- `echo.go` - Mock agent for testing

### Workflow Architecture

1. **Parallel Execution**: `spawnWorkloads()` starts all agents concurrently
2. **Fan-In**: `fanIn()` merges multiple agent channels into one ordered stream
3. **Orchestration**: `orchestrate()` ensures sequential output (Agent 0 → Agent 1 → Agent 2) by buffering out-of-order chunks
4. **SSE Streaming**: `sse()` converts the ordered stream to Server-Sent Events

### Key Dependencies
- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/openai/openai-go/v3` - OpenAI client
- `github.com/anthropics/anthropic-sdk-go` - Anthropic client
- `google.golang.org/genai` - Google Gemini client

### Module Configuration
- Go version: 1.25.3
- Module name: `example-agentic-workflow-engine-go`

## Container Deployment

### Docker Build and Run
```bash
# Build the Docker image
docker build -t agentic-workflow-engine .

# Run the container with environment variables
docker run -p 8080:8080 \
  -e OPENAI_API_KEY='your-openai-api-key-here' \
  -e ANTHROPIC_API_KEY='your-anthropic-api-key-here' \
  -e GOOGLE_API_KEY='your-google-api-key-here' \
  agentic-workflow-engine
```

### Docker Compose
```bash
# Set environment variables in .env file or export them
export OPENAI_API_KEY='your-openai-api-key-here'
export ANTHROPIC_API_KEY='your-anthropic-api-key-here'
export GOOGLE_API_KEY='your-google-api-key-here'

# Start the application
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# Stop the application
docker-compose down
```