# 🤖 AgentIA - Resumen Ejecutivo

## ¿Qué es AgentIA?

AgentIA es un **agente conversacional inteligente** desarrollado en Go que:
- 💬 Mantiene conversaciones naturales en español
- 🧠 Aprende de cada interacción
- 🎤 Entiende comandos de voz
- 📊 Mejora continuamente con el uso

## 📦 Estructura Creada

```
d:\DevOps\tools\agent/
│
├── 📂 cmd/agent/              → Aplicación principal
│   ├── main.go               → Punto de entrada
│   └── version.go            → Información de versión
│
├── 📂 internal/               → Código interno
│   ├── agent/                → Coordinador principal
│   ├── speech/               → Reconocimiento de voz
│   ├── nlp/                  → Procesamiento de lenguaje
│   └── learning/             → Sistema de aprendizaje
│
├── 📂 pkg/                    → Paquetes compartidos
│   ├── storage/              → Base de datos SQLite
│   └── logger/               → Sistema de logging
│
├── 📂 configs/                → Configuración
│   └── config.yaml           → Configuración principal
│
├── 📂 data/                   → Datos generados
│   ├── conversations/        → Conversaciones guardadas
│   └── models/               → Modelos aprendidos
│
├── 📂 logs/                   → Archivos de log
│
├── 📄 .env.example           → Plantilla de variables
├── 📄 go.mod                 → Dependencias de Go
├── 📄 README.md              → Documentación completa
├── 📄 QUICKSTART.md          → Inicio rápido
├── 📄 ARCHITECTURE.md        → Diseño del sistema
├── 📄 CHANGELOG.md           → Historia de cambios
├── 📄 setup.ps1              → Script de configuración
└── 📄 Makefile               → Tareas automatizadas
```

## 🚀 Inicio Rápido (3 pasos)

### 1️⃣ Configurar API Key
```powershell
Copy-Item .env.example .env
notepad .env  # Agregar tu OPENAI_API_KEY
```

### 2️⃣ Instalar y Compilar
```powershell
.\setup.ps1
```

### 3️⃣ Ejecutar
```powershell
.\agent.exe
```

## 💡 Características Principales

| Módulo | Descripción | Estado |
|--------|-------------|--------|
| 🗣️ **NLP** | Procesamiento con OpenAI GPT | ✅ Completo |
| 🎤 **Voz** | Reconocimiento de voz con Whisper | ✅ Completo |
| 🧠 **Aprendizaje** | Sistema de patrones y feedback | ✅ Completo |
| 💾 **Storage** | Persistencia en SQLite | ✅ Completo |
| 📊 **Stats** | Métricas y estadísticas | ✅ Completo |
| 🌐 **API REST** | Servidor HTTP | 🔜 Próximamente |
| 🎨 **Web UI** | Interface web | 🔜 Próximamente |

## 📝 Archivos Importantes

### Configuración
- **`.env`** → API keys y secretos
- **`configs/config.yaml`** → Configuración del sistema

### Documentación
- **`README.md`** → Documentación completa
- **`QUICKSTART.md`** → Guía de inicio rápido
- **`ARCHITECTURE.md`** → Diseño del sistema
- **`EXAMPLES.md`** → Ejemplos de uso
- **`DEVELOPMENT.md`** → Notas de desarrollo

### Scripts
- **`setup.ps1`** → Configuración automática (PowerShell)
- **`setup.bat`** → Configuración automática (CMD)
- **`Makefile`** → Tareas de desarrollo

## 🎯 Comandos del Agente

```
> Hola                  → Conversar normalmente
> /help                → Ver ayuda
> /stats               → Ver estadísticas
> /export              → Exportar conocimiento
> /exit                → Salir
```

## 🔧 Tecnologías Usadas

```go
// Core
Go 1.21+

// IA
OpenAI GPT-3.5/4      → Conversación
OpenAI Whisper        → Transcripción de voz

// Storage
SQLite                → Base de datos

// Audio
PortAudio             → Captura de micrófono

// Config
YAML                  → Archivos de configuración
dotenv                → Variables de entorno
```

## 📊 Arquitectura Simplificada

```
Usuario
  ↓
┌─────────────────┐
│  Entrada        │ → Texto o Voz
│  (Speech/Text)  │
└────────┬────────┘
         ↓
┌─────────────────┐
│  Procesamiento  │ → NLP + Intent Detection
│  (Agent Core)   │
└────────┬────────┘
         ↓
┌─────────────────┐
│  Aprendizaje    │ → Patterns + Feedback
│  (Learning)     │
└────────┬────────┘
         ↓
┌─────────────────┐
│  Storage        │ → SQLite + Logs
└─────────────────┘
```

## 🧪 Testing

```powershell
# Instalar dependencias
go mod download

# Compilar
go build -o agent.exe .\cmd\agent

# Ejecutar tests
go test ./...

# Ver coverage
go test -cover ./...

# Formatear código
go fmt ./...
```

## 📈 Próximos Pasos

1. ✅ **Completado** - Estructura base del proyecto
2. ⏳ **Siguiente** - Descargar dependencias: `go mod download`
3. ⏳ **Después** - Configurar API key en `.env`
4. ⏳ **Luego** - Compilar: `go build -o agent.exe .\cmd\agent`
5. ⏳ **Finalmente** - Ejecutar: `.\agent.exe`

## 🐛 Solución Rápida de Problemas

| Problema | Solución |
|----------|----------|
| Error API key | Verificar `.env` tiene `OPENAI_API_KEY=sk-...` |
| Error compilación | Ejecutar `go mod download` |
| Error PortAudio | Desactivar speech: `EnableSpeech: false` |
| No responde | Verificar conexión a internet y API key válida |

## 📚 Recursos

- **OpenAI API**: https://platform.openai.com/docs
- **Go Docs**: https://golang.org/doc/
- **PortAudio**: http://www.portaudio.com/

## 📞 Soporte

1. Ver `README.md` para documentación completa
2. Ver `QUICKSTART.md` para inicio rápido
3. Ver logs en `logs/agent.log` para debugging
4. Revisar `ARCHITECTURE.md` para entender el diseño

---

**Versión:** 1.0.0  
**Fecha:** 11 de Diciembre, 2025  
**Licencia:** MIT  
**Lenguaje:** Go 1.21+

🎉 **¡Proyecto completamente estructurado y listo para usar!**
