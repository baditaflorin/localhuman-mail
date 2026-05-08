package security

import (
	"context"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/config"
)

type Capability struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Detail    string `json:"detail"`
}

func DetectCapabilities(ctx context.Context, cfg config.Config) []Capability {
	return []Capability{
		toolCapability("libpst/readpst", cfg.ReadPSTBin, "PST import can be enabled through readpst"),
		toolCapability("libmagic/file", cfg.FileBin, "MIME detection can use libmagic through file"),
		toolCapability("tantivy-cli", cfg.TantivyCLI, "Tantivy CLI can be wired as the search engine"),
		toolCapability("age", "age", "Runtime artifacts can be encrypted with age"),
		sentenceTransformersCapability(ctx),
		ollamaCapability(ctx, cfg.OllamaURL),
		{Name: "sqlite-search", Available: true, Detail: "Bundled local message store and fallback search"},
	}
}

func toolCapability(name, binary, detail string) Capability {
	if binary == "" {
		return Capability{Name: name, Available: false, Detail: "not configured"}
	}
	if _, err := exec.LookPath(binary); err != nil {
		return Capability{Name: name, Available: false, Detail: binary + " not found on PATH"}
	}
	return Capability{Name: name, Available: true, Detail: detail}
}

func sentenceTransformersCapability(ctx context.Context) Capability {
	command := exec.CommandContext(ctx, "python3", "-c", "import sentence_transformers")
	if err := command.Run(); err != nil {
		return Capability{Name: "sentence-transformers", Available: false, Detail: "python sentence_transformers package not found"}
	}
	return Capability{Name: "sentence-transformers", Available: true, Detail: "python sentence_transformers package import succeeded"}
}

func ollamaCapability(ctx context.Context, rawURL string) Capability {
	if strings.TrimSpace(rawURL) == "" {
		return Capability{Name: "local-llm", Available: false, Detail: "Ollama URL is not configured"}
	}
	client := http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(rawURL, "/")+"/api/tags", nil)
	if err != nil {
		return Capability{Name: "local-llm", Available: false, Detail: "invalid Ollama URL"}
	}
	response, err := client.Do(request)
	if err != nil {
		return Capability{Name: "local-llm", Available: false, Detail: "Ollama is not reachable"}
	}
	defer response.Body.Close()
	return Capability{Name: "local-llm", Available: response.StatusCode < 400, Detail: "Ollama API probe completed"}
}
