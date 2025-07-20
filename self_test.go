package mistclient

import (
	"testing"
)

func TestGetSelf(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := New(&Config{BaseURL: s.URL, APIKey: "testAPIKey"}, nil)

	self, err := c.GetSelf()
	if err != nil {
		t.Errorf("APIClient.GetSelf(): Threw error: %s", err)
	}
	if self.Email != "test@mistsys.com" {
		t.Errorf("APIClient.GetSelf(): expected email 'test@mistsys.com', got: %s", self.Email)
	}
	if len(self.Privileges) != 3 {
		t.Errorf("APIClient.GetSelf(): expected 3 privileges, got: %d", len(self.Privileges))
	}
	var (
		orgID  string
		siteID string
		mspID  string
	)
	for _, p := range self.Privileges {
		if p.Scope == "org" {
			orgID = p.OrgID
			continue
		}
		if p.Scope == "site" {
			siteID = p.SiteID
			continue
		}
		if p.Scope == "msp" {
			mspID = p.MSPID
			continue
		}
	}
	if orgID != "9ff00eec-24f0-44d7-bda4-6238c81376ee" {
		t.Errorf("APIClient.GetSelf(): expected orgID '9ff00eec-24f0-44d7-bda4-6238c81376ee', got: %s", orgID)
	}
	if siteID != "d96e3952-53e8-4266-959a-45acd55f5114" {
		t.Errorf("APIClient.GetSelf(): expected siteID 'd96e3952-53e8-4266-959a-45acd55f5114', got: %s", siteID)
	}
	if mspID != "9520c63a-f7b3-670c-0944-727774d5a722" {
		t.Errorf("APIClient.GetSelf(): expected mspID '9520c63a-f7b3-670c-0944-727774d5a722', got: %s", mspID)
	}
}
