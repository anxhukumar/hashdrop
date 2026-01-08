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

// PostJSON to the server and receive response
func PostJSON(endpoint string, reqBody, respBody any, token string) error {

	// Encode out data as json
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error encoding data to json: %w", err)
	}

	url := config.BaseURL + endpoint
	// Create a post request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating post request: %w", err)
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
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("server error (%d): %s", res.StatusCode, body)
	}

	// Decode json data from response
	if respBody != nil {
		if err := json.NewDecoder(res.Body).Decode(respBody); err != nil {
			return fmt.Errorf("response decode failed: %w", err)
		}
	}

	return nil
}
