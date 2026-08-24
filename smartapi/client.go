package smartapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	debug        bool
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
		debug:        false,
		httpClient:   httpClient,
	}
}

type serverResponse[T any] struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      T      `json:"data"`
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// SetProxy configures an optional proxy for this client only.
// Pass an empty string to use a direct connection.
func (c *Client) SetProxy(proxyURL string) error {
	var proxy func(*http.Request) (*url.URL, error)
	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("invalid proxy URL %q", proxyURL)
		}
		proxy = http.ProxyURL(parsed)
	}

	var transport *http.Transport
	switch current := c.httpClient.Transport.(type) {
	case nil:
		transport = http.DefaultTransport.(*http.Transport).Clone()
	case *http.Transport:
		transport = current.Clone()
	default:
		return fmt.Errorf("cannot set proxy on a custom HTTP transport")
	}
	transport.Proxy = proxy
	c.httpClient.Transport = transport
	return nil
}

func (c *Client) doRequest(method string, path string, body any, result any) error {

	// Build URL
	url := c.baseUrl + path

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
		if c.debug {
			fmt.Printf("[DEBUG] Request Body: %s\n", string(b))
		}

	}

	if c.debug {
		fmt.Printf("[DEBUG] Request: %s %s\n", method, url)
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

	if c.debug {
		fmt.Printf("[DEBUG] Request Headers:\n")
		for k, v := range req.Header {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// if c.debug {
	// 	fmt.Printf("[DEBUG] Response Request: %+v\n", resp.Request)
	// }

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response Status: %s\n", resp.Status)
		fmt.Printf("[DEBUG] Response Body: %s\n", string(respBody))
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return newHTTPError(resp.StatusCode, respBody)
	}

	var serverResp serverResponse[json.RawMessage]
	// Decode into serverResponse wrapper
	if err := json.Unmarshal(respBody, &serverResp); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s \n server resp : %v", err, string(respBody), serverResp)
	}

	if serverResp.ErrorCode != "" {
		if serverResp.ErrorCode == "AB1021" {
			return &HTTPError{
				StatusCode: resp.StatusCode,
				Message:    serverResp.Message,
				RawBody:    string(respBody),
			}
		}
		return fmt.Errorf("API error: %s (code: %s : %s)", serverResp.Message, serverResp.ErrorCode, getAPIErrorMessage(serverResp.ErrorCode))
	}

	// Check for errors in server response
	if !serverResp.Status {
		return fmt.Errorf("API error: %s (code: %s)", serverResp.Message, serverResp.ErrorCode)
	}

	// If caller doesn't want the response decoded, return
	if result == nil {
		return nil
	}

	// Decode JSON into provided result pointer
	if err := json.Unmarshal(serverResp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s", err, string(serverResp.Data))
	}

	if c.debug {
		fmt.Printf("[DEBUG] Decoded Result: %+v\n\n\n", result)
	}

	return nil
}

func newHTTPError(statusCode int, body []byte) *HTTPError {
	rawBody := string(body)
	message := string(bytes.TrimSpace(body))

	var serverResp serverResponse[json.RawMessage]
	if json.Unmarshal(body, &serverResp) == nil && serverResp.Message != "" {
		message = serverResp.Message
	}
	if message == "" {
		message = http.StatusText(statusCode)
	}

	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		RawBody:    rawBody,
	}
}
