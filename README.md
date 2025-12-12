# AgentIA - Agente Inteligente con Aprendizaje y Reconocimiento de Voz

Un agente conversacional inteligente desarrollado en Go que aprende de las interacciones y cuenta con capacidades de reconocimiento de voz. **100% Open Source y funciona completamente offline** - no necesitas API keys ni conexión a internet.

## 🚀 Características

- **Procesamiento de Lenguaje Natural (NLP)**: Usa **Ollama** para ejecutar modelos como Llama, Mistral, Phi3 localmente
- **Reconocimiento de Voz**: Transcribe audio usando **Whisper.cpp** (ejecutado localmente)
- **Sistema de Aprendizaje**: Mejora continuamente con cada interacción
- **Detección de Intenciones**: Identifica qué quiere hacer el usuario
- **Almacenamiento Persistente**: Guarda conversaciones y patrones en SQLite
- **Estadísticas**: Rastrea métricas de rendimiento y feedback
- **Exportación de Conocimiento**: Permite guardar y cargar el conocimiento aprendido
- **🔒 100% Privado**: Todos los datos y procesamiento permanecen en tu máquina
- **🌐 Funciona Offline**: No requiere conexión a internet después de la instalación

## 📋 Requisitos Previos

- Go 1.21 o superior
- **Ollama** (para el modelo de lenguaje local)
- **Whisper.cpp** (opcional, solo para reconocimiento de voz)
- PortAudio (opcional, solo para captura de audio en vivo)

### Instalación de Ollama (REQUERIDO)

Ollama permite ejecutar modelos de IA localmente sin necesidad de API keys.

**Windows:**
1. Descarga el instalador desde: https://ollama.com/download/windows
2. Ejecuta el instalador
3. Abre PowerShell y descarga un modelo:
```powershell
ollama pull llama3.2:3b
```

**Linux:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3.2:3b
```

**macOS:**
```bash
brew install ollama
ollama pull llama3.2:3b
```

**Modelos recomendados:**
- `llama3.2:1b` - Ultrarrápido (1GB RAM) - ideal para equipos limitados
- `llama3.2:3b` - Balanceado (3GB RAM) - **RECOMENDADO**
- `mistral:7b` - Alta calidad (4GB RAM)
- `phi3:mini` - Eficiente (2GB RAM)
- `qwen2.5:3b` - Multilingüe (3GB RAM)

### Instalación de Whisper.cpp (OPCIONAL - solo para voz)

Solo necesario si quieres usar reconocimiento de voz.

**Windows:**
```powershell
# Clonar whisper.cpp
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp

# Compilar (necesitas Visual Studio o MinGW)
cmake -B build
cmake --build build --config Release

# Descargar modelo (base es suficiente para español)
cd models
.\download-ggml-model.ps1 base
cd ..\..
```

**Linux/macOS:**
```bash
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
make
bash ./models/download-ggml-model.sh base
cd ..
```

### Instalación de PortAudio (OPCIONAL - solo para captura de voz en vivo)

**Windows:**
```powershell
pacman -S mingw-w64-x86_64-portaudio  # Con MSYS2
```

**Linux:**
```bash
sudo apt-get install portaudio19-dev
```

**macOS:**
```bash
brew install portaudio
```

## 🛠️ Instalación del Agente

1. **Clonar el repositorio:**
```powershell
git clone https://github.com/yourusername/agent
cd agent
```

2. **Instalar dependencias Go:**
```powershell
go mod download
```

3. **Configurar rutas en config.yaml:**

Edita `configs/config.yaml` y ajusta las rutas según tu instalación:

```yaml
nlp:
  model: "llama3.2:3b"  # El modelo que descargaste con Ollama
  ollama_url: "http://localhost:11434"  # URL por defecto de Ollama

speech:
  provider: "whisper-cpp"  # o "whisper-api" si usas un servidor
  whisper_path: "./whisper.cpp/main"  # Ruta a tu instalación de whisper.cpp
  model_path: "./whisper.cpp/models/ggml-base.bin"
```

4. **Compilar el proyecto:**
```powershell
go build -o agent.exe ./cmd/agent
```

5. **Iniciar Ollama (si no está corriendo):**
```powershell
ollama serve
```

## 🎯 Uso

### Modo Texto (Recomendado para empezar)

Ejecutar el agente en modo texto:
```powershell
.\agent.exe
```

Luego simplemente escribe tus mensajes:
```
> Hola, ¿cómo estás?
AgentIA: ¡Hola! Estoy muy bien, gracias por preguntar. ¿En qué puedo ayudarte hoy?

> ¿Cuál es la capital de Francia?
AgentIA: La capital de Francia es París...
```

### Comandos Disponibles

Dentro del agente, puedes usar estos comandos:

- `/help` o `/ayuda` - Muestra la ayuda
- `/stats` - Muestra estadísticas del agente
- `/export` - Exporta el conocimiento aprendido a un archivo JSON
- `/exit`, `/salir` o `/quit` - Cierra el agente

### Modo Voz

Para habilitar el reconocimiento de voz, editar `cmd/agent/main.go`:
```go
EnableSpeech: true,  // Cambiar de false a true
```

Luego recompilar:
```powershell
go build -o agent.exe ./cmd/agent
```

## 📁 Estructura del Proyecto

```
agent/
├── cmd/
│   └── agent/
│       └── main.go              # Punto de entrada principal
├── internal/
│   ├── agent/
│   │   └── agent.go             # Coordinador principal
│   ├── speech/
│   │   ├── recognizer.go        # Captura de audio
│   │   └── transcriber.go       # Transcripción a texto
│   ├── nlp/
│   │   └── processor.go         # Procesamiento NLP
│   └── learning/
│       └── engine.go            # Motor de aprendizaje
├── pkg/
│   ├── storage/
│   │   └── storage.go           # Almacenamiento SQLite
│   └── logger/
│       └── logger.go            # Sistema de logging
├── configs/
│   └── config.yaml              # Configuración principal
├── data/
│   ├── agent.db                 # Base de datos (se crea automáticamente)
│   └── conversations/           # Conversaciones guardadas
├── logs/
│   └── agent.log                # Logs de la aplicación
├── .env                         # Variables de entorno (crear desde .env.example)
├── .env.example                 # Plantilla de variables de entorno
├── go.mod                       # Dependencias de Go
└── README.md                    # Este archivo
```

## ⚙️ Configuración

Editar `configs/config.yaml` para personalizar el comportamiento:

```yaml
agent:
  name: "AgentIA"
  version: "2.0.0"
  language: "es"

nlp:
  model: "llama3.2:3b"  # Cualquier modelo de Ollama
  max_tokens: 500
  temperature: 0.7
  ollama_url: "http://localhost:11434"

speech:
  provider: "whisper-cpp"
  whisper_path: "./whisper.cpp/main"
  model_path: "./whisper.cpp/models/ggml-base.bin"

learning:
  enabled: true
  learning_rate: 0.01
  confidence_threshold: 0.7
```

### Cambiar Modelo de IA

Para usar un modelo diferente:
```powershell
# Descargar otro modelo
ollama pull mistral:7b

# Actualizar config.yaml
# nlp:
#   model: "mistral:7b"
```

## 🧠 Sistema de Aprendizaje

El agente aprende de las interacciones de las siguientes formas:

1. **Patrones de Conversación**: Identifica intenciones comunes y respuestas exitosas
2. **Contexto**: Mantiene el historial de la conversación
3. **Feedback**: Mejora basándose en la retroalimentación (por implementar en UI)
4. **Estadísticas**: Rastrea métricas para mejorar continuamente

### Exportar e Importar Conocimiento

```
> /export
Conocimiento exportado a: data/knowledge_export_1234567890.json
```

Para importar conocimiento previo, agregar función en el código.

## 🔧 Desarrollo

### Ejecutar en modo desarrollo:
```powershell
go run ./cmd/agent
```

### Ejecutar tests:
```powershell
go test ./...
```

### Verificar código:
```powershell
go vet ./...
```

### Formatear código:
```powershell
go fmt ./...
```

## 📊 Logs y Depuración

Los logs se guardan en `logs/agent.log`. Para cambiar el nivel de logging, editar `configs/config.yaml`:

```yaml
logging:
  level: "debug"  # debug, info, warn, error
```

## 🐛 Solución de Problemas

### Error: "error llamando a Ollama" o "connection refused"
**Causa**: Ollama no está corriendo
**Solución**: 
```powershell
ollama serve
```
O reinicia el servicio de Ollama en Windows.

### Error: "modelo no encontrado"
**Causa**: No has descargado el modelo especificado en config.yaml
**Solución**:
```powershell
ollama pull llama3.2:3b  # O el modelo que especifiques
```

### Error con Whisper.cpp: "whisperPath no configurado"
**Causa**: No has configurado la ruta a whisper.cpp en config.yaml
**Solución**: 
1. Instala whisper.cpp (ver sección de instalación)
2. Actualiza `whisper_path` en config.yaml con la ruta correcta

### El agente responde muy lento
**Solución**: Usa un modelo más pequeño
```powershell
ollama pull llama3.2:1b  # Más rápido
# Actualiza config.yaml: model: "llama3.2:1b"
```

### Error de compilación con PortAudio
- Solo necesario para captura de voz en vivo (opcional)
- Verifica que PortAudio esté instalado correctamente
- En Windows, asegura que MSYS2 esté en el PATH

### Base de datos bloqueada
- Cierra todas las instancias del agente antes de iniciar una nueva

## 🎯 Ventajas de Esta Versión Open Source

✅ **Sin costos de API**: No pagas por cada consulta  
✅ **100% Privado**: Tus conversaciones nunca salen de tu PC  
✅ **Funciona offline**: No necesitas internet después de instalar  
✅ **Personalizable**: Cambia modelos según tus necesidades  
✅ **Rápido**: Respuestas instantáneas en hardware moderno  
✅ **Sin límites**: Usa el agente tanto como quieras  

## 🚧 Roadmap

- [x] Migración a Ollama (modelos locales)
- [x] Integración con Whisper.cpp local
- [ ] Interface web con React
- [ ] Soporte para múltiples modelos simultáneos
- [ ] API REST para integración externa
- [ ] Sistema de plugins
- [ ] Memoria vectorial con embeddings locales
- [ ] Soporte para RAG (Retrieval Augmented Generation)

## 📝 Licencia

MIT License - ver el archivo LICENSE para más detalles

## 🤝 Contribuciones

Las contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📧 Contacto

Para preguntas o sugerencias, abrir un issue en el repositorio.

---

Desarrollado con ❤️ usando Go
