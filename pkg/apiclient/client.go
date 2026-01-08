package apiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL      string
	AuthURL      string
	ClientID     string
	ClientSecret string
	Token        string
	HTTPClient   *http.Client
}

func NewClient(clientID string, clientSecret string) *Client {
	return NewClientWithProxy(clientID, clientSecret, "")
}

func NewClientWithProxy(clientID string, clientSecret string, proxyURL string) *Client {
	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			log.Printf("Invalid proxy URL: %v", err)
			return &Client{
				BaseURL:      "https://xmclouddeploy-api.sitecorecloud.io",
				AuthURL:      "https://auth.sitecorecloud.io/oauth/token",
				ClientID:     clientID,
				ClientSecret: clientSecret,
				HTTPClient:   client,
			}
		}
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &Client{
		BaseURL:      "https://xmclouddeploy-api.sitecorecloud.io",
		AuthURL:      "https://auth.sitecorecloud.io/oauth/token",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   client,
	}
}

// doRequest handles the common request logic including authentication
type RequestOptions struct {
	Method string
	Path   string
	Body   interface{}
}

func (c *Client) doRequest(opts RequestOptions) (*http.Response, error) {

	// Ensure we have a valid token
	err := c.EnsureTokenValid()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure valid token: %v", err)
	}

	// Create request URL
	requestURL := fmt.Sprintf("%s%s", c.BaseURL, opts.Path)

	// Create request body if needed
	var reqBody io.Reader
	if opts.Body != nil {
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	req, err := http.NewRequest(opts.Method, requestURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("Failure as status code is %d", resp.StatusCode)
	}

	return resp, nil
}
