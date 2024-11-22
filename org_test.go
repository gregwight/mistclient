package mistclient

import (
	"testing"
)

func TestGetOrgSites(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := NewClient(&Config{BaseURL: s.URL, APIKey: "testAPIKey"})

	orgID, siteID := "test-org-id", "4ac1dcf4-9d8b-7211-65c4-057819f0862b"
	sites, err := c.GetOrgSites(orgID)
	if err != nil {
		t.Errorf("Client.GetOrgSites(%s): Threw error: %s", orgID, err)
	}
	if len(sites) != 1 {
		t.Errorf("Client.GetOrgSites(%s): expected sites 1, got: %d", orgID, len(sites))
	} else if sites[0].ID != siteID {
		t.Errorf("Client.GetOrgSites(%s): expected site ID %s, got: %s", orgID, siteID, sites[0].ID)
	}

	orgID = "randon-org-id"
	sites, err = c.GetOrgSites(orgID)
	if err == nil {
		t.Errorf("Client.GetOrgSites(%s): Did not throw expected error.", orgID)
	}
	if sites != nil {
		t.Errorf("Client.GetOrgSites(%s): Returned %d unexpected sites.", orgID, len(sites))
	}
}

func TestCountOrgTickets(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := NewClient(&Config{BaseURL: s.URL, APIKey: "testAPIKey"})

	orgID := "test-org-id"
	ticketCounts, err := c.CountOrgTickets(orgID)
	if err != nil {
		t.Errorf("Client.CountOrgTickets(%s): Threw error: %s", orgID, err)
	}
	if len(ticketCounts) != 5 {
		t.Errorf("Client.CountOrgTickets(%s): expected 5 ticket types, got: %d", orgID, len(ticketCounts))
	} else if ticketCounts[Open] != 12 {
		t.Errorf("Client.CountOrgTickets(%s): expected 12 Open tickets, got: %d", orgID, ticketCounts[Open])
	}
}

func TestCountOrgAlarms(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := NewClient(&Config{BaseURL: s.URL, APIKey: "testAPIKey"})

	orgID := "test-org-id"
	alarmCounts, err := c.CountOrgAlarms(orgID)
	if err != nil {
		t.Errorf("Client.CountOrgAlarms(%s): Threw error: %s", orgID, err)
	}
	if len(alarmCounts) != 3 {
		t.Errorf("Client.CountOrgAlarms(%s): expected 3 alarm types, got: %d", orgID, len(alarmCounts))
	} else if alarmCounts["device_down"] != 1 {
		t.Errorf("Client.CountOrgAlarms(%s): expected 1 'device_down' alarm, got: %d", orgID, alarmCounts["device_down"])
	}
}
