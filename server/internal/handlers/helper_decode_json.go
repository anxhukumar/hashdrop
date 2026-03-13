package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const maxBodySize = 1 << 20 // 1 MB

// Decode json
func DecodeJson(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		err := fmt.Errorf("invalid json: %w", err)
		return err
	}

	if decoder.More() {
		return fmt.Errorf("only one JSON object allowed")
	}

	return nil
}
