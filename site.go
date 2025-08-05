package mistclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSiteDevices fetches and returns a list of all devices configured at a site.
func (c *APIClient) GetSiteDevices(siteID string) ([]Device, error) {
	resp, err := c.Get(c.baseURL.JoinPath(fmt.Sprintf("/api/v1/sites/%s/devices", siteID)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// GetSiteDeviceStats fetches and returns a list of all devices configured at a site, supplemented with operational statistics.
func (c *APIClient) GetSiteDeviceStats(siteID string) ([]DeviceStat, error) {
	resp, err := c.Get(c.baseURL.JoinPath(fmt.Sprintf("/api/v1/sites/%s/stats/devices", siteID)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var devices []DeviceStat
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// StreamSiteDeviceStats opens a websocket connection and subscribes to the device statistics stream
func (c *APIClient) StreamSiteDeviceStats(ctx context.Context, siteID string) (<-chan DeviceStat, error) {
	return streamSiteStats[DeviceStat](ctx, c, fmt.Sprintf("/sites/%s/stats/devices", siteID))

}

// GetSiteClientStats fetches and returns a list of all clients configured at a site
func (c *APIClient) GetSiteClientStats(siteID string) ([]Client, error) {
	resp, err := c.Get(c.baseURL.JoinPath(fmt.Sprintf("/api/v1/sites/%s/stats/clients", siteID)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var clients []Client
	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return nil, err
	}

	return clients, nil
}

// StreamSiteClientStats opens a websocket connection and subscribes to the client statistics stream
func (c *APIClient) StreamSiteClientStats(ctx context.Context, siteID string) (<-chan Client, error) {
	return streamSiteStats[Client](ctx, c, fmt.Sprintf("/sites/%s/stats/clients", siteID))
}

// streamSiteStats is a generic helper to subscribe to a websocket channel and stream typed data
func streamSiteStats[T any](ctx context.Context, c *APIClient, channel string) (<-chan T, error) {
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
