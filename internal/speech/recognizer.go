//go:build portaudio
// +build portaudio

package speech

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
)

// Recognizer maneja el reconocimiento de voz con PortAudio
type Recognizer struct {
	config Config
	stream *portaudio.Stream
}

// NewRecognizer crea una nueva instancia del reconocedor de voz
func NewRecognizer(config Config) (*Recognizer, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("error inicializando PortAudio: %w", err)
	}

	return &Recognizer{
		config: config,
	}, nil
}

// StartListening comienza a escuchar audio del micrófono
func (r *Recognizer) StartListening(ctx context.Context) (<-chan []byte, error) {
	audioChan := make(chan []byte, 10)

	buffer := make([]int16, 1024)
	stream, err := portaudio.OpenDefaultStream(
		r.config.Channels,
		0,
		float64(r.config.SampleRate),
		len(buffer),
		buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("error abriendo stream de audio: %w", err)
	}

	r.stream = stream

	if err := stream.Start(); err != nil {
		return nil, fmt.Errorf("error iniciando stream: %w", err)
	}

	go func() {
		defer close(audioChan)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := stream.Read(); err != nil {
					fmt.Printf("Error leyendo audio: %v\n", err)
					return
				}

				// Convertir int16 a bytes
				audioBytes := make([]byte, len(buffer)*2)
				for i, sample := range buffer {
					audioBytes[i*2] = byte(sample)
					audioBytes[i*2+1] = byte(sample >> 8)
				}

				select {
				case audioChan <- audioBytes:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return audioChan, nil
}

// StopListening detiene la escucha de audio
func (r *Recognizer) StopListening() error {
	if r.stream != nil {
		if err := r.stream.Stop(); err != nil {
			return fmt.Errorf("error deteniendo stream: %w", err)
		}
	}
	return nil
}

// RecordToFile graba audio a un archivo WAV
func (r *Recognizer) RecordToFile(ctx context.Context, filename string, duration time.Duration) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo: %w", err)
	}
	defer file.Close()

	// Escribir encabezado WAV
	if err := r.writeWAVHeader(file); err != nil {
		return err
	}

	audioChan, err := r.StartListening(ctx)
	if err != nil {
		return err
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			r.StopListening()
			return r.updateWAVHeader(file)
		case audioData, ok := <-audioChan:
			if !ok {
				return r.updateWAVHeader(file)
			}
			if _, err := file.Write(audioData); err != nil {
				return fmt.Errorf("error escribiendo audio: %w", err)
			}
		case <-ctx.Done():
			r.StopListening()
			return ctx.Err()
		}
	}
}

// writeWAVHeader escribe el encabezado WAV básico
func (r *Recognizer) writeWAVHeader(w io.Writer) error {
	header := []byte{
		'R', 'I', 'F', 'F',
		0, 0, 0, 0, // tamaño del archivo (se actualizará después)
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		16, 0, 0, 0, // tamaño del chunk fmt
		1, 0, // formato PCM
		byte(r.config.Channels), 0,
		byte(r.config.SampleRate), byte(r.config.SampleRate >> 8), byte(r.config.SampleRate >> 16), byte(r.config.SampleRate >> 24),
		0, 0, 0, 0, // byte rate (se calculará)
		2, 0, // block align
		16, 0, // bits por muestra
		'd', 'a', 't', 'a',
		0, 0, 0, 0, // tamaño de datos (se actualizará después)
	}

	_, err := w.Write(header)
	return err
}

// updateWAVHeader actualiza el tamaño del archivo WAV
func (r *Recognizer) updateWAVHeader(file *os.File) error {
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// Actualizar tamaño total del archivo
	if _, err := file.Seek(4, io.SeekStart); err != nil {
		return err
	}
	totalSize := uint32(size - 8)
	if err := writeLittleEndianUint32(file, totalSize); err != nil {
		return err
	}

	// Actualizar tamaño de datos
	if _, err := file.Seek(40, io.SeekStart); err != nil {
		return err
	}
	dataSize := uint32(size - 44)
	return writeLittleEndianUint32(file, dataSize)
}

func writeLittleEndianUint32(w io.Writer, value uint32) error {
	bytes := []byte{
		byte(value),
		byte(value >> 8),
		byte(value >> 16),
		byte(value >> 24),
	}
	_, err := w.Write(bytes)
	return err
}

// Close limpia los recursos
func (r *Recognizer) Close() error {
	if err := r.StopListening(); err != nil {
		return err
	}
	return portaudio.Terminate()
}
