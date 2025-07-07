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
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config represents the API access parameters.
type Config struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Timeout int    `yaml:"timeout"`
}

// APIClient represents the API client.
type APIClient struct {
	config *Config
	client *http.Client
}

// NewAPIClient sets a default timeout value and returns an instance of the API client.
func NewAPIClient(config *Config) *APIClient {
	timeout := 10
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	return &APIClient{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// doRequest performs that actual HTTP client request against the provided API endpoint.
// It constructs the URL based on the client configuration and supplied path, sets an
// authentication header based on the API key, and returns the response or any errors.
func (c *APIClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.config.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	log.Debugf("Making %s request to %s", method, url)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Errorf("Request failed: %v", err)
		return nil, err
	}
	log.Debugf("Received response with status code %d", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Unable to read response body: %v", err)
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
func (c *APIClient) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

// Post is a convenience function for performing HTTP POST requests using the API client.
func (c *APIClient) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, body)
}

// Put is a convenience function for performing HTTP PUT requests using the API client.
func (c *APIClient) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", path, body)
}

// Delete is a convenience function for performing HTTP DELETE requests using the API client.
func (c *APIClient) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}
