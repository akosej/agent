@echo off
REM Script de inicio rápido para AgentIA en Windows

echo ========================================
echo AgentIA - Configuracion Inicial
echo ========================================
echo.

REM Verificar si existe .env
if not exist .env (
    echo [1/5] Creando archivo .env desde plantilla...
    copy .env.example .env
    echo IMPORTANTE: Edita el archivo .env y agrega tu OPENAI_API_KEY
    echo.
    pause
) else (
    echo [1/5] Archivo .env encontrado
)

REM Crear directorios necesarios
echo [2/5] Creando directorios...
if not exist logs mkdir logs
if not exist data mkdir data
if not exist data\conversations mkdir data\conversations
if not exist data\models mkdir data\models

REM Descargar dependencias
echo [3/5] Descargando dependencias de Go...
go mod download
go get github.com/mattn/go-sqlite3

REM Compilar el proyecto
echo [4/5] Compilando el proyecto...
go build -o agent.exe .\cmd\agent

if %errorlevel% neq 0 (
    echo ERROR: La compilacion fallo
    echo Verifica que tengas Go instalado correctamente
    pause
    exit /b 1
)

echo [5/5] Compilacion exitosa!
echo.
echo ========================================
echo AgentIA esta listo para usar
echo ========================================
echo.
echo Para ejecutar: agent.exe
echo Para ayuda: Ver README.md
echo.
pause

REM Preguntar si quiere ejecutar ahora
echo ¿Desea ejecutar AgentIA ahora? (S/N)
set /p respuesta=

if /i "%respuesta%"=="S" (
    echo.
    echo Iniciando AgentIA...
    echo.
    agent.exe
)
