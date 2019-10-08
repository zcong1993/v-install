package core

import "testing"

func TestPrintConfig(t *testing.T) {
	cfgData, _ := BuildV2rayConfig(&Config{
		VmessPort:           20001,
		VmessUUID:           "xsxscd",
		ShadowsocksPassword: "test-pass",
		ShadowsocksPort:     20002,
	})
	cfg := ParseConfigByte(cfgData)
	PrintConfig(cfg)
}

func TestPutConfig(t *testing.T) {
	PutConfig("./test.json", []byte("[]"))
}
