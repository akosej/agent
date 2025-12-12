package speech

// Config contiene la configuración para el reconocimiento de voz
type Config struct {
	SampleRate int
	Channels   int
	Language   string
	Provider   string // "whisper-cpp", "whisper-api"
}
