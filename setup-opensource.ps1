# Script de configuración para AgentIA Open Source
# No necesitas API keys - Todo funciona localmente!

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "AgentIA - Configuración Open Source" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Verificar directorio
if (-not (Test-Path "go.mod")) {
    Write-Host "Error: Ejecuta este script desde la raíz del proyecto" -ForegroundColor Red
    exit 1
}

# Crear directorios
Write-Host "[1/6] Creando directorios..." -ForegroundColor Yellow
$dirs = @("logs", "data", "data/conversations", "models")
foreach ($dir in $dirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}
Write-Host "   ✓ Directorios creados" -ForegroundColor Green

# Verificar Ollama
Write-Host "`n[2/6] Verificando Ollama..." -ForegroundColor Yellow
$ollamaCmd = Get-Command ollama -ErrorAction SilentlyContinue
if ($ollamaCmd) {
    Write-Host "   ✓ Ollama instalado" -ForegroundColor Green
    
    # Verificar modelos
    Write-Host "   Modelos disponibles:" -ForegroundColor Cyan
    ollama list
    
    # Ofrecer descargar modelo si no hay ninguno
    $choice = Read-Host "`n   ¿Descargar llama3.2:3b si no lo tienes? (s/n)"
    if ($choice -eq "s") {
        Write-Host "   Descargando modelo (esto puede tardar unos minutos)..." -ForegroundColor Yellow
        ollama pull llama3.2:3b
        Write-Host "   ✓ Modelo descargado" -ForegroundColor Green
    }
} else {
    Write-Host "   ✗ Ollama NO está instalado" -ForegroundColor Red
    Write-Host "   Instálalo desde: https://ollama.com/download/windows" -ForegroundColor Yellow
    Write-Host "   Después ejecuta: ollama pull llama3.2:3b" -ForegroundColor Yellow
    $continue = Read-Host "`n   ¿Continuar de todos modos? (s/n)"
    if ($continue -ne "s") {
        exit 1
    }
}

# Descargar dependencias Go
Write-Host "`n[3/6] Descargando dependencias de Go..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✓ Dependencias descargadas" -ForegroundColor Green
} else {
    Write-Host "   ✗ Error descargando dependencias" -ForegroundColor Red
    exit 1
}

# Verificar configuración
Write-Host "`n[4/6] Verificando configuración..." -ForegroundColor Yellow
if (Test-Path "configs/config.yaml") {
    Write-Host "   ✓ Archivo de configuración encontrado" -ForegroundColor Green
    
    # Mostrar modelo configurado
    $config = Get-Content "configs/config.yaml" -Raw
    if ($config -match 'model:\s*"(.+?)"') {
        Write-Host "   Modelo configurado: $($Matches[1])" -ForegroundColor Cyan
    }
} else {
    Write-Host "   ✗ Archivo de configuración no encontrado" -ForegroundColor Red
}

# Compilar
Write-Host "`n[5/6] Compilando el agente..." -ForegroundColor Yellow
go build -o agent.exe ./cmd/agent
if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✓ Compilación exitosa" -ForegroundColor Green
} else {
    Write-Host "   ✗ Error en la compilación" -ForegroundColor Red
    exit 1
}

# Verificar Ollama está corriendo
Write-Host "`n[6/6] Verificando servicio de Ollama..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:11434/api/tags" -Method GET -TimeoutSec 2 -ErrorAction Stop
    Write-Host "   ✓ Ollama está corriendo" -ForegroundColor Green
} catch {
    Write-Host "   ⚠ Ollama no está corriendo" -ForegroundColor Yellow
    Write-Host "   Ejecuta [ollama serve] en otra terminal" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✓ Configuración completada!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para ejecutar el agente:" -ForegroundColor Yellow
Write-Host "  1. Asegúrate de que Ollama esté corriendo: ollama serve" -ForegroundColor White
Write-Host "  2. Ejecuta el agente: .\agent.exe" -ForegroundColor White
Write-Host ""
Write-Host "Comandos útiles:" -ForegroundColor Yellow
Write-Host "  /help   - Mostrar ayuda" -ForegroundColor White
Write-Host "  /stats  - Ver estadísticas" -ForegroundColor White
Write-Host "  /export - Exportar conocimiento" -ForegroundColor White
Write-Host "  /exit   - Salir" -ForegroundColor White
Write-Host ""
Write-Host "🎉 ¡Disfruta de tu agente IA local y privado!" -ForegroundColor Green
Write-Host ""
