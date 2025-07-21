# mistclient

[![Go Reference](https://pkg.go.dev/badge/github.com/gregwight/mistclient.svg)](https://pkg.go.dev/github.com/gregwight/mistclient)
<!-- Add other badges as you set them up, e.g., build status, coverage -->

`mistclient` is a Go client library for interacting with the Juniper Mist API.

This package was created with an initial focus on observability, providing access to endpoints useful for building monitoring and exporting tools (e.g., a Prometheus exporter).

## Installation

```sh
go get github.com/gregwight/mistclient
```

## Usage

First, you'll need a Mist API token. You can generate one from the Mist Dashboard under `My Account > API Tokens`.

Here's a simple example of how to create a client and list the sites in your organization:

```go
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gregwight/mistclient"
)

func main() {
	apiKey := os.Getenv("MIST_API_KEY")
	if apiKey == "" {
		log.Fatal("MIST_API_KEY environment variable not set")
	}
	orgID := os.Getenv("MIST_ORG_ID")
	if orgID == "" {
		log.Fatal("MIST_ORG_ID environment variable not set")
	}

	// Create a new client
	client := mistclient.New(&mistclient.Config{
		BaseURL: "https://api.mist.com", // Or your regional cloud, e.g., https://api.eu.mist.com
		APIKey:  apiKey,
	}, slog.Default())

	// Get all sites for the organization
	sites, err := client.GetOrgSites(orgID)
	if err != nil {
		log.Fatalf("Error getting sites: %v", err)
	}

	fmt.Printf("Found %d sites in organization %s:\n", len(sites), orgID)
	for _, site := range sites {
		fmt.Printf("- %s (ID: %s)\n", site.Name, site.ID)
	}
}
```

## Supported API Endpoints

The client currently supports the following endpoints:

*   `GetSelf()` -> `GET /api/v1/self`
*   `GetOrgSites(orgID string)` -> `GET /api/v1/orgs/:org_id/sites`
*   `CountOrgTickets(orgID string)` -> `GET /api/v1/orgs/:org_id/tickets/count`
*   `CountOrgAlarms(orgID string)` -> `GET /api/v1/orgs/:org_id/alarms/count`
*   `GetSiteDevices(siteID string)` -> `GET /api/v1/sites/:site_id/devices`
*   `GetSiteDeviceStats(siteID string)` -> `GET /api/v1/sites/:site_id/stats/devices`
*   `GetSiteClients(siteID string)` -> `GET /api/v1/sites/:site_id/stats/clients`

## Testing

To run the test suite:
```sh
go test -v ./...
```
