package robot

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var HETZNER_SERVER_URL = "https://robot-ws.your-server.de"

type Client struct {
	Username string
	Password string
}

// Invoke GET HTTP request using Hetzner API. Path is added to the base URL.
func (c Client) Get(path string) ([]byte, error) {
	return c.MakeRequest("GET", path, nil)
}

// Invoke GET HTTP request using Hetzner API. Path is added to the base URL.
func (c Client) Post(path string, values url.Values) ([]byte, error) {
	return c.MakeRequest("POST", path, values)
}

// Invoke Hetzner API. Path is added to the base URL.
func (c Client) MakeRequest(action, path string, values url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", HETZNER_SERVER_URL, path)
	
	var parameters io.Reader
	if values != nil {
		parameters = strings.NewReader(values.Encode())
	}

	req, err := http.NewRequest(action, url, parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Response Body:", string(body))
	return body, nil
}
