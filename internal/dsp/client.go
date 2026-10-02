package dsp

import (
	"adex/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration) *HTTPClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (c *HTTPClient) SendRequest(ctx context.Context, endpoint string, req domain.AuctionRequest) error {
	body, err := json.Marshal(req)

	if err != nil {
		return fmt.Errorf("ошибка marshal %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint, bytes.NewReader(body),
	)

	if err != nil {
		return fmt.Errorf("ошибка build request %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(httpReq)

	if err != nil {
		return fmt.Errorf("ошибка response %w", err)
	}

	defer func() {
		io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка StatusCode %d", response.StatusCode)
	}

	return nil
}
