package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func Delete(endpoint string, token string, queryParams map[string]string) error {

	url := config.BaseURL + endpoint

	// Create a delete request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating delete requests: %w", err)
	}

	// Add query params if needed
	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a new client and make the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("server error (%d): %s", res.StatusCode, body)
	}

	return nil
}
