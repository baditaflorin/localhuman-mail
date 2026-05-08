package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/config"
	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
)

type Service struct {
	ollamaURL string
	llmModel  string
	client    *http.Client
}

type ReplyResult struct {
	Draft  string `json:"draft"`
	Model  string `json:"model"`
	Source string `json:"source"`
}

func NewService(cfg config.Config) *Service {
	return &Service{
		ollamaURL: strings.TrimRight(cfg.OllamaURL, "/"),
		llmModel:  cfg.LLMModel,
		client: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
}

func (service *Service) AssistReply(ctx context.Context, message mailbox.Message, tone, instructions string) ReplyResult {
	if service.ollamaURL == "" {
		return fallbackReply(message, tone)
	}
	draft, err := service.generateWithOllama(ctx, message, tone, instructions)
	if err != nil || strings.TrimSpace(draft) == "" {
		return fallbackReply(message, tone)
	}
	return ReplyResult{Draft: strings.TrimSpace(draft), Model: service.llmModel, Source: "local_llm"}
}

func (service *Service) generateWithOllama(ctx context.Context, message mailbox.Message, tone, instructions string) (string, error) {
	payload := map[string]any{
		"model":  service.llmModel,
		"stream": false,
		"prompt": prompt(message, tone, instructions),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal ollama request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, service.ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create ollama request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := service.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("ollama status %d", response.StatusCode)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read ollama response: %w", err)
	}
	var decoded struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return "", fmt.Errorf("decode ollama response: %w", err)
	}
	return decoded.Response, nil
}

func prompt(message mailbox.Message, tone, instructions string) string {
	return fmt.Sprintf(`Draft a %s reply to this email. Keep private data local. Do not invent facts.

Extra instructions: %s

From: %s
Subject: %s
Body:
%s`, tone, instructions, message.From, message.Subject, message.Body)
}

func fallbackReply(message mailbox.Message, tone string) ReplyResult {
	opener := "Got it."
	switch tone {
	case "warm":
		opener = "Thanks for the thoughtful context."
	case "decisive":
		opener = "I can take this forward."
	}
	draft := fmt.Sprintf("%s\n\nI reviewed \"%s\" and will follow up with the key decision, owner, and next step separated clearly. I will keep the response concise and flag anything that needs approval.\n\nBest,\n", opener, message.Subject)
	return ReplyResult{Draft: draft, Model: "fallback-template", Source: "fallback"}
}
