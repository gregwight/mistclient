package mistclient

import (
	"encoding/json"
	"errors"
	"net/netip"
	"time"
)

// Seconds represents a time in seconds
type Seconds time.Duration

// MarshalJSON implements the [json.Marshaler] interface.
func (s *Seconds) MarshalJSON() ([]byte, error) {
	d := (*time.Duration)(s)
	b, err := json.Marshal(d)
	if err != nil {
		return nil, errors.New("Seconds.MarshalJSON: " + err.Error())
	}
	return b, nil
}

func (s *Seconds) UnmarshalJSON(b []byte) error {
	// The Mist API returns this value as a float for some unknown reason...
	var seconds float64
	if err := json.Unmarshal(b, &seconds); err != nil {
		return err
	}
	*s = Seconds(time.Duration(int(seconds)) * time.Second)
	return nil
}

// UnixTime represents the number of seconds since the Eunix epoch
type UnixTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (ut *UnixTime) MarshalJSON() ([]byte, error) {
	if y := ut.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("UnixTime.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(ut.Format(`"` + time.RFC3339 + `"`)), nil
}

func (ut *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

// Self represents the account associated with the authenticated API user
type Self struct {
	Email                string      `json:"email,omitempty"`
	FirstName            string      `json:"first_name,omitempty"`
	LastName             string      `json:"last_name,omitempty"`
	Phone                string      `json:"phone,omitempty"`
	SSO                  bool        `json:"via_sso,omitempty"`
	PasswordLastModified UnixTime    `json:"password_modified_time,omitzero"`
	Privileges           []Privilege `json:"privileges,omitempty"`
	Tags                 []string    `json:"tags,omitempty"`
}

// Privilege represents the permissions of the authenticated API user
type Privilege struct {
	Scope        string   `json:"scope,omitempty"`
	Name         string   `json:"name,omitempty"`
	Role         string   `json:"role,omitempty"`
	Views        []string `json:"views,omitempty"`
	OrgID        string   `json:"org_id,omitempty"`
	OrgName      string   `json:"org_name,omitempty"`
	SiteID       string   `json:"site_id,omitempty"`
	SiteName     string   `json:"site_name,omitempty"`
	MSPID        string   `json:"msp_id,omitempty"`
	MSPName      string   `json:"msp_name,omitempty"`
	MSPURL       string   `json:"msp_url,omitempty"`
	MSPLogoURL   string   `json:"msp_logo_url,omitempty"`
	OrgGroupIDs  []string `json:"orggroup_ids,omitempty"`
	SiteGroupIDs []string `json:"sitegroup_ids,omitempty"`
}

// Alarm represents an alarm created by an error condition
type Alarm struct {
	ID             string   `json:"id,omitzero"`
	Timestamp      UnixTime `json:"timestamp,omitzero"`
	SiteID         string   `json:"site_id,omitzero"`
	Type           string   `json:"type,omitzero"`
	Count          int      `json:"count,omitzero"`
	Acked          bool     `json:"acked,omitzero"`
	AckedTime      UnixTime `json:"acked_time,omitzero"`
	AckedAdminName string   `json:"ack_admin_name,omitzero"`
	AckedAdminID   string   `json:"ack_admin_id,omitzero"`
	Note           string   `json:"note,omitzero"`
}

// Site represents a physical location containing devices
type Site struct {
	ID           string             `json:"id,omitzero"`
	Name         string             `json:"name,omitzero"`
	Timezone     string             `json:"timezone,omitzero"`
	CountryCode  string             `json:"country_code,omitzero"`
	LatLng       map[string]float32 `json:"latlng,omitzero"`
	SiteGroupIDs []string           `json:"sitegroup_ids,omitzero"`
	Address      string             `json:"address,omitzero"`
}

// Device represents a physical piece of network equipment
type Device struct {
	ID           string     `json:"id,omitzero"`
	Name         string     `json:"name,omitzero"`
	Type         DeviceType `json:"type,omitzero"`
	Model        string     `json:"model,omitzero"`
	Serial       string     `json:"serial,omitzero"`
	HwRev        string     `json:"hw_rev,omitzero"`
	Mac          string     `json:"mac,omitzero"`
	OrgID        string     `json:"org_id,omitzero"`
	SiteID       string     `josn:"site_id"`
	ModifiedTime UnixTime   `json:"modified_time,omitzero"`
	CreatedTime  UnixTime   `json:"created_time,omitzero"`
}

// DeviceStat holds operational statistics and data relating to a Device
type DeviceStat struct {
	Device

	LastSeen   UnixTime     `json:"last_seen,omitzero"`
	NumClients int          `json:"num_clients,omitzero"`
	Version    string       `josn:"version"`
	Status     DeviceStatus `json:"status,omitzero"`
	IP         netip.Addr   `json:"ip,omitzero"`
	ExtIP      netip.Addr   `json:"ext_ip,omitzero"`
	NumWLANs   int          `json:"num_wlans,omitzero"`
	Uptime     Seconds      `json:"uptime,omitzero"`
	TxBps      int          `json:"tx_bps,omitzero"`
	RxBps      int          `json:"rx_bps,omitzero"`
	TxBytes    int          `json:"tx_bytes,omitzero"`
	RxBytes    int          `json:"rx_bytes,omitzero"`
	TxPkts     int          `json:"tx_pkts,omitzero"`
	RxPkts     int          `json:"rx_pkts,omitzero"`
	CPUUtil    int          `json:"cpu_util,omitzero"`
	MemUsedKB  int          `json:"mem_used_kb,omitzero"`
	PowerSrc   string       `json:"power_src,omitzero"`

	RadioStats map[Radio]RadioStat `json:"radio_stat,omitzero"`
}

// RadioStat holds operational statistics and data relating to the radio connectivity of a Device
type RadioStat struct {
	Mac                 string `json:"mac,omitzero"`
	NumClients          int    `json:"num_clients,omitzero"`
	NumWLANs            int    `json:"num_wlans,omitzero"`
	Channel             int    `json:"channel,omitzero"`
	Bandwidth           int    `json:"bandwidth,omitzero"`
	Power               int    `json:"power,omitzero"`
	TxBytes             int    `json:"tx_bytes,omitzero"`
	TxPkts              int    `json:"tx_pkts,omitzero"`
	RxBytes             int    `json:"rx_bytes,omitzero"`
	RxPkts              int    `json:"rx_pkts,omitzero"`
	UtilAll             int    `json:"util_all,omitzero"`
	UtilTx              int    `json:"util_tx,omitzero"`
	UtilRxInBSS         int    `json:"util_rx_in_bss,omitzero"`
	UtilRxOtherBSS      int    `json:"util_rx_other_bss,omitzero"`
	UtilUnknownWiFi     int    `json:"util_unknown_wifi,omitzero"`
	UtilNonWiFi         int    `json:"util_non_wifi,omitzero"`
	UtilUndecodableWiFi int    `json:"util_undecodable_wifi,omitzero"`
}

// Client represents an end-user device connected to the radio of a Device
type Client struct {
	Mac         string     `json:"mac,omitzero"`
	LastSeen    UnixTime   `json:"last_seen,omitzero"`
	Username    string     `json:"username,omitzero"`
	Hostname    string     `json:"hostname,omitzero"`
	OS          string     `json:"os,omitzero"`
	Manufacture string     `json:"manufacture,omitzero"`
	Family      string     `json:"family,omitzero"`
	Model       string     `json:"model,omitzero"`
	IP          netip.Addr `json:"ip,omitzero"`
	IP6         netip.Addr `json:"ip6,omitzero"`
	APMac       string     `json:"ap_mac,omitzero"`
	APID        string     `json:"ap_id,omitzero"`
	SSID        string     `json:"ssid,omitzero"`
	WLANID      string     `json:"wlan_id,omitzero"`
	PSKID       string     `json:"psk_id,omitzero"`

	Uptime     Seconds `json:"uptime,omitzero"`
	Idletime   Seconds `json:"idle_time,omitzero"`
	PowerSaing bool    `json:"power_saving,omitzero"`
	Band       string  `jsonn:"band"`
	Proto      string  `json:"proto,omitzero"`
	KeyMgmt    string  `json:"key_mgmt,omitzero"`
	DualBand   bool    `json:"dual_band,omitzero"`

	Channel        int     `json:"channel,omitzero"`
	VLANID         string  `json:"vlan_id,omitzero"`
	AirspaceIfname string  `json:"airespace_ifname,omitzero"`
	RSSI           int     `json:"rssi,omitzero"`
	SNR            int     `json:"snr,omitzero"`
	TxRate         float64 `json:"tx_rate,omitzero"`
	RxRate         float64 `json:"rx_rate,omitzero"`

	TxBytes   int `json:"tx_bytes,omitzero"`
	TxBps     int `json:"tx_bps,omitzero"`
	TxPackets int `json:"tx_packets,omitzero"`
	TxRetries int `json:"tx_retries,omitzero"`
	RxBytes   int `json:"rx_bytes,omitzero"`
	RxBps     int `json:"rx_bps,omitzero"`
	RxPackets int `json:"rx_packets,omitzero"`
	RxRetries int `json:"rx_retries,omitzero"`

	MapID          string  `json:"map_id,omitzero"`
	X              float64 `json:"x,omitzero"`
	Y              float64 `json:"y,omitzero"`
	Xm             float64 `json:"x_m,omitzero"`
	Ym             float64 `json:"y_m,omitzero"`
	NumLocatingAPs int     `json:"num_locating_aps,omitzero"`

	IsGuest  bool     `json:"is_guest,omitzero"`
	Guest    Guest    `json:"guest,omitzero"`
	Airwatch Airwatch `json:"airwatch,omitzero"`
	TTL      int      `json:"_ttl,omitzero"`
}

// Guest holds data relating to the `guest` status of a Client
type Guest struct {
	Authorized             bool     `json:"authorized,omitzero"`
	AuthorizedTime         UnixTime `json:"authorized_time,omitzero"`
	AuthorizedExpiringTime UnixTime `json:"authorized_expiring_time,omitzero"`

	Name      string `json:"name,omitzero"`
	Email     string `json:"email,omitzero"`
	Company   string `json:"company,omitzero"`
	Field1    string `json:"field1,omitzero"`
	CrossSite bool   `json:"cross_site,omitzero"`
}

// Airwatch holds information regarding the 'airwatch` status of a Client
type Airwatch struct {
	Authorized bool `json:"authorized,omitzero"`
}
