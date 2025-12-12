package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Level representa el nivel de logging
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger maneja el logging de la aplicación
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	level       Level
}

// Config contiene la configuración del logger
type Config struct {
	Level      string
	File       string
	MaxSize    int
	MaxBackups int
}

// NewLogger crea una nueva instancia del logger
func NewLogger(config Config) (*Logger, error) {
	level := parseLevel(config.Level)

	// Crear directorio de logs si no existe
	if config.File != "" {
		if err := os.MkdirAll(filepath.Dir(config.File), 0755); err != nil {
			return nil, fmt.Errorf("error creando directorio de logs: %w", err)
		}
	}

	// Abrir archivo de log
	var writer io.Writer = os.Stdout
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("error abriendo archivo de log: %w", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
	}

	return &Logger{
		debugLogger: log.New(writer, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(writer, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(writer, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(writer, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}, nil
}

// Debug registra un mensaje de debug
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info registra un mensaje informativo
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn registra una advertencia
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error registra un error
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// parseLevel convierte un string a Level
func parseLevel(levelStr string) Level {
	switch levelStr {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// LogInteraction registra una interacción
func (l *Logger) LogInteraction(userInput, response, intent string) {
	l.Info("Interacción - Intent: %s | Input: %s | Response: %s", intent, userInput, response)
}

// LogError registra un error con contexto
func (l *Logger) LogError(component, operation string, err error) {
	l.Error("Error en %s.%s: %v", component, operation, err)
}

// LogStartup registra el inicio de la aplicación
func (l *Logger) LogStartup(version string) {
	l.Info("=== AgentIA Iniciado - Versión %s - %s ===", version, time.Now().Format(time.RFC3339))
}

// LogShutdown registra el cierre de la aplicación
func (l *Logger) LogShutdown() {
	l.Info("=== AgentIA Detenido - %s ===", time.Now().Format(time.RFC3339))
}
