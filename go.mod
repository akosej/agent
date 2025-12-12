module github.com/akosej/agent

go 1.21

require (
	github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.18
	gopkg.in/yaml.v3 v3.0.1
)

// Dependencias removidas (ya no necesitamos OpenAI):
// github.com/sashabaranov/go-openai - Reemplazado por llamadas HTTP directas a Ollama
