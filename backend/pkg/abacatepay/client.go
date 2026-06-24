package abacatepay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.abacatepay.com/v2"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
	}
}

func (c *Client) do(method, path string, body, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("abacatepay: invalid API key (401)")
	}

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("abacatepay: forbidden (403)")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("abacatepay: rate limited (429)")
	}

	if resp.StatusCode >= 400 {
		errMsg := apiResp.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("unexpected status %d", resp.StatusCode)
		}
		return fmt.Errorf("abacatepay: %s", errMsg)
	}

	if result != nil && apiResp.Data != nil {
		data, err := json.Marshal(apiResp.Data)
		if err != nil {
			return fmt.Errorf("marshal api data: %w", err)
		}
		if err := json.Unmarshal(data, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}
