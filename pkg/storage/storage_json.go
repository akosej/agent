//go:build !cgo
// +build !cgo

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage maneja el almacenamiento persistente usando archivos JSON (sin CGO)
type Storage struct {
	config       Config
	dataDir      string
	mu           sync.RWMutex
	interactions []map[string]interface{}
	patterns     []map[string]interface{}
	stats        map[string]interface{}
}

// NewStorage crea una nueva instancia de almacenamiento basado en archivos JSON
func NewStorage(config Config) (*Storage, error) {
	// Crear directorio de datos si no existe
	dataDir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de datos: %w", err)
	}

	storage := &Storage{
		config:       config,
		dataDir:      dataDir,
		interactions: make([]map[string]interface{}, 0),
		patterns:     make([]map[string]interface{}, 0),
		stats:        make(map[string]interface{}),
	}

	// Cargar datos existentes
	if err := storage.loadData(); err != nil {
		// Si no existe, no es error, solo inicia vacío
		fmt.Printf("Iniciando almacenamiento nuevo en: %s\n", dataDir)
	}

	return storage, nil
}

// loadData carga datos desde archivos JSON
func (s *Storage) loadData() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cargar interacciones
	interactionsPath := filepath.Join(s.dataDir, "interactions.json")
	if data, err := os.ReadFile(interactionsPath); err == nil {
		json.Unmarshal(data, &s.interactions)
	}

	// Cargar patrones
	patternsPath := filepath.Join(s.dataDir, "patterns.json")
	if data, err := os.ReadFile(patternsPath); err == nil {
		json.Unmarshal(data, &s.patterns)
	}

	// Cargar estadísticas
	statsPath := filepath.Join(s.dataDir, "stats.json")
	if data, err := os.ReadFile(statsPath); err == nil {
		json.Unmarshal(data, &s.stats)
	}

	return nil
}

// saveData guarda datos en archivos JSON
func (s *Storage) saveData() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Guardar interacciones
	interactionsPath := filepath.Join(s.dataDir, "interactions.json")
	if data, err := json.MarshalIndent(s.interactions, "", "  "); err == nil {
		os.WriteFile(interactionsPath, data, 0644)
	}

	// Guardar patrones
	patternsPath := filepath.Join(s.dataDir, "patterns.json")
	if data, err := json.MarshalIndent(s.patterns, "", "  "); err == nil {
		os.WriteFile(patternsPath, data, 0644)
	}

	// Guardar estadísticas
	statsPath := filepath.Join(s.dataDir, "stats.json")
	if data, err := json.MarshalIndent(s.stats, "", "  "); err == nil {
		os.WriteFile(statsPath, data, 0644)
	}

	return nil
}

// SaveInteraction guarda una interacción
func (s *Storage) SaveInteraction(interaction map[string]interface{}) error {
	s.mu.Lock()
	s.interactions = append(s.interactions, interaction)
	s.mu.Unlock()

	return s.saveData()
}

// GetInteractions obtiene interacciones con límite
func (s *Storage) GetInteractions(limit int) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > len(s.interactions) {
		limit = len(s.interactions)
	}

	// Retornar las últimas 'limit' interacciones
	start := len(s.interactions) - limit
	if start < 0 {
		start = 0
	}

	return s.interactions[start:], nil
}

// SavePattern guarda un patrón aprendido
func (s *Storage) SavePattern(pattern map[string]interface{}) error {
	s.mu.Lock()
	s.patterns = append(s.patterns, pattern)
	s.mu.Unlock()

	return s.saveData()
}

// GetPatterns obtiene patrones aprendidos
func (s *Storage) GetPatterns() ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.patterns, nil
}

// UpdateStats actualiza las estadísticas
func (s *Storage) UpdateStats(stats map[string]interface{}) error {
	s.mu.Lock()
	for k, v := range stats {
		s.stats[k] = v
	}
	s.stats["last_updated"] = time.Now()
	s.mu.Unlock()

	return s.saveData()
}

// GetStats obtiene las estadísticas
func (s *Storage) GetStats() (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stats, nil
}

// SaveConversation guarda una conversación completa
func (s *Storage) SaveConversation(conversation map[string]interface{}) error {
	conversationsDir := filepath.Join(s.dataDir, "conversations")
	os.MkdirAll(conversationsDir, 0755)

	timestamp := time.Now().Unix()
	filename := filepath.Join(conversationsDir, fmt.Sprintf("conversation_%d.json", timestamp))

	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Close cierra el almacenamiento (guarda datos pendientes)
func (s *Storage) Close() error {
	return s.saveData()
}

// Backup crea una copia de seguridad
func (s *Storage) Backup() error {
	if !s.config.BackupEnabled {
		return nil
	}

	backupDir := filepath.Join(s.dataDir, "backups")
	os.MkdirAll(backupDir, 0755)

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%s.json", timestamp))

	s.mu.RLock()
	backupData := map[string]interface{}{
		"interactions": s.interactions,
		"patterns":     s.patterns,
		"stats":        s.stats,
		"timestamp":    time.Now(),
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0644)
}
