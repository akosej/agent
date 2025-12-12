//go:build cgo
// +build cgo

package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Storage maneja el almacenamiento persistente usando SQLite
type Storage struct {
	db     *sql.DB
	config Config
}

// NewStorage crea una nueva instancia de almacenamiento
func NewStorage(config Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, fmt.Errorf("error abriendo base de datos: %w", err)
	}

	storage := &Storage{
		db:     db,
		config: config,
	}

	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

// initTables inicializa las tablas de la base de datos
func (s *Storage) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS interactions (
		id TEXT PRIMARY KEY,
		timestamp DATETIME,
		user_input TEXT,
		response TEXT,
		intent TEXT,
		context TEXT,
		feedback_rating INTEGER,
		feedback_comment TEXT,
		feedback_timestamp DATETIME
	);

	CREATE TABLE IF NOT EXISTS patterns (
		pattern_key TEXT PRIMARY KEY,
		pattern TEXT,
		response TEXT,
		frequency INTEGER,
		confidence REAL,
		last_used DATETIME
	);

	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY,
		total_interactions INTEGER,
		positive_feedback INTEGER,
		negative_feedback INTEGER,
		average_rating REAL,
		last_updated DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_interactions_timestamp ON interactions(timestamp);
	CREATE INDEX IF NOT EXISTS idx_interactions_intent ON interactions(intent);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveInteraction guarda una interacción en la base de datos
func (s *Storage) SaveInteraction(interaction map[string]interface{}) error {
	contextJSON, _ := json.Marshal(interaction["context"])

	query := `
	INSERT INTO interactions (id, timestamp, user_input, response, intent, context)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		interaction["id"],
		interaction["timestamp"],
		interaction["user_input"],
		interaction["response"],
		interaction["intent"],
		string(contextJSON),
	)

	return err
}

// SavePattern guarda un patrón aprendido
func (s *Storage) SavePattern(patternKey string, pattern map[string]interface{}) error {
	query := `
	INSERT OR REPLACE INTO patterns (pattern_key, pattern, response, frequency, confidence, last_used)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		patternKey,
		pattern["pattern"],
		pattern["response"],
		pattern["frequency"],
		pattern["confidence"],
		pattern["last_used"],
	)

	return err
}

// GetRecentInteractions obtiene las interacciones más recientes
func (s *Storage) GetRecentInteractions(limit int) ([]map[string]interface{}, error) {
	query := `
	SELECT id, timestamp, user_input, response, intent, context
	FROM interactions
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []map[string]interface{}
	for rows.Next() {
		var id, userInput, response, intent, contextJSON string
		var timestamp time.Time

		if err := rows.Scan(&id, &timestamp, &userInput, &response, &intent, &contextJSON); err != nil {
			continue
		}

		var context map[string]interface{}
		json.Unmarshal([]byte(contextJSON), &context)

		interactions = append(interactions, map[string]interface{}{
			"id":         id,
			"timestamp":  timestamp,
			"user_input": userInput,
			"response":   response,
			"intent":     intent,
			"context":    context,
		})
	}

	return interactions, nil
}

// GetPattern obtiene un patrón por su clave
func (s *Storage) GetPattern(patternKey string) (map[string]interface{}, error) {
	query := `
	SELECT pattern, response, frequency, confidence, last_used
	FROM patterns
	WHERE pattern_key = ?
	`

	var pattern, response string
	var frequency int
	var confidence float64
	var lastUsed time.Time

	err := s.db.QueryRow(query, patternKey).Scan(&pattern, &response, &frequency, &confidence, &lastUsed)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pattern":    pattern,
		"response":   response,
		"frequency":  frequency,
		"confidence": confidence,
		"last_used":  lastUsed,
	}, nil
}

// UpdateStats actualiza las estadísticas
func (s *Storage) UpdateStats(stats map[string]interface{}) error {
	query := `
	INSERT OR REPLACE INTO stats (id, total_interactions, positive_feedback, negative_feedback, average_rating, last_updated)
	VALUES (1, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		stats["total_interactions"],
		stats["positive_feedback"],
		stats["negative_feedback"],
		stats["average_rating"],
		time.Now(),
	)

	return err
}

// Close cierra la conexión a la base de datos
func (s *Storage) Close() error {
	return s.db.Close()
}
