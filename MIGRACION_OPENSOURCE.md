# 🔄 Migración a Open Source - AgentIA v2.0

## Resumen de Cambios

Este proyecto ha sido migrado de usar servicios pagos de OpenAI a soluciones **100% open source** que funcionan completamente en local.

## ✨ Principales Cambios

### 1. Procesamiento de Lenguaje Natural (NLP)
- **ANTES**: OpenAI GPT-3.5/4 (requiere API key y conexión a internet)
- **AHORA**: Ollama con modelos locales (Llama, Mistral, Phi3, etc.)
- **Archivo modificado**: `internal/nlp/processor.go`

### 2. Transcripción de Voz
- **ANTES**: OpenAI Whisper API (servicio en la nube)
- **AHORA**: Whisper.cpp (ejecutable local)
- **Archivo modificado**: `internal/speech/transcriber.go`

### 3. Configuración
- **ANTES**: Requería `OPENAI_API_KEY` en `.env`
- **AHORA**: Solo necesita URLs locales de servicios
- **Archivos modificados**: 
  - `configs/config.yaml`
  - `cmd/agent/main.go`
  - `internal/agent/agent.go`

### 4. Dependencias
- **ANTES**: `github.com/sashabaranov/go-openai`
- **AHORA**: Cliente HTTP estándar de Go
- **Archivo modificado**: `go.mod`

## 🎯 Ventajas de la Migración

| Aspecto | Antes | Ahora |
|---------|-------|-------|
| **Costo** | ~$0.002 por request | $0 (gratis) |
| **Privacidad** | Datos enviados a OpenAI | 100% local |
| **Internet** | Requerido | No necesario |
| **Límites** | Rate limits de API | Ilimitado |
| **Velocidad** | Depende de internet | Inmediato en local |
| **Personalización** | Limitada | Total control |

## 📦 Nuevos Requisitos

### Software Necesario

1. **Ollama** (requerido para NLP)
   - Descarga: https://ollama.com/
   - Modelos recomendados:
     - `llama3.2:1b` - Ultrarrápido (1GB)
     - `llama3.2:3b` - Balanceado (3GB) ⭐
     - `mistral:7b` - Alta calidad (4GB)

2. **Whisper.cpp** (opcional, solo para voz)
   - Repositorio: https://github.com/ggerganov/whisper.cpp
   - Modelo recomendado: `ggml-base.bin`

### Hardware Mínimo

- **RAM**: 4GB (8GB recomendado)
- **Disco**: 5GB libres (para modelos)
- **CPU**: Cualquier CPU de 64 bits
- **GPU**: Opcional (acelera significativamente)

## 🔧 Cambios Técnicos Detallados

### 1. internal/nlp/processor.go

#### Estructuras Nuevas:
```go
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OllamaRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
    Options  Options   `json:"options,omitempty"`
}
```

#### Método Principal:
```go
func (p *Processor) callOllama(ctx context.Context, messages []Message) (string, error)
```

### 2. internal/speech/transcriber.go

#### Métodos Nuevos:
```go
func NewTranscriber(whisperPath, modelPath, language string) *Transcriber
func NewTranscriberWithAPI(apiURL, language string) *Transcriber
func (t *Transcriber) transcribeWithWhisperCpp(ctx context.Context, audioPath string) (string, error)
```

### 3. configs/config.yaml

#### Nuevos Campos:
```yaml
nlp:
  ollama_url: "http://localhost:11434"
  
speech:
  whisper_path: "./whisper.cpp/main"
  model_path: "./models/ggml-base.bin"
  api_url: "http://localhost:8000"
```

### 4. internal/agent/agent.go

#### Config Actualizado:
```go
type Config struct {
    OllamaURL     string  // Nuevo
    WhisperPath   string  // Nuevo
    WhisperModel  string  // Nuevo
    WhisperAPIURL string  // Nuevo
    // OpenAIKey removido
    // DeepgramKey removido
}
```

## 🚀 Pasos para Migrar una Instalación Existente

### 1. Instalar Ollama
```powershell
# Windows: Descargar desde ollama.com
# Luego:
ollama pull llama3.2:3b
```

### 2. Actualizar el Código
```powershell
git pull origin main
go mod download
go build -o agent.exe ./cmd/agent
```

### 3. Actualizar Configuración
Edita `configs/config.yaml`:
```yaml
nlp:
  model: "llama3.2:3b"
  ollama_url: "http://localhost:11434"
```

### 4. Iniciar Ollama
```powershell
ollama serve
```

### 5. Ejecutar el Agente
```powershell
.\agent.exe
```

## 🔍 Comparación de Funcionalidad

| Función | Antes | Ahora | Status |
|---------|-------|-------|--------|
| Chat conversacional | ✅ GPT-3.5 | ✅ Llama/Mistral | ✅ Mantenido |
| Detección de intenciones | ✅ GPT | ✅ Modelo local | ✅ Mantenido |
| Historial de conversación | ✅ | ✅ | ✅ Mantenido |
| Sistema de aprendizaje | ✅ | ✅ | ✅ Mantenido |
| Transcripción de voz | ✅ Whisper API | ✅ Whisper.cpp | ✅ Mantenido |
| Almacenamiento SQLite | ✅ | ✅ | ✅ Mantenido |
| Estadísticas | ✅ | ✅ | ✅ Mantenido |
| Exportar conocimiento | ✅ | ✅ | ✅ Mantenido |

## 📊 Rendimiento Esperado

### Tiempos de Respuesta (en hardware moderno)

| Modelo | Primera respuesta | Respuestas siguientes | RAM |
|--------|-------------------|----------------------|-----|
| llama3.2:1b | ~2s | ~0.5s | 1GB |
| llama3.2:3b | ~3s | ~1s | 3GB |
| mistral:7b | ~5s | ~2s | 4GB |

### Calidad de Respuestas

- **llama3.2:1b**: Adecuado para conversaciones simples
- **llama3.2:3b**: Excelente para uso general ⭐
- **mistral:7b**: Respuestas de alta calidad

## ⚠️ Limitaciones Conocidas

1. **Primera ejecución lenta**: El modelo se carga en memoria (una sola vez)
2. **Modelos grandes requieren RAM**: Verifica requisitos antes de descargar
3. **Sin GPU**: Las respuestas serán más lentas (pero funcionales)

## 🆘 Soporte

Si tienes problemas con la migración:

1. Lee `INSTALACION_RAPIDA.md`
2. Revisa la sección "Solución de Problemas" en `README.md`
3. Verifica que Ollama esté corriendo: `ollama serve`
4. Confirma que el modelo esté descargado: `ollama list`

## 📝 Notas para Desarrolladores

### Testing con Diferentes Modelos

```powershell
# Descargar modelo
ollama pull phi3:mini

# Actualizar config.yaml
# nlp:
#   model: "phi3:mini"

# Reiniciar el agente
```

### Modo Debug

Edita `configs/config.yaml`:
```yaml
logging:
  level: "debug"
```

### Monitorear Ollama

```powershell
# Ver logs de Ollama
ollama logs

# Ver modelos instalados
ollama list

# Ver uso de recursos
ollama ps
```

## 🎉 Conclusión

La migración a open source hace que AgentIA sea:
- **Más accesible**: Sin costos de API
- **Más privado**: Datos nunca salen de tu PC
- **Más flexible**: Cambia modelos según necesites
- **Más confiable**: No depende de servicios externos

¡Disfruta de tu agente IA completamente autónomo! 🚀
