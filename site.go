package mistclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSiteDevices fetches and returns a list of all devices configured at a site.
func (c *Client) GetSiteDevices(siteID string) ([]Device, error) {
	path := fmt.Sprintf("/api/v1/sites/%s/devices", siteID)
	resp, err := c.Get(path)
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
func (c *Client) GetSiteDeviceStats(siteID string) ([]DeviceStat, error) {
	path := fmt.Sprintf("/api/v1/sites/%s/stats/devices", siteID)
	resp, err := c.Get(path)
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
