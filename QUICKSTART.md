# 🚀 Inicio Rápido - AgentIA

## Pasos para Empezar (5 minutos)

### 1. Prerrequisitos
Asegúrate de tener instalado:
- ✅ Go 1.21 o superior: https://golang.org/dl/
- ✅ Git (opcional)

### 2. Configurar API Key de OpenAI

**Paso 1:** Obtén tu API Key de OpenAI
- Ve a: https://platform.openai.com/api-keys
- Crea una nueva API key

**Paso 2:** Configura el archivo .env
```powershell
# Copia el archivo de ejemplo
Copy-Item .env.example .env

# Edita .env y reemplaza 'your_openai_api_key_here' con tu API key real
notepad .env
```

Tu archivo `.env` debe verse así:
```env
OPENAI_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxx
DB_PATH=./data/agent.db
LOG_LEVEL=info
SPEECH_LANGUAGE=es-ES
```

### 3. Ejecutar el Script de Configuración

**Opción A: PowerShell (Recomendado)**
```powershell
.\setup.ps1
```

**Opción B: Command Prompt**
```cmd
setup.bat
```

**Opción C: Manual**
```powershell
# Descargar dependencias
go mod download

# Compilar
go build -o agent.exe .\cmd\agent

# Ejecutar
.\agent.exe
```

### 4. Primer Uso

Una vez que el agente esté corriendo, verás:
```
╔════════════════════════════════════════════╗
║                                            ║
║        AgentIA v1.0.0                      ║
║        Agente Inteligente con              ║
║        Aprendizaje y Reconocimiento        ║
║        de Voz                              ║
║                                            ║
╚════════════════════════════════════════════╝

Escribe tu mensaje o usa /help para ver comandos.

>
```

**Prueba estos comandos:**
```
> Hola, ¿cómo estás?
> ¿Cuál es la capital de España?
> Cuéntame un chiste
> /stats
> /help
> /exit
```

## 🔧 Solución de Problemas Comunes

### Error: "OPENAI_API_KEY no está configurado"
**Solución:** Verifica que el archivo `.env` existe y contiene tu API key correcta.

### Error: "go: module not found"
**Solución:** Ejecuta:
```powershell
go mod download
go mod tidy
```

### Error al compilar con portaudio
**Solución temporal:** El reconocimiento de voz está deshabilitado por defecto. Para habilitarlo:
1. Instala PortAudio (ver README.md)
2. En `cmd/agent/main.go` cambia `EnableSpeech: false` a `true`

### El agente no responde
**Solución:** 
- Verifica tu conexión a internet
- Verifica que tu API key de OpenAI sea válida
- Revisa el archivo `logs/agent.log` para más detalles

## 📝 Comandos Disponibles

| Comando | Descripción |
|---------|-------------|
| `/help` | Muestra ayuda |
| `/stats` | Muestra estadísticas |
| `/export` | Exporta conocimiento aprendido |
| `/exit` | Salir del agente |

## 🎯 Próximos Pasos

1. **Experimentar:** Prueba diferentes tipos de conversaciones
2. **Ver Estadísticas:** Usa `/stats` para ver el aprendizaje
3. **Personalizar:** Edita `configs/config.yaml` para ajustar el comportamiento
4. **Explorar Logs:** Revisa `logs/agent.log` para debugging

## 📚 Más Información

- Ver `README.md` para documentación completa
- Ver `EXAMPLES.md` para ejemplos de uso
- Ver `DEVELOPMENT.md` para desarrollo

## 💡 Consejos

- **Para mejores respuestas:** Usa GPT-4 en vez de GPT-3.5 (edita `configs/config.yaml`)
- **Para conversaciones largas:** El agente mantiene historial de las últimas 20 interacciones
- **Para feedback:** Las estadísticas mejoran con el uso continuo

---

¿Problemas? Abre un issue o revisa la documentación completa en README.md
