package mistclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
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

func TestStreamSiteDeviceStats(t *testing.T) {
	testDeviceStat := DeviceStat{
		Device: Device{
			ID:     "test-device-id",
			Name:   "test-device",
			SiteID: "test-site-id",
		},
		Status: Connected,
	}
	testData, err := json.Marshal(testDeviceStat)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	wsServer := testWebsocketServer(t, false, string(testData))
	defer wsServer.Close()

	wsURL, err := url.Parse(wsServer.URL)
	if err != nil {
		t.Fatalf("failed to parse websocket server URL: %v", err)
	}
	host, port, err := net.SplitHostPort(wsURL.Host)
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}
	testBaseURL := fmt.Sprintf("http://api.%s.nip.io:%s", host, port)

	c, err := New(&Config{BaseURL: testBaseURL, APIKey: "testAPIKey"}, nil)
	if err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}

	siteID := "test-site-id"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statChan, err := c.StreamSiteDeviceStats(ctx, siteID)
	if err != nil {
		t.Fatalf("APIClient.StreamSiteDeviceStats(%s) threw error: %v", siteID, err)
	}

	select {
	case stat, ok := <-statChan:
		if !ok {
			t.Fatalf("APIClient.StreamSiteDeviceStats(%s): channel closed unexpectedly", siteID)
		}
		if stat.ID != testDeviceStat.ID {
			t.Errorf("StreamSiteDeviceStats().ID: expected %q, got %q", testDeviceStat.ID, stat.ID)
		}
		if stat.Status != testDeviceStat.Status {
			t.Errorf("StreamSiteDeviceStats().Status: expected %q, got %q", testDeviceStat.Status, stat.Status)
		}
	case <-ctx.Done():
		t.Fatal("APIClient.StreamSiteDeviceStats(): timed out waiting for stat")
	}

	cancel()
	select {
	case _, ok := <-statChan:
		if ok {
			t.Error("APIClient.StreamSiteDeviceStats(): channel not closed after context cancellation")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("APIClient.StreamSiteDeviceStats(): channel not closed within 1s after context cancellation")
	}
}

func TestStreamSiteDeviceStats_SubFails(t *testing.T) {
	wsServer := testWebsocketServer(t, true)
	defer wsServer.Close()

	wsURL, _ := url.Parse(wsServer.URL)
	host, port, _ := net.SplitHostPort(wsURL.Host)
	testBaseURL := fmt.Sprintf("http://api.%s.nip.io:%s", host, port)

	c, _ := New(&Config{BaseURL: testBaseURL, APIKey: "testAPIKey"}, nil)

	_, err := c.StreamSiteDeviceStats(context.Background(), "test-site-id")
	if err == nil {
		t.Fatal("StreamSiteDeviceStats() expected an error, got nil")
	}
	expectedErr := "websocket subscription failed: subscription_failed"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("StreamSiteDeviceStats() error = %q, want to contain %q", err, expectedErr)
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
