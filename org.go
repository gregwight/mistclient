package mistclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOrgSites returns a list of all sites configured within an organisation.
func (c *APIClient) GetOrgSites(orgID string) ([]Site, error) {
	resp, err := c.Get(c.baseURL.JoinPath(fmt.Sprintf("/api/v1/orgs/%s/sites", orgID)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var sites []Site
	if err := json.NewDecoder(resp.Body).Decode(&sites); err != nil {
		return nil, err
	}

	return sites, nil
}

// CountOrgTickets returns a map of counts of all tickets related to an organisation, keyed by their status.
func (c *APIClient) CountOrgTickets(orgID string) (map[TicketStatus]int, error) {
	u := c.baseURL.JoinPath(fmt.Sprintf("/api/v1/orgs/%s/tickets/count", orgID))

	q := u.Query()
	q.Add("distinct", "status")

	u.RawQuery = q.Encode()

	resp, err := c.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	result := struct {
		Results []struct {
			Status string  `json:"status"`
			Count  float64 `json:"count"`
		} `json:"results"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	counts := make(map[TicketStatus]int)
	for _, data := range result.Results {
		counts[TicketStatusFromString(data.Status)] = int(data.Count)
	}
	return counts, nil
}

// CountOrgAlarms returns a map of counts of all alarms related to an organisation, keyed by their type.
func (c *APIClient) CountOrgAlarms(orgID string) (map[string]int, error) {
	u := c.baseURL.JoinPath(fmt.Sprintf("/api/v1/orgs/%s/alarms/count", orgID))

	q := u.Query()
	q.Add("distinct", "type")

	u.RawQuery = q.Encode()

	resp, err := c.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	result := struct {
		Results []struct {
			Type  string  `json:"type"`
			Count float64 `json:"count"`
		} `json:"results"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, data := range result.Results {
		counts[data.Type] = int(data.Count)
	}
	return counts, nil
}

// ListOrgDevices returns a map of device MAC addresses to names.
func (c *APIClient) ListOrgDevices(orgID string) (map[string]string, error) {
	resp, err := c.Get(c.baseURL.JoinPath(fmt.Sprintf("/api/v1/orgs/%s/devices", orgID)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	// Decode array of device objects, then build the map
	var list []struct {
		Mac  string `json:"mac"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	devices := make(map[string]string, len(list))
	for _, d := range list {
		devices[d.Mac] = d.Name
	}

	return devices, nil
}
