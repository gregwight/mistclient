package mistclient

import (
	"encoding/json"
	"net/http"
)

// GetSelf returns a ‘whoami’ and privileges of the account making the request
func (c *APIClient) GetSelf() (Self, error) {
	var self Self

	resp, err := c.Get(c.baseURL.JoinPath("/api/v1/self"))
	if err != nil {
		return self, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return self, extractError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&self); err != nil {
		return self, err
	}

	return self, nil
}
