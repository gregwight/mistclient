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
	case "pending":
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
		return "6"
	case Band5:
		return "5"
	case Band24:
		return "2.4"
	default:
		return "unknown"
	}
}

func RadioFromString(r string) Radio {
	switch r {
	case "6":
		return Band6
	case "5":
		return Band5
	case "24":
		return Band24
	default:
		return 0
	}
}

func (r *Radio) unmarshal(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal radio band from %q: %w", string(data), err)
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

type RadioConfig int

const (
	Band6Config RadioConfig = iota + 1
	Band5Config
	Band24Config
)

func (rc RadioConfig) String() string {
	switch rc {
	case Band6Config:
		return "band_6"
	case Band5Config:
		return "band_5"
	case Band24Config:
		return "band_24"
	default:
		return "unknown"
	}
}

func RadioConfigFromString(r string) RadioConfig {
	switch r {
	case "band_6":
		return Band6Config
	case "band_5":
		return Band5Config
	case "band_24":
		return Band24Config
	default:
		return 0
	}
}

func (rc *RadioConfig) unmarshal(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal radio config band from %q: %w", string(data), err)
	}
	*rc = RadioConfigFromString(s)
	return nil
}

func (rc *RadioConfig) UnmarshalText(b []byte) error {
	return rc.unmarshal(b)
}

func (rc *RadioConfig) UnmarshalJSON(b []byte) error {
	return rc.unmarshal(b)
}

type Dot11Proto int

const (
	A Dot11Proto = iota + 1
	AC
	AX
	B
	G
	N
)

func (dp Dot11Proto) String() string {
	switch dp {
	case A:
		return "a"
	case AC:
		return "ac"
	case AX:
		return "ax"
	case B:
		return "b"
	case G:
		return "g"
	case N:
		return "n"
	default:
		return "unknown"
	}
}

func Dot11ProtoFromString(dp string) Dot11Proto {
	switch dp {
	case "a":
		return A
	case "ac":
		return AC
	case "ax":
		return AX
	case "b":
		return B
	case "g":
		return G
	case "n":
		return N
	default:
		return 0
	}
}

func (dp *Dot11Proto) unmarshal(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal radio config band from %q: %w", string(data), err)
	}
	*dp = Dot11ProtoFromString(s)
	return nil
}

func (dp *Dot11Proto) UnmarshalText(b []byte) error {
	return dp.unmarshal(b)
}

func (dp *Dot11Proto) UnmarshalJSON(b []byte) error {
	return dp.unmarshal(b)
}
