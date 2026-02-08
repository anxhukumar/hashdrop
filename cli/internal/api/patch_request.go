package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func PatchJSON(endpoint string, reqBody, respBody any, token string) (int, error) {

	// Encode out data as json
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("error encoding data to json: %w", err)
	}

	url := config.BaseURL + endpoint
	// Create a patch request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating patch request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Set authorization header if token is provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Create a new client and make the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	status := res.StatusCode

	// If server returned error, read body and return error WITH status
	if status >= 400 {
		body, _ := io.ReadAll(res.Body)
		return status, fmt.Errorf("server error (%d): %s", status, body)
	}

	// Decode response if needed
	if respBody != nil {
		if err := json.NewDecoder(res.Body).Decode(respBody); err != nil {
			return status, fmt.Errorf("response decode failed: %w", err)
		}
	}

	return status, nil
}
