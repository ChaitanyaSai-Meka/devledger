package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseURL = "http://localhost:38080"
const requestTimeout = 5 * time.Second

func post(path string, body string) (*http.Response, error) {
	url := baseURL + path
	return doRequest(http.MethodPost, url, strings.NewReader(body), "application/json")
}

func get(path string) (*http.Response, error) {
	url := baseURL + path
	return doRequest(http.MethodGet, url, nil, "")
}

func deleteReq(path string) (*http.Response, error) {
	url := baseURL + path
	return doRequest(http.MethodDelete, url, nil, "")
}

func doRequest(method string, url string, body io.Reader, contentType string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return http.DefaultClient.Do(req)
}

func printResponse(resp *http.Response) error {
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		data, err := json.MarshalIndent(result["data"], "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response data: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("Error:", result["error"])
	return fmt.Errorf("request failed with status %s", resp.Status)
}
