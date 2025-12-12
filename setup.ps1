# Script de inicio rápido para PowerShell

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "AgentIA - Configuración Inicial" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Verificar si existe .env
if (-not (Test-Path .env)) {
    Write-Host "[1/5] Creando archivo .env desde plantilla..." -ForegroundColor Yellow
    Copy-Item .env.example .env
    Write-Host "IMPORTANTE: Edita el archivo .env y agrega tu OPENAI_API_KEY" -ForegroundColor Red
    Write-Host ""
    Read-Host "Presiona Enter para continuar después de configurar .env"
} else {
    Write-Host "[1/5] Archivo .env encontrado" -ForegroundColor Green
}

# Crear directorios necesarios
Write-Host "[2/5] Creando directorios..." -ForegroundColor Yellow
$dirs = @("logs", "data", "data\conversations", "data\models")
foreach ($dir in $dirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}
Write-Host "Directorios creados" -ForegroundColor Green

# Descargar dependencias
Write-Host "[3/5] Descargando dependencias de Go..." -ForegroundColor Yellow
go mod download
go get github.com/mattn/go-sqlite3

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Fallo al descargar dependencias" -ForegroundColor Red
    exit 1
}
Write-Host "Dependencias descargadas" -ForegroundColor Green

# Compilar el proyecto
Write-Host "[4/5] Compilando el proyecto..." -ForegroundColor Yellow
go build -o agent.exe .\cmd\agent

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: La compilación falló" -ForegroundColor Red
    Write-Host "Verifica que tengas Go instalado correctamente" -ForegroundColor Red
    Read-Host "Presiona Enter para salir"
    exit 1
}
Write-Host "Compilación exitosa!" -ForegroundColor Green

Write-Host "[5/5] Verificando configuración..." -ForegroundColor Yellow
$env_content = Get-Content .env -Raw
if ($env_content -match "your_openai_api_key_here") {
    Write-Host "ADVERTENCIA: Parece que no has configurado tu OPENAI_API_KEY" -ForegroundColor Red
    Write-Host "Edita el archivo .env antes de ejecutar el agente" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "AgentIA está listo para usar" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para ejecutar: .\agent.exe" -ForegroundColor Green
Write-Host "Para ayuda: Ver README.md" -ForegroundColor Green
Write-Host ""

# Preguntar si quiere ejecutar ahora
$respuesta = Read-Host "¿Desea ejecutar AgentIA ahora? (S/N)"
if ($respuesta -eq "S" -or $respuesta -eq "s") {
    Write-Host ""
    Write-Host "Iniciando AgentIA..." -ForegroundColor Cyan
    Write-Host ""
    .\agent.exe
}
