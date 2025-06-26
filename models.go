package mistclient

import (
	"encoding/json"
	"net/netip"
	"time"
)

// Seconds represents a time in seconds
type Seconds time.Duration

func (d *Seconds) UnmarshalJSON(b []byte) error {
	var seconds int64
	if err := json.Unmarshal(b, &seconds); err != nil {
		return err
	}
	*d = Seconds(time.Duration(seconds) * time.Second)
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

// Alarm represents an alarm created by an error condition
type Alarm struct {
	ID             string   `json:"id"`
	Timestamp      UnixTime `json:"timestamp"`
	SiteID         string   `json:"site_id"`
	Type           string   `json:"type"`
	Count          int      `json:"count"`
	Acked          bool     `json:"acked"`
	AckedTime      UnixTime `json:"acked_time"`
	AckedAdminName string   `json:"ack_admin_name"`
	AckedAdminID   string   `json:"ack_admin_id"`
	Note           string   `json:"note"`
}

// Site represents a physical location containing devices
type Site struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Timezone     string             `json:"timezone"`
	CountryCode  string             `json:"country_code"`
	LatLng       map[string]float32 `json:"latlng"`
	SiteGroupIDs []string           `json:"sitegroup_ids"`
	Address      string             `json:"address"`
}

// Device represents a physical piece of network equipment
type Device struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         DeviceType `json:"type"`
	Model        string     `json:"model"`
	Serial       string     `json:"serial"`
	HwRev        string     `json:"hw_rev"`
	Mac          string     `json:"mac"`
	OrgID        string     `json:"org_id"`
	SiteID       string     `josn:"site_id"`
	ModifiedTime UnixTime   `json:"modified_time"`
	CreatedTime  UnixTime   `json:"created_time"`
}

// DeviceStat holds operational statistics and data relating to a Device
type DeviceStat struct {
	Device

	LastSeen   UnixTime     `json:"last_seen"`
	NumClients int          `json:"num_clients"`
	Version    string       `josn:"version"`
	Status     DeviceStatus `json:"status"`
	IP         netip.Addr   `json:"ip"`
	ExtIP      netip.Addr   `json:"ext_ip"`
	NumWLANs   int          `json:"num_wlans"`
	Uptime     Seconds      `json:"uptime"`
	TxBps      int          `json:"tx_bps"`
	RxBps      int          `json:"rx_bps"`
	TxBytes    int          `json:"tx_bytes"`
	RxBytes    int          `json:"rx_bytes"`
	TxPkts     int          `json:"tx_pkts"`
	RxPkts     int          `json:"rx_pkts"`
	CPUUtil    int          `json:"cpu_util"`
	MemUsedKB  int          `json:"mem_used_kb"`
	PowerSrc   string       `json:"power_src"`

	RadioStats map[Radio]RadioStat `json:"radio_stat"`
}

// RadioStat holds operational statistics and data relating to the radio connectivity of a Device
type RadioStat struct {
	Mac                 string `json:"mac"`
	NumClients          int    `json:"num_clients"`
	NumWLANs            int    `json:"num_wlans"`
	Channel             int    `json:"channel"`
	Bandwidth           int    `json:"bandwidth"`
	Power               int    `json:"power"`
	TxBytes             int    `json:"tx_bytes"`
	TxPkts              int    `json:"tx_pkts"`
	RxBytes             int    `json:"rx_bytes"`
	RxPkts              int    `json:"rx_pkts"`
	UtilAll             int    `json:"util_all"`
	UtilTx              int    `json:"util_tx"`
	UtilRxInBSS         int    `json:"util_rx_in_bss"`
	UtilRxOtherBSS      int    `json:"util_rx_other_bss"`
	UtilUnknownWiFi     int    `json:"util_unknown_wifi"`
	UtilNonWiFi         int    `json:"util_non_wifi"`
	UtilUndecodableWiFi int    `json:"util_undecodable_wifi"`
}
