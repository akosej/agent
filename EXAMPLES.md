# Scripts de ejemplo

## Ejemplo de uso básico con Python (para pruebas)

```python
import subprocess
import json

# Iniciar el agente
agent = subprocess.Popen(['./agent.exe'], 
                        stdin=subprocess.PIPE, 
                        stdout=subprocess.PIPE,
                        stderr=subprocess.PIPE,
                        text=True)

# Enviar mensaje
mensaje = "Hola, ¿cómo estás?"
agent.stdin.write(mensaje + "\n")
agent.stdin.flush()

# Leer respuesta
respuesta = agent.stdout.readline()
print(f"Respuesta: {respuesta}")
```

## Script de PowerShell para pruebas automatizadas

```powershell
# test-agent.ps1

# Compilar
go build -o agent.exe ./cmd/agent

# Ejecutar y enviar comandos
$commands = @(
    "Hola",
    "¿Cuál es la capital de España?",
    "/stats",
    "/exit"
)

foreach ($cmd in $commands) {
    Write-Host "> $cmd"
    Start-Sleep -Seconds 2
}
```

## Integración con API REST (próximamente)

```bash
# Enviar mensaje por HTTP
curl -X POST http://localhost:8080/api/message \
  -H "Content-Type: application/json" \
  -d '{"text": "Hola, ¿cómo estás?"}'
```
