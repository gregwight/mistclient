# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-08-07

### Added

-   **Initial Release** of the `mistclient` library.
-   **API Client**:
    -   New client constructor `mistclient.New()` with support for `BaseURL`, `APIKey`, and `Timeout` configuration.
    -   Context-aware methods for request cancellation and timeouts.
    -   Integrated logging using `slog`.
    -   Convenience methods for `GET`, `POST`, `PUT`, `DELETE` HTTP requests.
-   **WebSocket Support**:
    -   Functionality to connect to Mist WebSocket stream endpoints.
    -   `Subscribe` method to listen to specific channels, returning a Go channel for message consumption.
    -   Automatic handling of unsubscribe and connection closing via `context.Context`.
    -   Generic `streamStats` helper for easy implementation of typed data streams.
-   **Typed Data Models**:
    -   Go structs for core Mist API objects: `Self`, `Org`, `Site`, `Device`, `Client`, and their statistics variants (`OrgStat`, `SiteStat`, `DeviceStat`).
    -   Custom types for handling API-specific data formats like `UnixTime` and `Seconds`.
    -   Enum-like types for `DeviceType`, `DeviceStatus`, `Radio`, `Dot11Proto`, etc., with JSON marshalling/unmarshalling.
-   **Supported Endpoints**:
    -   **Self**: `GetSelf()`
    -   **Organization**: `GetOrgSites()`, `CountOrgTickets()`, `CountOrgAlarms()`
    -   **Site (REST)**: `GetSiteStats()`, `GetSiteDevices()`, `GetSiteDeviceStats()`, `GetSiteClientStats()`
    -   **Site (WebSocket)**: `StreamSiteDevices()`, `StreamSiteDeviceStats()`, `StreamSiteClientStats()`
-   **Project Documentation**:
    -   `README.md` with installation instructions, usage examples for both REST and WebSocket APIs, and a list of supported endpoints.
    -   `CONTRIBUTING.md` outlining the development process.
-   **Testing**:
    -   Comprehensive unit tests for the client, REST endpoints, and WebSocket streaming functionality using `httptest`.
