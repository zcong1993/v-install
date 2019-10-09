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
		VmessWsUUID:         GeneratePassword(100),
		VmessWsPort:         20003,
		VmessWsPath:         "/" + GeneratePassword(5),
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
	ip, err := GetPublicIp()
	if err != nil {
		fmt.Printf("warn: get ip error %s\n\n", err.Error())
	}
	for _, inbound := range cfg.Inbounds {
		if inbound.Protocol == "vmess" {
			fmt.Println("-------------------------------------------")
			fmt.Printf("IP: %s\n", ip)
			fmt.Printf("Type: %s\n", inbound.Protocol)
			fmt.Printf("Port: %d\n", inbound.Port)
			if inbound.StreamSettings.Network != "" {
				fmt.Printf("Network: %s\n", inbound.StreamSettings.Network)
				fmt.Printf("WsPath: %s\n", inbound.StreamSettings.WsSettings.Path)
			}
			// clients length > 0
			fmt.Printf("ID: %s\n", inbound.Settings.Clients[0].ID)
			fmt.Printf("AlertID: %d\n", inbound.Settings.Clients[0].AlterID)
		} else if inbound.Protocol == "shadowsocks" {
			fmt.Println("-------------------------------------------")
			fmt.Printf("IP: %s\n", ip)
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
      "port": {{ .VmessWsPort }},
      "listen": "0.0.0.0",
      "tag": "vmess-ws",
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "{{ .VmessWsUUID }}",
            "alterId": 64
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "{{ .VmessWsPath }}"
        }
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
