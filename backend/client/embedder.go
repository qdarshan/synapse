package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	LMSTUDIO_EMBEDDINGS_URL = "http://127.0.0.1:1234/v1/embeddings"
	LMSTUDIO_MODEL          = "text-embedding-nomic-embed-text-v1.5"
	HTTP_TIMEOUT            = 30 * time.Second
)

type LMStudioRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type LMStudioResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func GenerateEmbedding(ctx context.Context, input string) ([]float64, error) {
	requestPayload := LMStudioRequest{
		Model: LMSTUDIO_MODEL,
		Input: input,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		logger.Error("Failed to marshal request", "error", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", LMSTUDIO_EMBEDDINGS_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Error("Failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: HTTP_TIMEOUT}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to execute request to LM Studio", "error", err)
		return nil, fmt.Errorf("failed to execute request to LM Studio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("LM Studio returned non-200 status code", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("LM Studio returned non-200 status code: %d", resp.StatusCode)
	}

	var lmStudioResponse LMStudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&lmStudioResponse); err != nil {
		logger.Error("Failed to decode response body", "error", err)
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(lmStudioResponse.Data) == 0 || len(lmStudioResponse.Data[0].Embedding) == 0 {
		logger.Error("LM Studio returned empty embedding data")
		return nil, fmt.Errorf("LM Studio returned empty embedding data")
	}

	logger.Debug("Successfully generated embedding", "input_length", len(input))
	return lmStudioResponse.Data[0].Embedding, nil
}
