package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/darkseear/go-musthave/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type RateLimitError struct {
	RetryAfter time.Duration
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %v", e.RetryAfter)
}

func (c *Client) GetAccrual(orderNumber string) (*models.Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)
	req, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if req.StatusCode == http.StatusOK {
		var accrual models.Accrual
		if err := json.Unmarshal(body, &accrual); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return &accrual, nil
	} else if req.StatusCode == http.StatusNoContent {
		return nil, nil
	} else if req.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("system url config: %s returned 404", url)
	} else {
		return nil, fmt.Errorf("failed to get accrual: %s", string(body))
	}
}
