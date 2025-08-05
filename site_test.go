package mistclient

import (
	"testing"
	"time"
)

func TestGetSiteDevices(t *testing.T) {
	c := newTestClient(t)

	siteID := "test-site-id"
	devices, err := c.GetSiteDevices(siteID)
	if err != nil {
		t.Errorf("APIClient.GetSiteDevices(%s): Threw error: %s", siteID, err)
	}
	if len(devices) != 1 {
		t.Errorf("APIClient.GetSiteDevices(%s): expected 1 device, got: %d", siteID, len(devices))
	} else {
		d := devices[0]
		if d.Mac != "5c5b350e0001" {
			t.Errorf("APIClient.GetSiteDevices(%s)[0].Mac: expected '5c5b350e0001', got: %s", siteID, d.Mac)
		}
		if d.Type != AP {
			t.Errorf("APIClient.GetSiteDevices(%s)[0].Type: expected 'ap', got: %s", siteID, d.Type)

		}
	}
}

func TestGetSiteDeviceStats(t *testing.T) {
	c := newTestClient(t)

	siteID := "test-site-id"
	deviceStats, err := c.GetSiteDeviceStats(siteID)
	if err != nil {
		t.Errorf("APIClient.GetSiteDeviceStats(%s): Threw error: %s", siteID, err)
	}
	if len(deviceStats) != 1 {
		t.Errorf("APIClient.GetSiteDeviceStats(%s): expected 1 device, got: %d", siteID, len(deviceStats))
	} else {
		ds := deviceStats[0]
		if ds.Status != Connected {
			t.Errorf("APIClient.GetSiteDeviceStats(%s)[0]: expected status 'connected', got: %s", siteID, ds.Status)
		}
		if ds.RadioStats[Band5Config].TxBytes != 50877568 {
			t.Errorf("APIClient.GetSiteDeviceStats(%s)[0].RadioStat[Band5Config]: expected 50877568 TxBytes, got: %d", siteID, ds.RadioStats[Band5Config].TxBytes)
		}
	}
}

func TestGetSiteClients(t *testing.T) {
	c := newTestClient(t)

	siteID := "test-site-id"
	clients, err := c.GetSiteClientStats(siteID)
	if err != nil {
		t.Errorf("APIClient.GetSiteClients(%s): Threw error: %s", siteID, err)
	}
	if len(clients) != 1 {
		t.Errorf("APIClient.GetSiteClients(%s): expected 1 client, got: %d", siteID, len(clients))
	} else {
		client := clients[0]
		if client.Mac != "5684dae9ac8b" {
			t.Errorf("APIClient.GetSiteClients(%s)[0].Mac: expected 5684dae9ac8b, got: %s", siteID, client.Mac)
		}
		if client.Uptime != Seconds(time.Duration(3568)*time.Second) {
			t.Errorf("APIClient.GetSiteClients(%s)[0].Uptime: expected 3568s, got: %d", siteID, client.Uptime)
		}
		if !client.Guest.Authorized {
			t.Errorf("APIClient.GetSiteClients(%s)[0].Guest.Authorized: expected true, got: %t", siteID, client.Guest.Authorized)
		}
		if client.Band != Band24 {
			t.Errorf("APIClient.GetSiteClients(%s)[0].Band: expected '5', got: %s", siteID, client.Band)
		}
	}
}
