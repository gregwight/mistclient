package mistclient

import (
	"encoding/json"
	"fmt"
)

type TicketStatus int

const (
	Open TicketStatus = iota + 1
	Pending
	Solved
	Closed
	Hold
)

func (ts TicketStatus) String() string {
	switch ts {
	case Open:
		return "open"
	case Pending:
		return "pending"
	case Solved:
		return "solved"
	case Closed:
		return "closed"
	case Hold:
		return "hold"
	default:
		return "unknown"
	}
}

func TicketStatusFromString(status string) TicketStatus {
	switch status {
	case "open":
		return Open
	case "pendng":
		return Pending
	case "solved":
		return Solved
	case "closed":
		return Closed
	case "hold":
		return Hold
	default:
		return 0
	}
}

type DeviceType int

const (
	AP DeviceType = iota + 1
	Switch
	Gateway
)

func (dt DeviceType) String() string {
	switch dt {
	case AP:
		return "ap"
	case Switch:
		return "switch"
	case Gateway:
		return "gateway"
	default:
		return "unknown"
	}
}

func DeviceTypeFromString(dt string) DeviceType {
	switch dt {
	case "ap":
		return AP
	case "switch":
		return Switch
	case "gateway":
		return Gateway
	default:
		return 0
	}
}

func (dt *DeviceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*dt = DeviceTypeFromString(s)
	return nil
}

type DeviceStatus int

const (
	Connected DeviceStatus = iota + 1
	Disconnected
	Restarting
	Upgrading
)

func (ds DeviceStatus) String() string {
	switch ds {
	case Connected:
		return "connected"
	case Disconnected:
		return "disconnected"
	case Restarting:
		return "restarting"
	case Upgrading:
		return "upgrading"
	default:
		return "unknown"
	}
}

func DeviceStatusFromString(ds string) DeviceStatus {
	switch ds {
	case "connected":
		return Connected
	case "disconnected":
		return Disconnected
	case "restarting":
		return Restarting
	case "upgrading":
		return Upgrading
	default:
		return 0
	}
}

func (ds *DeviceStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ds = DeviceStatusFromString(s)
	return nil
}

type Radio int

const (
	Band6 Radio = iota + 1
	Band5
	Band24
)

func (r Radio) String() string {
	switch r {
	case Band6:
		return "band_6"
	case Band5:
		return "band_5"
	case Band24:
		return "band_24"
	default:
		return "unknown"
	}
}

func RadioFromString(r string) Radio {
	switch r {
	case "band_6":
		return Band6
	case "band_5":
		return Band5
	case "band_24":
		return Band24
	default:
		return 0
	}
}

func (r *Radio) unmarshal(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("JSON error: %+q", data)
	}
	*r = RadioFromString(s)
	return nil
}

func (r *Radio) UnmarshalText(b []byte) error {
	return r.unmarshal(b)
}

func (r *Radio) UnmarshalJSON(b []byte) error {
	return r.unmarshal(b)
}
