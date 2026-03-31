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

const baseURL = "http://localhost:8080"
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

func printResponse(resp *http.Response) {
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error: failed to decode response:", err)
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}
