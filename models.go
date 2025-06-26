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

// Client represents an end-user device connected to the radio of a Device
type Client struct {
	Mac         string     `json:"mac"`
	LastSeen    UnixTime   `json:"last_seen"`
	Username    string     `json:"username"`
	Hostname    string     `json:"hostname"`
	OS          string     `json:"os"`
	Manufacture string     `json:"manufacture"`
	Family      string     `json:"family"`
	Model       string     `json:"model"`
	IP          netip.Addr `json:"ip"`
	IP6         netip.Addr `json:"ip6"`
	APMac       string     `json:"ap_mac"`
	APID        string     `json:"ap_id"`
	SSID        string     `json:"ssid"`
	WLANID      string     `json:"wlan_id"`
	PSKID       string     `json:"psk_id"`

	Uptime     Seconds `json:"uptime"`
	Idletime   Seconds `json:"idle_time"`
	PowerSaing bool    `json:"power_saving"`
	Band       string  `jsonn:"band"`
	Proto      string  `json:"proto"`
	KeyMgmt    string  `json:"key_mgmt"`
	DualBand   bool    `json:"dual_band"`

	Channel        int    `json:"channel"`
	VLANID         string `json:"vlan_id"`
	AirspaceIfname string `json:"airespace_ifname"`
	RSSI           int    `json:"rssi"`
	SNR            int    `json:"snr"`
	TxRate         int    `json:"tx_rate"`
	RxRate         int    `json:"rx_rate"`

	TxBytes   int `json:"tx_bytes"`
	TxBps     int `json:"tx_bps"`
	TxPackets int `json:"tx_packets"`
	TxRetries int `json:"tx_retries"`
	RxBytes   int `json:"rx_bytes"`
	RxBps     int `json:"rx_bps"`
	RxPackets int `json:"rx_packets"`
	RxRetries int `json:"rx_retries"`

	MapID          string  `json:"map_id"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	Xm             float64 `json:"x_m"`
	Ym             float64 `json:"y_m"`
	NumLocatingAPs int     `json:"num_locating_aps"`

	IsGuest  bool     `json:"is_guest"`
	Guest    Guest    `json:"guest"`
	Airwatch Airwatch `json:"airwatch"`
	TTL      int      `json:"_ttl"`
}

// Guest holds data relating to the `guest` status of a Client
type Guest struct {
	Authorized             bool     `json:"authorized"`
	AuthorizedTime         UnixTime `json:"authorized_time"`
	AuthorizedExpiringTime UnixTime `json:"authorized_expiring_time"`

	Name      string `json:"name"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Field1    string `json:"field1"`
	CrossSite bool   `json:"cross_site"`
}

// Airwatch holds information regarding the 'airwatch` status of a Client
type Airwatch struct {
	Authorized bool `json:"authorized"`
}
