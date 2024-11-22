package mistclient

import (
	"testing"
)

func TestGetSiteDevices(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := NewClient(&Config{BaseURL: s.URL, APIKey: "testAPIKey"})

	siteID := "test-site-id"
	devices, err := c.GetSiteDevices(siteID)
	if err != nil {
		t.Errorf("Client.GetSiteDevices(%s): Threw error: %s", siteID, err)
	}
	if len(devices) != 1 {
		t.Errorf("Client.GetSiteDevices(%s): expected 1 device, got: %d", siteID, len(devices))
	} else {
		d := devices[0]
		if d.Mac != "5c5b350e0001" {
			t.Errorf("Client.GetSiteDevices(%s)[0].Mac: expected '5c5b350e0001', got: %s", siteID, d.Mac)
		}
		if d.Type != AP {
			t.Errorf("Client.GetSiteDevices(%s)[0].Type: expected 'ap', got: %s", siteID, d.Type)

		}
	}
}

func TestGetSiteDeviceStats(t *testing.T) {
	s := testAPIServer(t)
	defer s.Close()

	c := NewClient(&Config{BaseURL: s.URL, APIKey: "testAPIKey"})

	siteID := "test-site-id"
	deviceStats, err := c.GetSiteDeviceStats(siteID)
	if err != nil {
		t.Errorf("Client.GetSiteDeviceStats(%s): Threw error: %s", siteID, err)
	}
	if len(deviceStats) != 1 {
		t.Errorf("Client.GetSiteDeviceStats(%s): expected 1 device, got: %d", siteID, len(deviceStats))
	} else {
		ds := deviceStats[0]
		if ds.Status != Connected {
			t.Errorf("Client.GetSiteDeviceStats(%s)[0]: expected status 'connected', got: %s", siteID, ds.Status)
		}
		if ds.RadioStats[Band5].TxBytes != 50877568 {
			t.Errorf("Client.GetSiteDeviceStats(%s)[0].RadioStat[Band5]: expected 50877568 TxBytes, got: %d", siteID, ds.RadioStats[Band5].TxBytes)
		}
	}
}
