package mistclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
