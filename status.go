package finago

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Status struct {
	Status struct {
		Indicator   string `json:"indicator"`
		Description string `json:"description"`
	} `json:"status"`
}

// Gets status from https://status.24sevenoffice.com/
func GetStatus() (*Status, error) {
	req, err := http.NewRequest("GET", "https://status.24sevenoffice.com/", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("finago: failed to create new request: %w", err)
	}
	req.Header.Set("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("finago: failed to do request: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("finago: failed to read body: %w", err)
	}

	var status Status
	if err = json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("finago: failed to unmarshal body: %w", err)
	}

	return &status, nil
}
