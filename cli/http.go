package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const baseURL = "http://localhost:8080"

func post(path string, body string) (*http.Response, error) {
	url := baseURL + path
	return http.Post(url, "application/json", strings.NewReader(body))
}

func get(path string) (*http.Response, error) {
	url := baseURL + path
	return http.Get(url)
}

func deleteReq(path string) (*http.Response, error) {
	url := baseURL + path
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
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
