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
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Config represents the API access parameters.
type Config struct {
	BaseURL string        `yaml:"base_url,omitempty"`
	APIKey  string        `yaml:"api_key,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

// APIClient represents the API client.
type APIClient struct {
	baseURL *url.URL
	apiKey  string
	client  *http.Client
	logger  *slog.Logger
}

// SubscriptionRequest represents a websocket subscription request
type SubscriptionRequest struct {
	Subscribe string `json:"subscribe"`
}

// SubscriptionResponse represents a websocket subscription response
type SubscriptionResponse struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

// UnsubscribeRequest represents a websocket unsubscribe request
type UnsubscribeRequest struct {
	Unsubscribe string `json:"unsubscribe"`
}

// WebsocketMessage represents a message streamed by a websocket connection
type WebsocketMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// New returns an instance of the API client.
func New(config *Config, logger *slog.Logger) (*APIClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

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

// GetWebsocketURL maps the base API URL into the appropriate websocket endpoint.
// See https://www.juniper.net/documentation/us/en/software/mist/api/http/guides/websockets/hosts
func (c *APIClient) GetWebsocketURL() (*url.URL, error) {
	u := *c.baseURL

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return nil, fmt.Errorf("unsupported websocket URL scheme: %s", u.Scheme)
	}

	if !strings.HasPrefix(u.Host, "api.") {
		return nil, fmt.Errorf("unable to determin websocket endpoint address, base URL is not prefixed with 'api.': %s", u.Host)
	}

	u.Host = strings.Replace(u.Host, "api.", "api-ws.", 1)
	u.Path = "/api-ws/v1/stream"

	return &u, nil
}

// ConnectWebSocket opens a websocket connection to the appropriate websocket endpoint.
func (c *APIClient) ConnectWebSocket() (*websocket.Conn, error) {
	u, err := c.GetWebsocketURL()
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	wsConfig, err := websocket.NewConfig(u.String(), c.baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket config: %w", err)
	}

	wsConfig.TlsConfig = &tls.Config{
		ServerName: u.Host,
	}

	wsConfig.Dialer = &net.Dialer{
		Timeout: c.client.Timeout,
	}

	wsConfig.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))
	wsConfig.Header.Set("Content-Type", "application/json")

	return websocket.DialConfig(wsConfig)
}

// Subscribe sends a subscription request over a new websocket connection and returns a channel over which received messages will be sent.
func (c *APIClient) Subscribe(ctx context.Context, channel string) (<-chan WebsocketMessage, error) {
	conn, err := c.ConnectWebSocket()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}
	c.logger.Debug("successfully connected to websocket", "url", conn.Config().Location.String())

	if err := websocket.JSON.Send(conn, SubscriptionRequest{Subscribe: channel}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send websocket subscription request: %w", err)
	}
	c.logger.Debug("successfully sent websocket subscription request", "channel", channel)

	var subResp SubscriptionResponse
	conn.SetReadDeadline(time.Now().Add(c.client.Timeout))
	if err := websocket.JSON.Receive(conn, &subResp); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to receive websocket subscription response: %w", err)
	}
	conn.SetReadDeadline(time.Time{})
	c.logger.Debug("successfully received websocket subscription response", "channel", channel)

	if subResp.Event != "channel_subscribed" {
		conn.Close()
		return nil, fmt.Errorf("websocket subscription failed: %s", subResp.Event)
	}
	c.logger.Debug("successfully subscribed to websocket channel", "channel", channel)

	// Spawn a watcher to unsubscribe and close the conn when the context is done
	go func() {
		<-ctx.Done()
		err := c.Unsubscribe(conn, channel)
		if err != nil {
			c.logger.Error("failed to unsubscribe from websocket channel", "error", err, "channel", channel)
		}
		conn.Close()
	}()

	msgChan := make(chan WebsocketMessage)

	go func() {
		defer close(msgChan)

		for {
			var msg WebsocketMessage
			err := websocket.JSON.Receive(conn, &msg)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				c.logger.Error("websocket receive error", "error", err)
				return
			}
			c.logger.Debug("received websocket message", "channel", channel)
			msgChan <- msg
		}
	}()

	return msgChan, nil
}

// Unsubscribe sends an unsubscribe request over an existing websocket connection.
func (c *APIClient) Unsubscribe(conn *websocket.Conn, channel string) error {
	if err := websocket.JSON.Send(conn, UnsubscribeRequest{Unsubscribe: channel}); err != nil {
		return fmt.Errorf("failed to send websocket unsubscribe request: %w", err)
	}
	c.logger.Debug("successfully unsubscribed from websocket channel", "channel", channel)

	return nil
}

// streamStats is a generic helper to subscribe to a websocket channel and stream typed data
func streamStats[T any](ctx context.Context, c *APIClient, channel string) (<-chan T, error) {
	msgChan, err := c.Subscribe(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to websocket channel %s: %w", channel, err)
	}

	statChan := make(chan T)

	go func() {
		defer close(statChan)

		for msg := range msgChan {
			var stat T
			if err := json.Unmarshal([]byte(msg.Data), &stat); err != nil {
				c.logger.Error("failed to unmarshal websocket message", "error", err, "channel", channel)
				continue
			}
			statChan <- stat
		}
	}()

	return statChan, nil
}
