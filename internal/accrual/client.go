package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darkseear/go-musthave/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) GetAccrual(orderNumber int) (*models.Accrual, error) {
	url := fmt.Sprintf("%s/api/v1/accrual/%d", c.baseURL, orderNumber)
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
