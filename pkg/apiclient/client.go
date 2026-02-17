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
	"os"
	"strings"
)

type Client struct {
	BaseURL      string
	AuthURL      string
	ClientID     string
	ClientSecret string
	CliConfig    *CLIUserConfig
	Token        string
	HTTPClient   *http.Client
}

// NewClientFromCLI attempts to create a client using CLI authentication
func NewClientFromCLI(configPath string) (*Client, error) {

	if configPath == "" {
		foundConfigPath, err := findCLIUserConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %v", err)
		}

		configPath = foundConfigPath
	}

	return NewClientWithAllConfig("", "", "", "", configPath, &http.Client{})
}

func NewClientFromEnv() (*Client, error) {
	clientID := os.Getenv("SITECOREAI_CLIENT_ID")
	clientSecret := os.Getenv("SITECOREAI_CLIENT_SECRET")

	return NewClientWithAllConfig("", "", clientID, clientSecret, "", &http.Client{})
}

func NewClient(clientID string, clientSecret string) (*Client, error) {
	return NewClientWithAllConfig("", "", clientID, clientSecret, "", &http.Client{})
}

func NewClientWithAllConfig(baseUrl string, authUrl string, clientId string, clientSecret string, cliUserConfigPath string, httpClient *http.Client) (*Client, error) {

	BaseURL := "https://xmclouddeploy-api.sitecorecloud.io"
	AuthURL := "https://auth.sitecorecloud.io"

	if len(baseUrl) > 0 {
		BaseURL = baseUrl
	}

	if len(authUrl) > 0 {
		AuthURL = authUrl
	}

	var cliConfig *CLIUserConfig

	if cliUserConfigPath != "" {
		cfg, err := readCLIUserConfig(cliUserConfigPath)
		if cfg == nil || err != nil {
			return nil, fmt.Errorf("failed to read specified cli config from %s: %v", cliUserConfigPath, err)
		}

		cliConfig = cfg

		// If urls are not explicitly overriden, then use values from config
		if len(baseUrl) == 0 {
			BaseURL = cliConfig.Endpoints.XMCloud.Host
		}

		if len(authUrl) == 0 {
			AuthURL = cliConfig.Endpoints.XMCloud.Authority
		}
	}

	if cliConfig == nil && (len(clientId) == 0 || len(clientSecret) == 0) {
		return nil, fmt.Errorf("client_id and client_secret must be provided")
	}

	setupProxy(httpClient)
	return &Client{
		BaseURL:      strings.TrimSuffix(BaseURL, "/"),
		AuthURL:      strings.TrimSuffix(AuthURL, "/"),
		ClientID:     clientId,
		ClientSecret: clientSecret,
		CliConfig:    cliConfig,
		HTTPClient:   httpClient,
	}, nil
}

func setupProxy(client *http.Client) {
	proxy := os.Getenv("SITECOREAI_PROXY")

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Printf("Invalid proxy URL: %v", err)
		}
		log.Printf("Using insecure proxy for %s", proxy)
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
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
		return resp, fmt.Errorf("failure as status code is %d", resp.StatusCode)
	}

	return resp, nil
}
