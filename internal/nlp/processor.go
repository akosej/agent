package nlp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Intent representa una intención detectada
type Intent struct {
	Name       string
	Confidence float64
	Entities   map[string]string
}

// Config contiene la configuración del procesador NLP
type Config struct {
	Model       string
	MaxTokens   int
	Temperature float32
	OllamaURL   string // URL del servidor Ollama local
}

// Processor maneja el procesamiento de lenguaje natural
type Processor struct {
	client    *http.Client
	config    Config
	ollamaURL string
}

// Message representa un mensaje en la conversación
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaRequest representa una solicitud a Ollama
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  Options   `json:"options,omitempty"`
}

// Options representa opciones para la generación
type Options struct {
	Temperature float32 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"` // max_tokens en Ollama
}

// OllamaResponse representa la respuesta de Ollama
type OllamaResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

// NewProcessor crea una nueva instancia del procesador NLP
func NewProcessor(ollamaURL string, config Config) *Processor {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434" // Puerto por defecto de Ollama
	}
	return &Processor{
		client:    &http.Client{},
		config:    config,
		ollamaURL: ollamaURL,
	}
}

// ProcessText procesa texto y genera una respuesta
func (p *Processor) ProcessText(ctx context.Context, text string, conversationHistory []Message) (string, error) {
	messages := append(conversationHistory, Message{
		Role:    "user",
		Content: text,
	})

	return p.callOllama(ctx, messages)
}

// callOllama realiza una llamada al servidor Ollama local
func (p *Processor) callOllama(ctx context.Context, messages []Message) (string, error) {
	request := OllamaRequest{
		Model:    p.config.Model,
		Messages: messages,
		Stream:   false,
		Options: Options{
			Temperature: p.config.Temperature,
			NumPredict:  p.config.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error codificando request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.ollamaURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error llamando a Ollama: %w (Asegúrate de que Ollama esté corriendo con: ollama serve)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama respondió con error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// DetectIntent detecta la intención del usuario en el texto
func (p *Processor) DetectIntent(ctx context.Context, text string) (*Intent, error) {
	systemPrompt := `Eres un asistente que detecta intenciones. 
Analiza el siguiente texto y responde SOLO con el formato:
INTENCIÓN: [nombre_intención]
CONFIANZA: [0.0-1.0]
ENTIDADES: [clave1=valor1, clave2=valor2]

Intenciones posibles: saludo, despedida, pregunta, comando, conversacion, ayuda`

	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: text,
		},
	}

	// Usar temperatura más baja para detección de intenciones
	oldTemp := p.config.Temperature
	p.config.Temperature = 0.3
	oldTokens := p.config.MaxTokens
	p.config.MaxTokens = 150

	response, err := p.callOllama(ctx, messages)

	// Restaurar configuración
	p.config.Temperature = oldTemp
	p.config.MaxTokens = oldTokens

	if err != nil {
		return nil, fmt.Errorf("error detectando intención: %w", err)
	}

	return p.parseIntent(response), nil
}

// parseIntent parsea la respuesta del modelo en un Intent
func (p *Processor) parseIntent(response string) *Intent {
	intent := &Intent{
		Name:       "conversacion",
		Confidence: 0.5,
		Entities:   make(map[string]string),
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "INTENCIÓN:") {
			intent.Name = strings.TrimSpace(strings.TrimPrefix(line, "INTENCIÓN:"))
		} else if strings.HasPrefix(line, "CONFIANZA:") {
			var conf float64
			fmt.Sscanf(line, "CONFIANZA: %f", &conf)
			intent.Confidence = conf
		} else if strings.HasPrefix(line, "ENTIDADES:") {
			entitiesStr := strings.TrimSpace(strings.TrimPrefix(line, "ENTIDADES:"))
			entities := strings.Split(entitiesStr, ",")
			for _, entity := range entities {
				parts := strings.Split(strings.TrimSpace(entity), "=")
				if len(parts) == 2 {
					intent.Entities[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return intent
}

// GenerateResponse genera una respuesta basada en el contexto
func (p *Processor) GenerateResponse(ctx context.Context, userInput string, context map[string]interface{}) (string, error) {
	systemPrompt := fmt.Sprintf(`Eres un asistente virtual inteligente llamado AgentIA. 
Respondes en español de manera amigable y útil.
Contexto actual: %v`, context)

	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userInput,
		},
	}

	response, err := p.callOllama(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("error generando respuesta: %w", err)
	}

	return response, nil
}

// SummarizeConversation resume una conversación
func (p *Processor) SummarizeConversation(ctx context.Context, messages []Message) (string, error) {
	var conversationText strings.Builder
	for _, msg := range messages {
		conversationText.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	summaryMessages := []Message{
		{
			Role:    "system",
			Content: "Resume la siguiente conversación en 2-3 oraciones, capturando los puntos principales.",
		},
		{
			Role:    "user",
			Content: conversationText.String(),
		},
	}

	// Guardar y modificar temporalmente la configuración
	oldTemp := p.config.Temperature
	p.config.Temperature = 0.5
	oldTokens := p.config.MaxTokens
	p.config.MaxTokens = 200

	response, err := p.callOllama(ctx, summaryMessages)

	// Restaurar configuración
	p.config.Temperature = oldTemp
	p.config.MaxTokens = oldTokens

	if err != nil {
		return "", fmt.Errorf("error resumiendo conversación: %w", err)
	}

	return response, nil
}
