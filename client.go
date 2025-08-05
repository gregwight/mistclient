// Package mistclient provides a client for interacting Juniper's MIST management API.
//
// It has been written with an initial focus on observability, with the currently supported
// endpoints chosen for integration with a Prometheus exporter to produce operational metrics.
//
// The client currently supports the following Organization endpoints:
//   - /api/v1/orgs/:org_id/sites
//   - /api/v1/orgs/:org_id/tickets/count
//   - /api/v1/orgs/:org_id/alarms/count
//
// The client currently supports the following Site endpoints:
//   - /api/v1/sites/:site_id/devices
//   - /api/v1/sites/:site_id/stats/devices
//   - /api/v1/sites/:site_id/stats/clients
package mistclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Config represents the API access parameters.
type Config struct {
	BaseURL string        `yaml:"base_url"`
	APIKey  string        `yaml:"api_key"`
	Timeout time.Duration `yaml:"timeout"`
}

// APIClient represents the API client.
type APIClient struct {
	baseURL *url.URL
	apiKey  string
	client  *http.Client
	logger  *slog.Logger
}
}

// New returns an instance of the API client.
func New(config *Config, logger *slog.Logger) (*APIClient, error) {
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &APIClient{
		baseURL: baseURL,
		apiKey:  config.APIKey,
		logger:  logger.With("module", "mistclient"),
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// doRequest performs that actual HTTP client request against the provided API endpoint.
// It constructs the URL based on the client configuration and supplied path, sets an
// authentication header based on the API key, and returns the response or any errors.
func (c *APIClient) doRequest(method string, u *url.URL, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	c.logger.Debug("making API request", "method", method, "url", u.String())
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("request failed", "error", err)
		return nil, err
	}
	c.logger.Debug("API response received", "status", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Debug("unable to read response body", "error", err)
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	return resp, nil
}

// extractError is a convenience method for decoding the response body and returning it as an error.
// It is typically called when a returned HTTP status code does not match the expected value.
func extractError(r *http.Response) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("API request failed with status %d and error reading body: %v", r.StatusCode, err)
	}
	return fmt.Errorf("API request failed with status %d: %s", r.StatusCode, string(body))
}

// Get is a convenience function for performing HTTP GET requests using the API client.
func (c *APIClient) Get(u *url.URL) (*http.Response, error) {
	return c.doRequest("GET", u, nil)
}

// Post is a convenience function for performing HTTP POST requests using the API client.
func (c *APIClient) Post(u *url.URL, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", u, body)
}

// Put is a convenience function for performing HTTP PUT requests using the API client.
func (c *APIClient) Put(u *url.URL, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", u, body)
}

// Delete is a convenience function for performing HTTP DELETE requests using the API client.
func (c *APIClient) Delete(u *url.URL) (*http.Response, error) {
	return c.doRequest("DELETE", u, nil)
}

}
