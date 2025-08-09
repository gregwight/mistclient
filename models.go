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
	var seconds float32
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

// MarshalJSON implements the [json.Marshaler] interface.
func (ut UnixTime) MarshalJSON() ([]byte, error) {
	s := ut.Unix()
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.New("UnixTime.MarshalJSON: " + err.Error())
	}
	return b, nil
}

func (ut *UnixTime) UnmarshalJSON(b []byte) error {
	var dirtyTimestamp float32
	if err := json.Unmarshal(b, &dirtyTimestamp); err != nil {
		return err
	}
	timestamp := int64(dirtyTimestamp)
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
	ID                string             `json:"id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Timezone          string             `json:"timezone,omitempty"`
	CountryCode       string             `json:"country_code,omitempty"`
	LatLng            map[string]float32 `json:"latlng,omitempty"`
	SiteGroupIDs      []string           `json:"sitegroup_ids,omitempty"`
	Address           string             `json:"address,omitempty"`
	OrgID             string             `json:"org_id,omitempty"`
	ModifiedTime      UnixTime           `json:"modified_time,omitzero"`
	CreatedTime       UnixTime           `json:"created_time,omitzero"`
	Notes             string             `json:"notes,omitempty"`
	AlarmTemplateID   string             `json:"alarmtemplate_id,omitempty"`
	NetworkTemplateID string             `json:"networktemplate_id,omitempty"`
	RFTemplateID      string             `json:"rftemplate_id,omitempty"`
	APTemplateID      string             `json:"aptemplate_id,omitempty"`
	GatewayTemplateID string             `json:"gatewaytemplate_id,omitempty"`
	ApportTemplateID  string             `json:"apporttemplate_id,omitempty"`
	SecPolicyID       string             `json:"secpolicy_id,omitempty"`
}

// Org represents an organization
type Org struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	MSPID           string   `json:"msp_id,omitempty"`
	AllowMist       bool     `json:"allow_mist,omitempty"`
	AlarmTemplateID string   `json:"alarmtemplate_id,omitempty"`
	OrgGroupIDs     []string `json:"orggroup_ids,omitempty"`
	SessionExpiry   Seconds  `json:"session_expiry,omitempty"`
	ModifiedTime    UnixTime `json:"modified_time,omitzero"`
	CreatedTime     UnixTime `json:"created_time,omitzero"`
}

// OrgStat holds operational statistics and data relating to an org
type OrgStat struct {
	NumAps                int `json:"num_aps,omitempty"`
	NumGateways           int `json:"num_gateways,omitempty"`
	NumMxedges            int `json:"num_mxedges,omitempty"`
	NumSwitches           int `json:"num_switches,omitempty"`
	NumUnassignedAps      int `json:"num_unassigned_aps,omitempty"`
	NumUnassignedGateways int `json:"num_unassigned_gateways,omitempty"`
	NumUnassignedSwitches int `json:"num_unassigned_switches,omitempty"`
}

// SiteStat holds operational statistics and data relating to a Site
type SiteStat struct {
	Site

	Lat                 float32 `json:"lat,omitempty"`
	Lng                 float32 `json:"lng,omitempty"`
	MSPID               string  `json:"msp_id,omitempty"`
	NumAP               int     `json:"num_ap,omitempty"`
	NumAPConnected      int     `json:"num_ap_connected,omitempty"`
	NumClients          int     `json:"num_clients,omitempty"`
	NumDevices          int     `json:"num_devices,omitempty"`
	NumDevicesConnected int     `json:"num_devices_connected,omitempty"`
	NumGateway          int     `json:"num_gateway,omitempty"`
	NumGatewayConnected int     `json:"num_gateway_connected,omitempty"`
	NumSwitch           int     `json:"num_switch,omitempty"`
	NumSwitchConnected  int     `json:"num_switch_connected,omitempty"`
	TZOffset            int     `json:"tzoffset,omitempty"`
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

	CPUUtil          int          `json:"cpu_util,omitempty"`
	EnvStat          EnvStat      `json:"env_stat,omitempty"`
	LastSeen         UnixTime     `json:"last_seen,omitzero"`
	NumClients       int          `json:"num_clients,omitempty"`
	Version          string       `json:"version,omitempty"`
	Status           DeviceStatus `json:"status,omitempty"`
	IP               netip.Addr   `json:"ip,omitempty"`
	ExtIP            netip.Addr   `json:"ext_ip,omitempty"`
	NumWLANs         int          `json:"num_wlans,omitempty"`
	Uptime           Seconds      `json:"uptime,omitempty"`
	TxBps            int          `json:"tx_bps,omitempty"`
	RxBps            int          `json:"rx_bps,omitempty"`
	TxBytes          int          `json:"tx_bytes,omitempty"`
	RxBytes          int          `json:"rx_bytes,omitempty"`
	TxPkts           int          `json:"tx_pkts,omitempty"`
	RxPkts           int          `json:"rx_pkts,omitempty"`
	PowerSrc         string       `json:"power_src,omitempty"`
	PowerBudget      int          `json:"power_budget,omitempty"`
	PowerOpMode      string       `json:"power_op_mode,omitempty"`
	PowerConstrained bool         `json:"power_constrained,omitempty"`

	RadioStats map[RadioConfig]RadioStat `json:"radio_stat,omitempty"`
}

// EnvStat holds operational data relating to a device's environment
type EnvStat struct {
	AccelX       float32 `json:"accel_x,omitempty"`
	AccelY       float32 `json:"accel_y,omitempty"`
	AccelZ       float32 `json:"accel_z,omitempty"`
	AmbientTemp  int     `json:"ambient_temp,omitempty"`
	Attitude     int     `json:"attitude,omitempty"`
	CPUTemp      int     `json:"cpu_temp,omitempty"`
	Humidity     int     `json:"humidity,omitempty"`
	MagneX       float32 `json:"magne_x,omitempty"`
	MagneY       float32 `json:"magne_y,omitempty"`
	MagneZ       float32 `json:"magne_z,omitempty"`
	Pressure     float32 `json:"pressure,omitempty"`
	VcoreVoltage int     `json:"vcore_voltage,omitempty"`
}

// RadioStat holds operational statistics and data relating to the radio connectivity of a Device
type RadioStat struct {
	Mac                    string `json:"mac,omitempty"`
	NumClients             int    `json:"num_clients,omitempty"`
	NumWLANs               int    `json:"num_wlans,omitempty"`
	Channel                int    `json:"channel,omitempty"`
	Bandwidth              int    `json:"bandwidth,omitempty"`
	DynamicChainingEnabled bool   `json:"dynamic_chaining_enabled,omitempty"`
	Power                  int    `json:"power,omitempty"`
	NoiseFloor             int    `json:"noise_floor,omitempty"`
	TxBytes                int    `json:"tx_bytes,omitempty"`
	TxPkts                 int    `json:"tx_pkts,omitempty"`
	RxBytes                int    `json:"rx_bytes,omitempty"`
	RxPkts                 int    `json:"rx_pkts,omitempty"`
	Usage                  string `json:"usage,omitempty"`
	UtilAll                int    `json:"util_all,omitempty"`
	UtilTx                 int    `json:"util_tx,omitempty"`
	UtilRxInBSS            int    `json:"util_rx_in_bss,omitempty"`
	UtilRxOtherBSS         int    `json:"util_rx_other_bss,omitempty"`
	UtilUnknownWiFi        int    `json:"util_unknown_wifi,omitempty"`
	UtilNonWiFi            int    `json:"util_non_wifi,omitempty"`
	UtilUndecodableWiFi    int    `json:"util_undecodable_wifi,omitempty"`
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

	Uptime      Seconds    `json:"uptime,omitempty"`
	Idletime    Seconds    `json:"idle_time,omitempty"`
	PowerSaving bool       `json:"power_saving,omitempty"`
	Band        Radio      `json:"band,omitempty"`
	Proto       Dot11Proto `json:"proto,omitempty"`
	KeyMgmt     string     `json:"key_mgmt,omitempty"`
	DualBand    bool       `json:"dual_band,omitempty"`

	Channel        int     `json:"channel,omitempty"`
	VLANID         string  `json:"vlan_id,omitempty"`
	AirspaceIfname string  `json:"airespace_ifname,omitempty"`
	RSSI           int     `json:"rssi,omitempty"`
	SNR            int     `json:"snr,omitempty"`
	TxRate         float32 `json:"tx_rate,omitempty"`
	RxRate         float32 `json:"rx_rate,omitempty"`

	TxBytes   int `json:"tx_bytes,omitempty"`
	TxBps     int `json:"tx_bps,omitempty"`
	TxPackets int `json:"tx_packets,omitempty"`
	TxRetries int `json:"tx_retries,omitempty"`
	RxBytes   int `json:"rx_bytes,omitempty"`
	RxBps     int `json:"rx_bps,omitempty"`
	RxPackets int `json:"rx_packets,omitempty"`
	RxRetries int `json:"rx_retries,omitempty"`

	MapID          string  `json:"map_id,omitempty"`
	X              float32 `json:"x,omitempty"`
	Y              float32 `json:"y,omitempty"`
	Xm             float32 `json:"x_m,omitempty"`
	Ym             float32 `json:"y_m,omitempty"`
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

// StreamedDeviceStat holds information regarding a device returned by the websockets streaming stats API
type StreamedDeviceStat struct {
	Mac        string                            `json:"mac,omitempty"`
	Version    string                            `json:"version,omitempty"`
	IP         string                            `json:"ip,omitempty"`
	ExtIP      string                            `json:"ext_ip,omitempty"`
	PowerSrc   string                            `json:"power_src,omitempty"`
	Uptime     Seconds                           `json:"uptime,omitempty"`
	LastSeen   UnixTime                          `json:"last_seen,omitempty"`
	NumClients int                               `json:"num_clients,omitempty"`
	IPStat     StreamedIPStat                    `json:"ip_stat,omitempty"`
	RadioStats map[RadioConfig]StreamedRadioStat `json:"radio_stat,omitempty"`
	PortStats  map[string]StreamedPortStat       `json:"port_stat,omitempty"`
	LldpStat   StreamedLldpStat                  `json:"lldp_stat,omitempty"`
	RxBytes    int                               `json:"rx_bytes,omitempty"`
	RxPkts     int                               `json:"rx_pkts,omitempty"`
	TxBytes    int                               `json:"tx_bytes,omitempty"`
	TxPkts     int                               `json:"tx_pkts,omitempty"`
	TxBps      int                               `json:"tx_bps,omitempty"`
	RxBps      int                               `json:"rx_bps,omitempty"`
	CpuStat    StreamedCpuStat                   `json:"cpu_stat,omitempty"`
	MemStat    StreamedMemStat                   `json:"memory_stat,omitempty"`
}

// StreamedIPStat holds the IP addressing information of a device returned by the websockets streaming stats API
type StreamedIPStat struct {
	IP       netip.Addr        `json:"ip,omitempty"`
	Netmask  netip.Addr        `json:"netmask,omitempty"`
	Gateway  netip.Addr        `json:"gateway,omitempty"`
	IP6      netip.Addr        `json:"ip6,omitempty"`
	Netmask6 string            `json:"netmask6,omitempty"`
	Gateway6 netip.Addr        `json:"gateway6,omitempty"`
	DNS      []string          `json:"dns,omitempty"`
	IPs      map[string]string `json:"ips,omitempty"`
}

// StreamedRadioStat holds the radio statistics of a device returned by the websockets streaming stats API
type StreamedRadioStat struct {
	Bandwidth  int    `json:"bandwidth,omitempty"`
	Channel    int    `json:"channel,omitempty"`
	Mac        string `json:"mac,omitempty"`
	NumClients int    `json:"num_clients,omitempty"`
	Power      int    `json:"power,omitempty"`
	RxBytes    int    `json:"rx_bytes,omitempty"`
	RxPkts     int    `json:"rx_pkts,omitempty"`
	TxBytes    int    `json:"tx_bytes,omitempty"`
	TxPkts     int    `json:"tx_pkts,omitempty"`
}

// StreamedPortStat holds the wired port statistics of a device returned by the websockets streaming stats API
type StreamedPortStat struct {
	TxPkts     int  `json:"tx_pkts,omitempty"`
	TxBytes    int  `json:"tx_bytes,omitempty"`
	RxPkts     int  `json:"rx_pkts,omitempty"`
	RxBytes    int  `json:"rx_bytes,omitempty"`
	RxPeakBps  int  `json:"rx_peak_bps,omitempty"`
	TxPeakBps  int  `json:"tx_peak_bps,omitempty"`
	FullDuplex bool `json:"full_duplex,omitempty"`
	Speed      int  `json:"speed,omitempty"`
	Up         bool `json:"up"`
	RxErrors   int  `json:"rx_errors,omitempty"`
}

// StreamedLldpStat holds the LLDP information of a device returned by the websockets streaming stats API
type StreamedLldpStat struct {
	SystemName        string `json:"system_name,omitempty"`
	SystemDesc        string `json:"system_desc,omitempty"`
	MgmtAddr          string `json:"mgmt_addr,omitempty"`
	LldpMedSupported  bool   `json:"lldp_med_supported,omitempty"`
	PortDesc          string `json:"port_desc,omitempty"`
	PortID            string `json:"port_id,omitempty"`
	PowerRequestCount int    `json:"power_request_count,omitempty"`
	PowerAllocated    int    `json:"power_allocated,omitempty"`
	PowerDraw         int    `json:"power_draw,omitempty"`
	PowerRequested    int    `json:"power_requested,omitempty"`
}

// StreamedCpuStat holds the CPU statistics of a device returned by the websockets streaming stats API
type StreamedCpuStat struct {
	System    int       `json:"system,omitempty"`
	Idle      int       `json:"idle,omitempty"`
	Interrupt int       `json:"interrupt,omitempty"`
	User      int       `json:"user,omitempty"`
	LoadAvg   []float32 `json:"load_avg,omitempty"`
}

// StreamedMemStat holds the memory statistics of a device returned by the websockets streaming stats API
type StreamedMemStat struct {
	Usage int `json:"usage,omitempty"`
}

// StreamedClientStat holds information about a client returned by the websockets streaming stats API
type StreamedClientStat struct {
	Client
}
