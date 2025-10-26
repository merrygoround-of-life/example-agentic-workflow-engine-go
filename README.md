# example-agentic-workflow-engine-go

https://devocean.sk.com/blog/techBoardDetail.do?ID=167973


## Quickstart

1. Setup API key

   You should export your api-key as an environment variable as shown below.
   ```
   $ export OPENAI_API_KEY='your-openai-api-key-here'
   $ export ANTHROPIC_API_KEY='your-anthropic-api-key-here'
   $ export GOOGLE_API_KEY='your-google-api-key-here'
   ```

2. Run server

   ```
   $ go run main.go
   ```

3. Test with curl
   ```
   $ curl -N -H 'Content-Type: application/json' \
   -d '{"message": "What is the meaning of life?"}' \
   'localhost:8080/orchestrate'
   data:{"name":"ChatGPT","content":"The"}			# ChatGPT Agent Response

   data:{"name":"ChatGPT","content":" meaning"}

   data:{"name":"ChatGPT","content":" of"}

   data:{"name":"ChatGPT","content":" life"}

   data:{"name":"ChatGPT","content":" is"}

   data:{"name":"ChatGPT","content":" a"}

   data:{"name":"ChatGPT","content":" deeply"}

   data:{"name":"ChatGPT","content":" personal"}

   data:{"name":"ChatGPT","content":" and"}

   data:{"name":"ChatGPT","content":" philosophical"}

   data:{"name":"ChatGPT","content":" question"}

   ...

   data:{"name":"Claude","content":"This"}			# Claude Agent Response

   data:{"name":"Claude","content":" is one of humanity"}

   data:{"name":"Claude","content":"'s oldest and most profound questions, an"}

   data:{"name":"Claude","content":"d there's no single answer"}

   ...

   data:{"name":"Gemini","content":"Ah"}			# Gemini Agent Response

   data:{"name":"Gemini","content":", the age-old question! \"What is the meaning of life?\" It"}

   data:{"name":"Gemini","content":"'s a question that has been pondered by philosophers, theologians, scientists, artists, and everyday people for millennia. And the truth is, there's no single, universally agreed-upon answer.\n\nHere's a breakdown of why it's"}

   ...

   $ 
   ```
