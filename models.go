package mistclient

type Alarm struct {
	ID             string `json:"id"`
	Timestamp      int    `json:"timestamp"`
	SiteID         string `json:"site_id"`
	Type           string `json:"type"`
	Count          int    `json:"count"`
	Acked          bool   `json:"acked"`
	AckedTime      int    `json:"acked_time"`
	AckedAdminName string `json:"ack_admin_name"`
	AckedAdminID   string `json:"ack_admin_id"`
	Note           string `json:"note"`
}

type Site struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Timezone     string             `json:"timezone"`
	CountryCode  string             `json:"country_code"`
	LatLng       map[string]float32 `json:"latlng"`
	SiteGroupIDs []string           `json:"sitegroup_ids"`
	Address      string             `json:"address"`
}

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
	ModifiedTime int        `json:"modified_time"`
	CreatedTime  int        `json:"created_time"`
}

type DeviceStat struct {
	Device

	LastSeen   int          `json:"last_seen"`
	NumClients int          `json:"num_clients"`
	Version    string       `josn:"version"`
	Status     DeviceStatus `json:"status"`
	IP         string       `json:"ip"`
	ExtIP      string       `json:"ext_ip"`
	NumWLANs   int          `json:"num_wlans"`
	Uptime     int          `json:"uptime"`
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
