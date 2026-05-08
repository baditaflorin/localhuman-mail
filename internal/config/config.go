package config

import (
	"log/slog"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type LogLevel string

type Config struct {
	Addr           string   `envconfig:"ADDR" default:":8080"`
	DataDir        string   `envconfig:"DATA_DIR" default:"./runtime"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"https://baditaflorin.github.io,http://localhost:5173,http://127.0.0.1:5173"`
	LogLevel       LogLevel `envconfig:"LOG_LEVEL" default:"info"`
	OllamaURL      string   `envconfig:"OLLAMA_URL" default:"http://localhost:11434"`
	LLMModel       string   `envconfig:"LLM_MODEL" default:"llama3.2"`
	EmbeddingModel string   `envconfig:"EMBEDDING_MODEL" default:"nomic-embed-text"`
	AgeRecipient   string   `envconfig:"AGE_RECIPIENT"`
	TantivyCLI     string   `envconfig:"TANTIVY_CLI" default:"tantivy"`
	ReadPSTBin     string   `envconfig:"READPST_BIN" default:"readpst"`
	FileBin        string   `envconfig:"FILE_BIN" default:"file"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("LOCALHUMAN", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (level LogLevel) Level() slog.Level {
	switch strings.ToLower(string(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
