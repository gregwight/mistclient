# mistclient

[![Go Reference](https://pkg.go.dev/badge/github.com/gregwight/mistclient.svg)](https://pkg.go.dev/github.com/gregwight/mistclient)

`mistclient` is a Go client library for interacting with the Juniper Mist API.

This package was created with an initial focus on observability, providing access to endpoints useful for building monitoring and exporting tools (e.g. a Prometheus exporter).

## Features

-   Typed Go models for all supported API responses.
-   Support for both REST and WebSocket streaming endpoints.
-   Context-aware for handling cancellations and timeouts.
-   Configurable logger for integrating with your application's logging.

## Requirements

*   **Go 1.22 or newer is required.** This library uses features of Go's `for` loop semantics that were introduced in version 1.22. Using an older version of Go may lead to subtle concurrency bugs.

## Installation

```sh
go get github.com/gregwight/mistclient
```

## Usage

First, you'll need a Mist API token. You can generate one from the Mist Dashboard under `My Account > API Tokens`. The following examples assume you have set the `MIST_API_KEY` and `MIST_ORG_ID` environment variables.

### Basic API Calls

Here's a simple example of how to create a client and list the sites in an organization using a standard REST API endpoint.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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
	client, err := mistclient.New(&mistclient.Config{
		BaseURL: "https://api.mist.com", // Or your regional cloud, e.g., https://api.eu.mist.com
		APIKey:  apiKey,
	}, slog.Default())
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	
	// Get all sites for the organization
	sites, err := client.GetOrgSites(orgID)
	if err != nil {
		log.Fatalf("Error getting sites: %v", err)
	}

	fmt.Printf("Found %d sites in organization %s:\n", len(sites), orgID)
	for _, site := range sites {
		fmt.Printf("- %s (ID: %s)\n", site.Name, site.ID)
	}

	// Example: Streaming device statistics for a specific site
	// For this to work, you must have at least one site. We'll use the first one found.
	if len(sites) > 0 {
		siteID := sites[0].ID

		// Create a context that will be cancelled after 30 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("\nStreaming device stats for site %s (%s) for 30 seconds...\n", sites[0].Name, siteID)

		// Start streaming device stats
		statsChan, err := client.StreamSiteDeviceStats(ctx, siteID)
		if err != nil {
			log.Fatalf("Error starting device stats stream: %v", err)
		}

		// Process messages from the stream until the context is cancelled
		for stat := range statsChan {
			fmt.Printf("Received update for device %s (%s): Status=%s, Uptime=%s\n",
				stat.Name, stat.Mac, stat.Status, time.Duration(stat.Uptime).String())
		}

		fmt.Println("Stream finished.")
	}
}
```

## Supported API Endpoints

The client is organized by Mist API resources (e.g., Self, Organization, Site).

### Self Endpoints
| Method Signature | API Endpoint |
|---|---|
| `GetSelf() (Self, error)` | `GET /api/v1/self` |

### Organization Endpoints
| Method Signature | API Endpoint |
|---|---|
| `GetOrgSites(orgID string) ([]Site, error)` | `GET /api/v1/orgs/:org_id/sites` |
| `CountOrgTickets(orgID string) (Count, error)` | `GET /api/v1/orgs/:org_id/tickets/count` |
| `CountOrgAlarms(orgID string) (Count, error)` | `GET /api/v1/orgs/:org_id/alarms/count` |

### Site Endpoints
| Method Signature | API Endpoint | Type |
|---|---|---|
| `GetSiteStats(siteID string) (SiteStat, error)` | `GET /api/v1/sites/:site_id/stats` | REST |
| `GetSiteDevices(siteID string) ([]Device, error)` | `GET /api/v1/sites/:site_id/devices` | REST |
| `GetSiteDeviceStats(siteID string) ([]DeviceStat, error)` | `GET /api/v1/sites/:site_id/stats/devices` | REST |
| `GetSiteClientStats(siteID string) ([]Client, error)` | `GET /api/v1/sites/:site_id/stats/clients` | REST |
| `StreamSiteDevices(ctx context.Context, siteID string) (<-chan Device, error)` | `stream /sites/:site_id/devices` | WebSocket |
| `StreamSiteDeviceStats(ctx context.Context, siteID string) (<-chan DeviceStat, error)` | `stream /sites/:site_id/stats/devices` | WebSocket |
| `StreamSiteClientStats(ctx context.Context, siteID string) (<-chan Client, error)` | `stream /sites/:site_id/stats/clients` | WebSocket |

## Testing

To run the test suite:
```sh
go test -v ./...
```
