package core

type Config struct {
	VmessUUID           string
	VmessPort           uint
	ShadowsocksPassword string
	ShadowsocksPort     uint
	VmessWsPort         uint
	VmessWsUUID         string
	VmessWsPath         string
}

type Inbound struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Settings struct {
		Method   string `json:"method"`
		Password string `json:"password"`
		Clients  []struct {
			ID      string `json:"id"`
			Level   int    `json:"level"`
			AlterID int    `json:"alterId"`
		} `json:"clients"`
	} `json:"settings"`
	Listen         string `json:"listen,omitempty"`
	Tag            string `json:"tag,omitempty"`
	StreamSettings struct {
		Network    string `json:"network"`
		WsSettings struct {
			Path string `json:"path"`
		} `json:"wsSettings"`
	} `json:"streamSettings,omitempty"`
}

type V2rayConfig struct {
	Inbounds  []Inbound `json:"inbounds"`
	Outbounds []struct {
		Protocol string `json:"protocol"`
		Settings struct {
		} `json:"settings"`
		Tag string `json:"tag,omitempty"`
	} `json:"outbounds"`
	Routing struct {
		Rules []struct {
			Type        string   `json:"type"`
			IP          []string `json:"ip,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			InboundTag  []string `json:"inboundTag,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
}

type IpRes struct {
	Origin string `json:"origin"`
}
