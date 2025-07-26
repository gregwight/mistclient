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
func (s Seconds) MarshalJSON() ([]byte, error) {
	d := (time.Duration)(s)
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
	ID             string   `json:"id,omitempty"`
	Timestamp      UnixTime `json:"timestamp,omitzero"`
	SiteID         string   `json:"site_id,omitempty"`
	Type           string   `json:"type,omitempty"`
	Count          int      `json:"count,omitempty"`
	Acked          bool     `json:"acked,omitempty"`
	AckedTime      UnixTime `json:"acked_time,omitzero"`
	AckedAdminName string   `json:"ack_admin_name,omitempty"`
	AckedAdminID   string   `json:"ack_admin_id,omitempty"`
	Note           string   `json:"note,omitempty"`
}

// Site represents a physical location containing devices
type Site struct {
	ID           string             `json:"id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Timezone     string             `json:"timezone,omitempty"`
	CountryCode  string             `json:"country_code,omitempty"`
	LatLng       map[string]float32 `json:"latlng,omitempty"`
	SiteGroupIDs []string           `json:"sitegroup_ids,omitempty"`
	Address      string             `json:"address,omitempty"`
}

// Device represents a physical piece of network equipment
type Device struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	Type         DeviceType `json:"type,omitempty"`
	Model        string     `json:"model,omitempty"`
	Serial       string     `json:"serial,omitempty"`
	HwRev        string     `json:"hw_rev,omitempty"`
	Mac          string     `json:"mac,omitempty"`
	OrgID        string     `json:"org_id,omitempty"`
	SiteID       string     `json:"site_id,omitempty"`
	ModifiedTime UnixTime   `json:"modified_time,omitzero"`
	CreatedTime  UnixTime   `json:"created_time,omitzero"`
}

// DeviceStat holds operational statistics and data relating to a Device
type DeviceStat struct {
	Device

	LastSeen   UnixTime     `json:"last_seen,omitzero"`
	NumClients int          `json:"num_clients,omitempty"`
	Version    string       `json:"version,omitempty"`
	Status     DeviceStatus `json:"status,omitempty"`
	IP         netip.Addr   `json:"ip,omitempty"`
	ExtIP      netip.Addr   `json:"ext_ip,omitempty"`
	NumWLANs   int          `json:"num_wlans,omitempty"`
	Uptime     Seconds      `json:"uptime,omitempty"`
	TxBps      int          `json:"tx_bps,omitempty"`
	RxBps      int          `json:"rx_bps,omitempty"`
	TxBytes    int          `json:"tx_bytes,omitempty"`
	RxBytes    int          `json:"rx_bytes,omitempty"`
	TxPkts     int          `json:"tx_pkts,omitempty"`
	RxPkts     int          `json:"rx_pkts,omitempty"`
	PowerSrc   string       `json:"power_src,omitempty"`

	RadioStats map[RadioConfig]RadioStat `json:"radio_stat,omitempty"`
}

// RadioStat holds operational statistics and data relating to the radio connectivity of a Device
type RadioStat struct {
	Mac                 string `json:"mac,omitempty"`
	NumClients          int    `json:"num_clients,omitempty"`
	NumWLANs            int    `json:"num_wlans,omitempty"`
	Channel             int    `json:"channel,omitempty"`
	Bandwidth           int    `json:"bandwidth,omitempty"`
	Power               int    `json:"power,omitempty"`
	TxBytes             int    `json:"tx_bytes,omitempty"`
	TxPkts              int    `json:"tx_pkts,omitempty"`
	RxBytes             int    `json:"rx_bytes,omitempty"`
	RxPkts              int    `json:"rx_pkts,omitempty"`
	UtilAll             int    `json:"util_all,omitempty"`
	UtilTx              int    `json:"util_tx,omitempty"`
	UtilRxInBSS         int    `json:"util_rx_in_bss,omitempty"`
	UtilRxOtherBSS      int    `json:"util_rx_other_bss,omitempty"`
	UtilUnknownWiFi     int    `json:"util_unknown_wifi,omitempty"`
	UtilNonWiFi         int    `json:"util_non_wifi,omitempty"`
	UtilUndecodableWiFi int    `json:"util_undecodable_wifi,omitempty"`
}

// Client represents an end-user device connected to the radio of a Device
type Client struct {
	Mac         string     `json:"mac,omitempty"`
	LastSeen    UnixTime   `json:"last_seen,omitzero"`
	Username    string     `json:"username,omitempty"`
	Hostname    string     `json:"hostname,omitempty"`
	OS          string     `json:"os,omitempty"`
	Manufacture string     `json:"manufacture,omitempty"`
	Family      string     `json:"family,omitempty"`
	Model       string     `json:"model,omitempty"`
	IP          netip.Addr `json:"ip,omitempty"`
	IP6         netip.Addr `json:"ip6,omitempty"`
	APMac       string     `json:"ap_mac,omitempty"`
	APID        string     `json:"ap_id,omitempty"`
	SSID        string     `json:"ssid,omitempty"`
	WLANID      string     `json:"wlan_id,omitempty"`
	PSKID       string     `json:"psk_id,omitempty"`

	Uptime     Seconds `json:"uptime,omitempty"`
	Idletime   Seconds `json:"idle_time,omitempty"`
	PowerSaing bool    `json:"power_saving,omitempty"`
	Band       Radio   `json:"band,omitempty"`
	Proto      string  `json:"proto,omitempty"`
	KeyMgmt    string  `json:"key_mgmt,omitempty"`
	DualBand   bool    `json:"dual_band,omitempty"`

	Channel        int     `json:"channel,omitempty"`
	VLANID         string  `json:"vlan_id,omitempty"`
	AirspaceIfname string  `json:"airespace_ifname,omitempty"`
	RSSI           int     `json:"rssi,omitempty"`
	SNR            int     `json:"snr,omitempty"`
	TxRate         float64 `json:"tx_rate,omitempty"`
	RxRate         float64 `json:"rx_rate,omitempty"`

	TxBytes   int `json:"tx_bytes,omitempty"`
	TxBps     int `json:"tx_bps,omitempty"`
	TxPackets int `json:"tx_packets,omitempty"`
	TxRetries int `json:"tx_retries,omitempty"`
	RxBytes   int `json:"rx_bytes,omitempty"`
	RxBps     int `json:"rx_bps,omitempty"`
	RxPackets int `json:"rx_packets,omitempty"`
	RxRetries int `json:"rx_retries,omitempty"`

	MapID          string  `json:"map_id,omitempty"`
	X              float64 `json:"x,omitempty"`
	Y              float64 `json:"y,omitempty"`
	Xm             float64 `json:"x_m,omitempty"`
	Ym             float64 `json:"y_m,omitempty"`
	NumLocatingAPs int     `json:"num_locating_aps,omitempty"`

	IsGuest  bool     `json:"is_guest,omitempty"`
	Guest    Guest    `json:"guest,omitempty"`
	Airwatch Airwatch `json:"airwatch,omitempty"`
	TTL      int      `json:"_ttl,omitempty"`
}

// Guest holds data relating to the `guest` status of a Client
type Guest struct {
	Authorized             bool     `json:"authorized,omitempty"`
	AuthorizedTime         UnixTime `json:"authorized_time,omitzero"`
	AuthorizedExpiringTime UnixTime `json:"authorized_expiring_time,omitzero"`

	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Company   string `json:"company,omitempty"`
	Field1    string `json:"field1,omitempty"`
	CrossSite bool   `json:"cross_site,omitempty"`
}

// Airwatch holds information regarding the 'airwatch` status of a Client
type Airwatch struct {
	Authorized bool `json:"authorized,omitempty"`
}
