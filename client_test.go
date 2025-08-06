package mistclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestNew(t *testing.T) {
	testURL := "https://test.url.com"
	testKey := "xxKEYxx"

	c, err := New(&Config{
		BaseURL: testURL,
		APIKey:  testKey,
	}, nil)
	if err != nil {
		t.Fatalf("NewClient: unexpected error: %v", err)
	}
	if c.baseURL.Scheme != "https" && c.baseURL.Host != "test.url.com" {
		t.Errorf("NewClient: expected baseURL: %s, got: %s", testURL, c.baseURL.String())
	}
	if c.apiKey != testKey {
		t.Errorf("NewClient: expected APIKey: %s, got: %s", testKey, c.apiKey)
	}
	if c.client.Timeout != time.Duration(10)*time.Second {
		t.Errorf("NewClient: expected Timeout: 10, got: %d", c.client.Timeout)
	}
}

func TestGetWebsocketURL(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		wantURL   string
		expectErr bool
	}{
		{
			name:    "Standard HTTPS URL",
			baseURL: "https://api.mist.com",
			wantURL: "wss://api-ws.mist.com/api-ws/v1/stream",
		},
		{
			name:    "EU HTTPS URL",
			baseURL: "https://api.eu.mist.com",
			wantURL: "wss://api-ws.eu.mist.com/api-ws/v1/stream",
		},
		{
			name:    "HTTP URL",
			baseURL: "http://api.mist.com",
			wantURL: "ws://api-ws.mist.com/api-ws/v1/stream",
		},
		{
			name:      "URL without api prefix",
			baseURL:   "https://mist.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(&Config{BaseURL: tt.baseURL}, nil)
			wsURL, err := c.GetWebsocketURL()
			if (err != nil) != tt.expectErr {
				t.Fatalf("GetWebsocketURL() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && wsURL.String() != tt.wantURL {
				t.Errorf("GetWebsocketURL() = %v, want %v", wsURL.String(), tt.wantURL)
			}
		})
	}
}

func testAPIServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorized := false
		if values, ok := r.Header["Authorization"]; ok {
			for _, value := range values {
				if value == "Token testAPIKey" {
					authorized = true
				}
			}
		}
		if !authorized {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Test API server: missing or invalid authorization token"))
			return
		}
		respDataPath := "testdata" + r.URL.Path
		respData, err := os.ReadFile(respDataPath)
		if err != nil {
			if os.IsNotExist(err) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Test API server: response data not found at: %s", respDataPath)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Test API server: error reading response data: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respData)
	}))
}

func newTestClient(t *testing.T) *APIClient {
	t.Helper()

	s := testAPIServer(t)
	t.Cleanup(s.Close)

	c, err := New(&Config{BaseURL: s.URL, APIKey: "testAPIKey"}, nil)
	if err != nil {
		t.Fatalf("newTestClient: unexpected error: %v", err)
	}

	return c
}

// testWebsocketServer creates a test server that mimics the Mist websocket API.
// It can be configured to send good data or fail the subscription.
func testWebsocketServer(t *testing.T, subShouldFail bool, dataToSend ...string) *httptest.Server {
	t.Helper()
	handler := websocket.Handler(func(ws *websocket.Conn) {
		// The client code closes the connection, so we don't need to defer ws.Close() here.

		// 1. Expect a subscription request
		var subReq SubscriptionRequest
		if err := websocket.JSON.Receive(ws, &subReq); err != nil {
			// This can happen if the client closes the connection early.
			return
		}

		// 2. Send subscription response (success or failure)
		if subShouldFail {
			subResp := SubscriptionResponse{Event: "subscription_failed"}
			if err := websocket.JSON.Send(ws, subResp); err != nil {
				t.Logf("testWebsocketServer: failed to send subscription failure response: %v", err)
			}
			return // End the handler here for failure case.
		}

		subResp := SubscriptionResponse{
			Event:   "channel_subscribed",
			Channel: subReq.Subscribe,
		}
		if err := websocket.JSON.Send(ws, subResp); err != nil {
			t.Logf("testWebsocketServer: failed to send subscription response: %v", err)
			return
		}

		// 3. Send data messages
		for _, data := range dataToSend {
			dataMsg := WebsocketMessage{
				Event:   "data",
				Channel: subReq.Subscribe,
				Data:    data,
			}
			if err := websocket.JSON.Send(ws, dataMsg); err != nil {
				t.Logf("testWebsocketServer: failed to send data message: %v", err)
				return
			}
		}

		// 4. Wait for an unsubscribe request. The client should send this when its context is cancelled.
		var unsubReq UnsubscribeRequest
		// Set a deadline to avoid blocking forever if the client doesn't unsubscribe.
		ws.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err := websocket.JSON.Receive(ws, &unsubReq); err != nil {
			return
		}

		if unsubReq.Unsubscribe != subReq.Subscribe {
			t.Errorf("testWebsocketServer: expected unsubscribe from %q, got %q", subReq.Subscribe, unsubReq.Unsubscribe)
		}
	})
	mux := http.NewServeMux()
	mux.Handle("/api-ws/v1/stream", handler)
	return httptest.NewServer(mux)
}
