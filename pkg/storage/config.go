package storage

// Config contiene la configuración de almacenamiento
type Config struct {
	Type           string
	Path           string
	BackupEnabled  bool
	BackupInterval int
}
