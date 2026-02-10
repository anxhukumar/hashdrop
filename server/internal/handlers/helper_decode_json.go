package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Decode json
func DecodeJson(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		err := fmt.Errorf("invalid json: %w", err)
		return err
	}
	return nil
}
