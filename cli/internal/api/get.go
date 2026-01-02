package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Make GET requests to the server and receive response
func GetJSON(endpoint string, respBody any, token string, queryParams map[string]string) error {

	url := config.BaseURL + endpoint

	// Create a get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating get requests: %w", err)
	}

	// Add query params if needed
	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set request headers
	req.Header.Set("Accept", "application/json")

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
