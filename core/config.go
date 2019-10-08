package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
)

var tpl = template.Must(template.New("config").Parse(tplString))

func GenerateDefaultConfig() *Config {
	return &Config{
		VmessUUID:           GeneratePassword(100),
		VmessPort:           20001,
		ShadowsocksPassword: GeneratePassword(8),
		ShadowsocksPort:     20002,
	}
}

func BuildV2rayConfig(config *Config) ([]byte, error) {
	var out bytes.Buffer
	err := tpl.Execute(&out, config)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func ParseConfig(configPath string) *V2rayConfig {
	out, err := ioutil.ReadFile(configPath)
	Failed(err, nil)
	return ParseConfigByte(out)
}

func ParseConfigByte(data []byte) *V2rayConfig {
	var config V2rayConfig
	err := json.Unmarshal(data, &config)
	Failed(err, nil)
	return &config
}

func PrintConfig(cfg *V2rayConfig) {
	for _, inbound := range cfg.Inbounds {
		if inbound.Protocol == "vmess" {
			fmt.Println("-------------------------------------------")
			fmt.Printf("Type: %s\n", inbound.Protocol)
			fmt.Printf("Port: %d\n", inbound.Port)
			// clients length > 0
			fmt.Printf("ID: %s\n", inbound.Settings.Clients[0].ID)
		} else if inbound.Protocol == "shadowsocks" {
			fmt.Println("-------------------------------------------")
			fmt.Printf("Type: %s\n", inbound.Protocol)
			fmt.Printf("Port: %d\n", inbound.Port)
			fmt.Printf("Method: %s\n", inbound.Settings.Method)
			fmt.Printf("Password: %s\n", inbound.Settings.Password)
		}
	}
}

func PrintByPath(configPath string) {
	cfg := ParseConfig(configPath)
	PrintConfig(cfg)
}

func PutConfig(filePath string, data []byte) {
	err := ioutil.WriteFile(filePath, data, 0600)
	Failed(err, data)
}

func SetupConfig(configPath string) {
	config := GenerateDefaultConfig()
	cfg, err := BuildV2rayConfig(config)
	Failed(err, cfg)
	fmt.Println("Writing config...")
	PutConfig(configPath, cfg)
}

const tplString = `
{
  "inbounds": [
    {
      "port": {{ .VmessPort }},
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "{{ .VmessUUID }}",
            "level": 1,
            "alterId": 64
          }
        ]
      }
    },
    {
      "tag": "ss",
      "port": {{ .ShadowsocksPort }},
      "protocol": "shadowsocks",
      "settings": {
        "email": "love@v2ray.com",
        "method": "aes-256-cfb",
        "password": "{{ .ShadowsocksPassword }}",
        "level": 0,
        "ota": false,
        "network": "tcp"
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    },
    {
      "protocol": "blackhole",
      "settings": {},
      "tag": "blocked"
    },
    {
      "protocol": "freedom",
      "settings": {},
      "tag": "out"
    },
    {
      "tag": "tg-out",
      "protocol": "mtproto",
      "settings": {}
    }
  ],
  "routing": {
    "rules": [
      {
        "type": "field",
        "ip": [
          "geoip:private"
        ],
        "outboundTag": "blocked"
      },
      {
        "type": "field",
        "inboundTag": [
          "tg-in"
        ],
        "outboundTag": "tg-out"
      },
      {
        "type": "field",
        "inboundTag": [
          "socks5",
          "ss"
        ],
        "outboundTag": "out"
      }
    ]
  }
}
`
