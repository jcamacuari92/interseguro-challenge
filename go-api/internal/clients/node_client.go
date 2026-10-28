package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/interseguro/go-api/internal/models"
)

// NodeClient encapsula la comunicación HTTP con la API de Node.js.
type NodeClient struct {
	baseURL string
	http    *http.Client
}

func NewNodeClient(baseURL string) *NodeClient {
	return &NodeClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStatistics envía las matrices a Node.js y obtiene las estadísticas.
func (c *NodeClient) GetStatistics(ctx context.Context, matrices [][][]float64) (*models.StatsResponse, error) {
	payload, err := json.Marshal(models.StatsRequest{Matrices: matrices})
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/statistics", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error llamando a Node API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Node API respondió %d: %s", resp.StatusCode, string(body))
	}

	var stats models.StatsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de Node: %w", err)
	}
	return &stats, nil
}