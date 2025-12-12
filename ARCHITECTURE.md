# Arquitectura del Sistema AgentIA

## Vista General

```
┌─────────────────────────────────────────────────────────────┐
│                         Usuario                             │
│                    (Texto o Voz)                            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    CAPA DE ENTRADA                          │
│  ┌──────────────────┐         ┌──────────────────┐         │
│  │ Speech Recognizer│         │   Text Input     │         │
│  │  (PortAudio)    │         │   (Terminal)     │         │
│  └────────┬─────────┘         └────────┬─────────┘         │
│           │                            │                    │
│           ▼                            │                    │
│  ┌──────────────────┐                 │                    │
│  │  Transcriber     │                 │                    │
│  │  (Whisper API)   │                 │                    │
│  └────────┬─────────┘                 │                    │
│           └────────────────┬───────────┘                    │
└────────────────────────────┼────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  CAPA DE PROCESAMIENTO                      │
│                                                             │
│  ┌──────────────────────────────────────────────────┐      │
│  │           Agent Core (Orquestador)               │      │
│  │                                                  │      │
│  │  ┌──────────────┐    ┌──────────────┐          │      │
│  │  │ NLP Processor│◄───┤   Context    │          │      │
│  │  │  (OpenAI)    │    │   Manager    │          │      │
│  │  └──────┬───────┘    └──────────────┘          │      │
│  │         │                                        │      │
│  │         ▼                                        │      │
│  │  ┌──────────────┐    ┌──────────────┐          │      │
│  │  │    Intent    │───►│   Response   │          │      │
│  │  │  Detection   │    │  Generation  │          │      │
│  │  └──────────────┘    └──────┬───────┘          │      │
│  │                              │                  │      │
│  └──────────────────────────────┼──────────────────┘      │
│                                 │                          │
└─────────────────────────────────┼──────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   CAPA DE APRENDIZAJE                       │
│                                                             │
│  ┌──────────────────────────────────────────────────┐      │
│  │          Learning Engine                         │      │
│  │                                                  │      │
│  │  ┌──────────────┐    ┌──────────────┐          │      │
│  │  │  Interaction │    │   Pattern    │          │      │
│  │  │   Recorder   │───►│  Recognition │          │      │
│  │  └──────────────┘    └──────────────┘          │      │
│  │         │                     │                 │      │
│  │         ▼                     ▼                 │      │
│  │  ┌──────────────┐    ┌──────────────┐          │      │
│  │  │   Feedback   │    │ Knowledge    │          │      │
│  │  │   System     │    │     Base     │          │      │
│  │  └──────────────┘    └──────┬───────┘          │      │
│  │                              │                  │      │
│  └──────────────────────────────┼──────────────────┘      │
│                                 │                          │
└─────────────────────────────────┼──────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  CAPA DE PERSISTENCIA                       │
│                                                             │
│  ┌──────────────────┐         ┌──────────────────┐         │
│  │     Storage      │         │      Logger      │         │
│  │    (SQLite)      │         │    (File/Console)│         │
│  │                  │         │                  │         │
│  │ • Interactions   │         │ • Debug logs     │         │
│  │ • Patterns       │         │ • Info logs      │         │
│  │ • Stats          │         │ • Error logs     │         │
│  │ • Conversations  │         │ • Interactions   │         │
│  └──────────────────┘         └──────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## Flujo de Datos Principal

### 1. Entrada de Usuario
```
Usuario → [Texto/Voz] → Speech Recognizer → Transcriber → Texto Normalizado
```

### 2. Procesamiento
```
Texto → Agent Core
         ↓
      NLP Processor
         ↓
   Intent Detection → [Buscar en Knowledge Base]
         ↓                     ↓
   ¿Pattern Found?          [Yes] → Usar Pattern
         ↓                           ↓
       [No]                   [No] → Generate Response (OpenAI)
         ↓                           ↓
   Generate Response (OpenAI)       │
         └───────────────────────────┘
                     ↓
              Response Final
```

### 3. Aprendizaje
```
Interaction → Learning Engine
                ↓
        Record Interaction
                ↓
        Extract Patterns
                ↓
        Update Knowledge Base
                ↓
        Save to Storage
```

## Componentes Detallados

### Speech Module (`internal/speech/`)
- **Recognizer**: Captura audio del micrófono
- **Transcriber**: Convierte audio a texto usando Whisper API
- **Formato**: WAV, 16kHz, Mono

### NLP Module (`internal/nlp/`)
- **Processor**: Interfaz con OpenAI GPT
- **Intent Detection**: Clasifica la intención del usuario
- **Response Generation**: Genera respuestas contextuales
- **Summarization**: Resume conversaciones largas

### Learning Module (`internal/learning/`)
- **Engine**: Motor principal de aprendizaje
- **Interaction**: Modelo de interacción usuario-agente
- **Pattern**: Patrones aprendidos con confianza
- **Feedback**: Sistema de retroalimentación
- **Stats**: Estadísticas de uso

### Agent Core (`internal/agent/`)
- **Orchestration**: Coordina todos los módulos
- **Context Management**: Mantiene contexto de conversación
- **History**: Historial de mensajes (últimos 20)
- **Lifecycle**: Inicio, ejecución, cierre

### Storage Layer (`pkg/storage/`)
- **Database**: SQLite para persistencia
- **Tables**: interactions, patterns, stats
- **Backup**: Sistema de backup configurable
- **Queries**: CRUD operations optimizadas

### Logger (`pkg/logger/`)
- **Levels**: DEBUG, INFO, WARN, ERROR
- **Outputs**: Archivo y consola simultáneos
- **Format**: Timestamp, nivel, archivo, mensaje
- **Rotation**: Tamaño máximo y backups

## Configuración por Capas

```yaml
configs/config.yaml:
  ├── agent: Configuración general
  ├── speech: Parámetros de voz
  ├── nlp: Configuración de IA
  ├── learning: Parámetros de aprendizaje
  ├── storage: Base de datos
  └── logging: Sistema de logs

.env:
  ├── OPENAI_API_KEY: Clave de OpenAI
  ├── DEEPGRAM_API_KEY: Clave opcional
  └── DB_PATH: Ruta de base de datos
```

## Patrones de Diseño Utilizados

1. **Singleton**: Logger, Storage
2. **Factory**: Creación de componentes
3. **Strategy**: Diferentes providers de voz/NLP
4. **Observer**: Sistema de feedback y aprendizaje
5. **Repository**: Capa de persistencia

## Flujo de Errores

```
Error Occurred
    ↓
Logger.Error()
    ↓
[Is Critical?] → Yes → Close Gracefully
    ↓
   No
    ↓
Return Error to User
    ↓
Continue Execution
```

## Dependencias Externas

```
OpenAI API
    ├── GPT-3.5/4: Generación de respuestas
    └── Whisper: Transcripción de voz

PortAudio
    └── Captura de audio del sistema

SQLite
    └── Almacenamiento local
```

## Escalabilidad Futura

### Horizontal
- API REST para múltiples clientes
- Load balancing
- Cache distribuido

### Vertical
- Vector database para similitud semántica
- Fine-tuning de modelos
- Procesamiento paralelo de requests

### Features
- Múltiples idiomas
- Personalidades configurables
- Integración con servicios externos
- Sistema de plugins
