package smartapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	ClientId     string
	Mpin         string
	APIKey       string
	TOTPKey      string
	JWTToken     string
	RefreshToken string
	FeedToken    string
	baseUrl      string
	localIP      string
	publicIP     string
	macID        string
	httpClient   *http.Client
}

func NewClient(clientId, mpin, apiKey, totpKey string) *Client {

	httpClient := &http.Client{
		Timeout: REQUEST_TIMEOUT,
	}

	return &Client{
		ClientId:     clientId,
		Mpin:         mpin,
		APIKey:       apiKey,
		TOTPKey:      totpKey,
		JWTToken:     "",
		RefreshToken: "",
		FeedToken:    "",
		localIP:      localIP(),
		publicIP:     publicIP(),
		macID:        macID(),
		baseUrl:      BASE_URL,
		httpClient:   httpClient,
	}
}

type serverResponse[T any] struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      T      `json:"data"`
}

func (c *Client) SetBaseUrl(url string) {
	c.baseUrl = url
}

func (c *Client) doRequest(method string, path string, body any, result any) error {

	// Build URL
	url := c.baseUrl + path

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		// If caller passed an io.Reader, allow that through
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		default:
			b, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Default headers for JSON API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", c.APIKey)
	req.Header.Set("X-ClientLocalIP", c.localIP)
	req.Header.Set("X-ClientPublicIP", c.publicIP)
	req.Header.Set("X-MACAddress", c.macID)
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")

	// Authentication and metadata headers
	if c.JWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.JWTToken)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var serverResp serverResponse[json.RawMessage]
	// Decode into serverResponse wrapper
	if err := json.Unmarshal(respBody, &serverResp); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s", err, string(respBody))
	}

	if serverResp.ErrorCode != "" {
		return fmt.Errorf("API error: %s (code: %s : %s)", serverResp.Message, serverResp.ErrorCode, getAPIErrorMessage(serverResp.ErrorCode))
	}

	// Check for errors in server response
	if serverResp.Status != "success" {
		return fmt.Errorf("API error: %s (code: %s)", serverResp.Message, serverResp.ErrorCode)
	}

	// If caller doesn't want the response decoded, return
	if result == nil {
		return nil
	}

	// Decode JSON into provided result pointer
	if err := json.Unmarshal(serverResp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s", err, string(respBody))
	}

	return nil
}
