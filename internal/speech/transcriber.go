package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Transcriber maneja la transcripción de audio a texto usando whisper.cpp local
type Transcriber struct {
	language      string
	modelPath     string // Ruta al modelo de Whisper
	whisperPath   string // Ruta al ejecutable de whisper.cpp
	useWhisperCpp bool   // Si usar whisper.cpp o API local
	apiURL        string // URL de API local de Whisper (si se usa)
}

// NewTranscriber crea una nueva instancia del transcriptor
// whisperPath: ruta al ejecutable main de whisper.cpp (ej: ./whisper.cpp/main)
// modelPath: ruta al modelo .bin (ej: ./models/ggml-base.bin)
func NewTranscriber(whisperPath, modelPath, language string) *Transcriber {
	return &Transcriber{
		language:      language,
		modelPath:     modelPath,
		whisperPath:   whisperPath,
		useWhisperCpp: true,
		apiURL:        "",
	}
}

// NewTranscriberWithAPI crea un transcriptor que usa una API local de Whisper
func NewTranscriberWithAPI(apiURL, language string) *Transcriber {
	if apiURL == "" {
		apiURL = "http://localhost:8000" // Puerto por defecto
	}
	return &Transcriber{
		language:      language,
		useWhisperCpp: false,
		apiURL:        apiURL,
	}
}

// TranscribeFile transcribe un archivo de audio a texto usando whisper.cpp
func (t *Transcriber) TranscribeFile(ctx context.Context, audioPath string) (string, error) {
	if t.useWhisperCpp {
		return t.transcribeWithWhisperCpp(ctx, audioPath)
	}
	return t.transcribeWithAPI(ctx, audioPath)
}

// transcribeWithWhisperCpp usa el ejecutable de whisper.cpp directamente
func (t *Transcriber) transcribeWithWhisperCpp(ctx context.Context, audioPath string) (string, error) {
	if t.whisperPath == "" {
		return "", fmt.Errorf("whisperPath no configurado. Instala whisper.cpp y configura la ruta")
	}

	if t.modelPath == "" {
		return "", fmt.Errorf("modelPath no configurado. Descarga un modelo de Whisper")
	}

	// Ejecutar whisper.cpp: ./main -m model.bin -f audio.wav -l es
	args := []string{
		"-m", t.modelPath,
		"-f", audioPath,
		"-l", t.language,
		"-nt", // Sin timestamps
	}

	cmd := exec.CommandContext(ctx, t.whisperPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error ejecutando whisper.cpp: %w\nOutput: %s", err, string(output))
	}

	// Parsear output de whisper.cpp (buscar líneas que comienzan con texto)
	lines := strings.Split(string(output), "\n")
	var transcription strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Filtrar líneas de log, mantener solo transcripción
		if line != "" && !strings.HasPrefix(line, "[") && !strings.Contains(line, "whisper_") {
			transcription.WriteString(line)
			transcription.WriteString(" ")
		}
	}

	result := strings.TrimSpace(transcription.String())
	if result == "" {
		return "", fmt.Errorf("no se pudo transcribir el audio")
	}

	return result, nil
}

// transcribeWithAPI usa una API local de Whisper (como whisper-api o faster-whisper)
func (t *Transcriber) transcribeWithAPI(ctx context.Context, audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	// Agregar parámetros
	writer.WriteField("language", t.language)
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", t.apiURL+"/transcribe", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error llamando a API de Whisper: %w (Asegúrate de que el servidor esté corriendo)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API respondió con error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parsear respuesta JSON
	var result struct {
		Text string `json:"text"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	return result.Text, nil
}

// TranscribeStream transcribe audio desde un stream
func (t *Transcriber) TranscribeStream(ctx context.Context, audioData []byte, format string) (string, error) {
	// Crear archivo temporal
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("audio-*.%s", format))
	if err != nil {
		return "", fmt.Errorf("error creando archivo temporal: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(audioData); err != nil {
		return "", fmt.Errorf("error escribiendo datos de audio: %w", err)
	}

	return t.TranscribeFile(ctx, tmpFile.Name())
}
