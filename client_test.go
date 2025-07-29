package mistclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	testURL := "https://test.url.com"
	testKey := "xxKEYxx"

	c := New(&Config{
		BaseURL: testURL,
		APIKey:  testKey,
	}, nil)
	if c.config.BaseURL != testURL {
		t.Errorf("NewClient: expected BaseURL: %s, got: %s", testURL, c.config.BaseURL)
	}
	if c.config.APIKey != testKey {
		t.Errorf("NewClient: expected APIKey: %s, got: %s", testKey, c.config.APIKey)
	}
	if c.client.Timeout != time.Duration(10)*time.Second {
		t.Errorf("NewClient: expected Timeout: 10, got: %d", c.client.Timeout)
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
