//go:build !portaudio
// +build !portaudio

package speech

import (
	"context"
	"fmt"
)

// Recognizer stub cuando PortAudio no está disponible
type Recognizer struct {
	config Config
}

// NewRecognizer crea un reconocedor sin soporte de PortAudio
func NewRecognizer(config Config) (*Recognizer, error) {
	return &Recognizer{config: config}, nil
}

// StartListening no hace nada sin PortAudio
func (r *Recognizer) StartListening(ctx context.Context) error {
	return fmt.Errorf("reconocimiento de voz en vivo no disponible: compila con tag 'portaudio' para habilitarlo")
}

// StopListening no hace nada sin PortAudio
func (r *Recognizer) StopListening() error {
	return nil
}

// GetAudioData retorna error
func (r *Recognizer) GetAudioData() ([]byte, error) {
	return nil, fmt.Errorf("reconocimiento de voz en vivo no disponible: compila con tag 'portaudio' para habilitarlo")
}

// Close cierra el reconocedor
func (r *Recognizer) Close() error {
	return nil
}

// IsListening retorna false
func (r *Recognizer) IsListening() bool {
	return false
}
