package gotour

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sync"
)

type Drone struct {
	ID  string  `json:"id"`
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
}

type DroneService struct {
	BaseURL string
	HTTPClient *http.Client
	APIKey string

	mu sync.Mutex // Protects AccessToken
	AccessToken string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
}

func NewDroneService(apiKey string, timeout time.Duration) *DroneService {
	return &DroneService{
		BaseURL: "https://api.picogrid.com/v1",
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}


func (s *DroneService) FindDrones(ctx context.Context, lat, lon float64, radius int) ([]Drone, error) {
	payload := struct {
		Center struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"center"`
		Radius int `json:"radius_meters"`
	}{}
	payload.Center.Lat = lat
	payload.Center.Lon = lon
	payload.Radius = radius

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/drones/search", s.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err	
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error")
	}

	// Security: Do not attempt to check the size of the response body. Hackers can manipulate
	// the Content-Length header to bypass checks. Instead, use http.MaxBytesReader to limit.

	const maxResponseSize = 1024 * 1024 // 1MB
	limitedBody := http.MaxBytesReader(nil, resp.Body, maxResponseSize)

	var result []Drone
	if err := json.NewDecoder(limitedBody).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return result, nil

}

func (s *DroneService) GetAccessToken(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Skip if we already have a valid token
	if s.AccessToken != "" {
		return s.AccessToken, nil
	}
	// Exchange API key for access token
	reqBody, _ := json.Marshal(map[string]string{
		"api_key": s.APIKey,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", s.BaseURL+"/auth/token", bytes.NewBuffer(reqBody))

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tk TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		return "", fmt.Errorf("Failed to decode token response: %w", err)
	}

	s.AccessToken = tk.AccessToken
	return s.AccessToken, nil
}




