package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func GenerateEmbedding(input string) ([]float64, error) {
	requestPayload := LMStudioRequest{
		Model: LMSTUDIO_MODEL,
		Input: input,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", LMSTUDIO_EMBEDDINGS_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: HTTP_TIMEOUT}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to LM Studio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LM Studio returned non-200 status code: %d", resp.StatusCode)
	}

	var lmStudioResponse LMStudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&lmStudioResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(lmStudioResponse.Data) == 0 || len(lmStudioResponse.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("LM Studio returned empty embedding data")
	}

	return lmStudioResponse.Data[0].Embedding, nil
}
