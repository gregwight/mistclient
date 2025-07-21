package mistclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOrgSites returns a list of all sites configured within an organisation.
func (c *APIClient) GetOrgSites(orgID string) ([]Site, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/sites", orgID)
	resp, err := c.Get(path)
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
	path := fmt.Sprintf("/api/v1/orgs/%s/tickets/count?distinct=status", orgID)
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var result OrgTicketCountResult
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
	path := fmt.Sprintf("/api/v1/orgs/%s/alarms/count?distinct=type", orgID)
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	var OrgAlarmCountResult OrgAlarmCountResult
	if err := json.NewDecoder(resp.Body).Decode(&OrgAlarmCountResult); err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, data := range OrgAlarmCountResult.Results {
		counts[data.Type] = int(data.Count)
	}
	return counts, nil
}
