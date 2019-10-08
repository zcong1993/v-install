package core

type Config struct {
	VmessUUID           string
	VmessPort           uint
	ShadowsocksPassword string
	ShadowsocksPort     uint
}

type Inbound struct {
	Protocol string `json:"protocol"`
	Port     uint   `json:"port"`
	Settings struct {
		Method   string `json:"method"`
		Password string `json:"password"`
		Clients  []struct {
			ID string `json:"id"`
		} `json:"clients"`
	} `json:"settings"`
}

type V2rayConfig struct {
	Inbounds []Inbound `json:"inbounds"`
}

type IpRes struct {
	Origin string `json:"origin"`
}
