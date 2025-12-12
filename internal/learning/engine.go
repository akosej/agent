package learning

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Interaction representa una interacción del usuario
type Interaction struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	UserInput string                 `json:"user_input"`
	Response  string                 `json:"response"`
	Intent    string                 `json:"intent"`
	Feedback  *Feedback              `json:"feedback,omitempty"`
	Context   map[string]interface{} `json:"context"`
}

// Feedback representa retroalimentación del usuario
type Feedback struct {
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Pattern representa un patrón aprendido
type Pattern struct {
	Pattern    string    `json:"pattern"`
	Response   string    `json:"response"`
	Frequency  int       `json:"frequency"`
	Confidence float64   `json:"confidence"`
	LastUsed   time.Time `json:"last_used"`
}

// KnowledgeBase almacena el conocimiento aprendido
type KnowledgeBase struct {
	Patterns     map[string]*Pattern `json:"patterns"`
	Interactions []*Interaction      `json:"interactions"`
	Stats        *Stats              `json:"stats"`
	mu           sync.RWMutex
}

// Stats contiene estadísticas del agente
type Stats struct {
	TotalInteractions int       `json:"total_interactions"`
	PositiveFeedback  int       `json:"positive_feedback"`
	NegativeFeedback  int       `json:"negative_feedback"`
	AverageRating     float64   `json:"average_rating"`
	LastUpdated       time.Time `json:"last_updated"`
}

// Config contiene la configuración del sistema de aprendizaje
type Config struct {
	LearningRate        float64
	ConfidenceThreshold float64
	MaxInteractions     int
	SaveInterval        int
}

// Engine maneja el aprendizaje del agente
type Engine struct {
	kb     *KnowledgeBase
	config Config
}

// NewEngine crea una nueva instancia del motor de aprendizaje
func NewEngine(config Config) *Engine {
	return &Engine{
		kb: &KnowledgeBase{
			Patterns:     make(map[string]*Pattern),
			Interactions: make([]*Interaction, 0),
			Stats: &Stats{
				LastUpdated: time.Now(),
			},
		},
		config: config,
	}
}

// RecordInteraction registra una nueva interacción
func (e *Engine) RecordInteraction(interaction *Interaction) {
	e.kb.mu.Lock()
	defer e.kb.mu.Unlock()

	interaction.ID = fmt.Sprintf("int_%d", time.Now().UnixNano())
	interaction.Timestamp = time.Now()

	e.kb.Interactions = append(e.kb.Interactions, interaction)
	e.kb.Stats.TotalInteractions++
	e.kb.Stats.LastUpdated = time.Now()

	// Limitar el número de interacciones almacenadas
	if len(e.kb.Interactions) > e.config.MaxInteractions {
		e.kb.Interactions = e.kb.Interactions[len(e.kb.Interactions)-e.config.MaxInteractions:]
	}

	// Intentar extraer un patrón
	e.learnPattern(interaction)
}

// learnPattern intenta aprender un patrón de la interacción
func (e *Engine) learnPattern(interaction *Interaction) {
	// Simplificación: usar el intent como clave del patrón
	patternKey := interaction.Intent

	if pattern, exists := e.kb.Patterns[patternKey]; exists {
		pattern.Frequency++
		pattern.LastUsed = time.Now()
		// Ajustar confianza basado en feedback
		if interaction.Feedback != nil && interaction.Feedback.Rating >= 4 {
			pattern.Confidence = min(1.0, pattern.Confidence+e.config.LearningRate)
		}
	} else {
		e.kb.Patterns[patternKey] = &Pattern{
			Pattern:    interaction.UserInput,
			Response:   interaction.Response,
			Frequency:  1,
			Confidence: 0.5,
			LastUsed:   time.Now(),
		}
	}
}

// AddFeedback añade retroalimentación a una interacción
func (e *Engine) AddFeedback(interactionID string, rating int, comment string) error {
	e.kb.mu.Lock()
	defer e.kb.mu.Unlock()

	for _, interaction := range e.kb.Interactions {
		if interaction.ID == interactionID {
			interaction.Feedback = &Feedback{
				Rating:    rating,
				Comment:   comment,
				Timestamp: time.Now(),
			}

			// Actualizar estadísticas
			if rating >= 4 {
				e.kb.Stats.PositiveFeedback++
			} else if rating <= 2 {
				e.kb.Stats.NegativeFeedback++
			}

			// Recalcular rating promedio
			totalRatings := 0.0
			count := 0
			for _, inter := range e.kb.Interactions {
				if inter.Feedback != nil {
					totalRatings += float64(inter.Feedback.Rating)
					count++
				}
			}
			if count > 0 {
				e.kb.Stats.AverageRating = totalRatings / float64(count)
			}

			return nil
		}
	}

	return fmt.Errorf("interacción no encontrada: %s", interactionID)
}

// FindSimilarPattern busca un patrón similar en la base de conocimiento
func (e *Engine) FindSimilarPattern(intent string) (*Pattern, bool) {
	e.kb.mu.RLock()
	defer e.kb.mu.RUnlock()

	pattern, exists := e.kb.Patterns[intent]
	if !exists {
		return nil, false
	}

	// Solo retornar si la confianza supera el umbral
	if pattern.Confidence < e.config.ConfidenceThreshold {
		return nil, false
	}

	return pattern, true
}

// GetStats obtiene las estadísticas actuales
func (e *Engine) GetStats() *Stats {
	e.kb.mu.RLock()
	defer e.kb.mu.RUnlock()

	return e.kb.Stats
}

// GetRecentInteractions obtiene las N interacciones más recientes
func (e *Engine) GetRecentInteractions(n int) []*Interaction {
	e.kb.mu.RLock()
	defer e.kb.mu.RUnlock()

	if len(e.kb.Interactions) <= n {
		return e.kb.Interactions
	}

	return e.kb.Interactions[len(e.kb.Interactions)-n:]
}

// Export exporta la base de conocimiento a JSON
func (e *Engine) Export() ([]byte, error) {
	e.kb.mu.RLock()
	defer e.kb.mu.RUnlock()

	return json.MarshalIndent(e.kb, "", "  ")
}

// Import importa una base de conocimiento desde JSON
func (e *Engine) Import(data []byte) error {
	e.kb.mu.Lock()
	defer e.kb.mu.Unlock()

	return json.Unmarshal(data, e.kb)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
