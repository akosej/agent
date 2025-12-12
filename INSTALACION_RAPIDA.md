# 🚀 Instalación Rápida - AgentIA Open Source

Esta guía te llevará de 0 a tener el agente funcionando en menos de 10 minutos.

## Paso 1: Instalar Ollama (3 minutos)

### Windows
1. Ve a https://ollama.com/download/windows
2. Descarga e instala el ejecutable
3. Abre PowerShell y ejecuta:
```powershell
ollama pull llama3.2:3b
```

### Linux
```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3.2:3b
```

### macOS
```bash
brew install ollama
ollama pull llama3.2:3b
```

⏱️ **Tiempo**: ~3 minutos (descarga de ~2GB)

## Paso 2: Instalar el Agente (2 minutos)

```powershell
# Clonar el repositorio
git clone https://github.com/yourusername/agent
cd agent

# Instalar dependencias
go mod download

# Compilar
go build -o agent.exe ./cmd/agent
```

## Paso 3: Configurar (1 minuto)

Edita `configs/config.yaml` (opcional, la configuración por defecto funciona):

```yaml
nlp:
  model: "llama3.2:3b"
  ollama_url: "http://localhost:11434"
```

## Paso 4: ¡Ejecutar! (5 segundos)

```powershell
# Asegúrate de que Ollama está corriendo
ollama serve

# En otra terminal, ejecuta el agente
.\agent.exe
```

## ✅ ¡Listo!

Ahora puedes chatear con tu agente:

```
> Hola, ¿cómo estás?
AgentIA: ¡Hola! Estoy muy bien...

> ¿Cuál es la capital de España?
AgentIA: La capital de España es Madrid...
```

## 🔧 Comandos Útiles

- `/help` - Ver ayuda
- `/stats` - Ver estadísticas
- `/export` - Exportar conocimiento aprendido
- `/exit` - Salir

## 🎯 Próximos Pasos (Opcional)

### Usar un modelo más rápido (si va lento):
```powershell
ollama pull llama3.2:1b
# Actualiza model en config.yaml a "llama3.2:1b"
```

### Usar un modelo de mejor calidad:
```powershell
ollama pull mistral:7b
# Actualiza model en config.yaml a "mistral:7b"
```

### Agregar reconocimiento de voz:
Ve a la sección "Instalación de Whisper.cpp" en el README.md

## ❓ Problemas Comunes

**Error: "connection refused"**
```powershell
# Inicia Ollama
ollama serve
```

**Error: "modelo no encontrado"**
```powershell
# Descarga el modelo
ollama pull llama3.2:3b
```

**El agente va lento**
```powershell
# Usa un modelo más pequeño
ollama pull llama3.2:1b
```

## 📚 Más Información

- [README completo](README.md)
- [Arquitectura del proyecto](ARCHITECTURE.md)
- [Ejemplos de uso](EXAMPLES.md)

---

**¡Disfruta de tu agente IA completamente privado y offline!** 🎉
