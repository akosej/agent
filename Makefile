# Makefile para AgentIA

.PHONY: all build run clean test install deps

# Variables
BINARY_NAME=agent
BINARY_WINDOWS=$(BINARY_NAME).exe
MAIN_PATH=./cmd/agent
BUILD_DIR=./build

# Comandos principales
all: deps build

# Instalar dependencias
deps:
	@echo "Instalando dependencias..."
	go mod download
	go get github.com/mattn/go-sqlite3

# Compilar el proyecto
build:
	@echo "Compilando $(BINARY_NAME)..."
	go build -o $(BINARY_WINDOWS) $(MAIN_PATH)

# Compilar con optimizaciones
build-release:
	@echo "Compilando versión release..."
	go build -ldflags="-s -w" -o $(BINARY_WINDOWS) $(MAIN_PATH)

# Ejecutar en modo desarrollo
run:
	@echo "Ejecutando en modo desarrollo..."
	go run $(MAIN_PATH)

# Ejecutar tests
test:
	@echo "Ejecutando tests..."
	go test -v ./...

# Ejecutar tests con coverage
test-coverage:
	@echo "Ejecutando tests con coverage..."
	go test -cover ./...

# Verificar código
vet:
	@echo "Verificando código..."
	go vet ./...

# Formatear código
fmt:
	@echo "Formateando código..."
	go fmt ./...

# Limpiar archivos compilados
clean:
	@echo "Limpiando archivos..."
	@if exist $(BINARY_WINDOWS) del $(BINARY_WINDOWS)
	@if exist data\agent.db del data\agent.db
	@if exist logs\agent.log del logs\agent.log
	@echo "Limpieza completada"

# Crear directorios necesarios
setup:
	@echo "Creando estructura de directorios..."
	@if not exist logs mkdir logs
	@if not exist data mkdir data
	@if not exist data\conversations mkdir data\conversations
	@if not exist data\models mkdir data\models
	@echo "Estructura creada"

# Instalar herramientas de desarrollo
install-tools:
	@echo "Instalando herramientas de desarrollo..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint del código
lint:
	@echo "Ejecutando linter..."
	golangci-lint run

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make deps           - Instalar dependencias"
	@echo "  make build          - Compilar el proyecto"
	@echo "  make build-release  - Compilar versión optimizada"
	@echo "  make run            - Ejecutar en modo desarrollo"
	@echo "  make test           - Ejecutar tests"
	@echo "  make test-coverage  - Ejecutar tests con coverage"
	@echo "  make vet            - Verificar código"
	@echo "  make fmt            - Formatear código"
	@echo "  make clean          - Limpiar archivos generados"
	@echo "  make setup          - Crear estructura de directorios"
	@echo "  make lint           - Ejecutar linter"
	@echo "  make help           - Mostrar esta ayuda"
