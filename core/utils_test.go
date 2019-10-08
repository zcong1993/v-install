package core

import "testing"

func TestGeneratePassword(t *testing.T) {
	t.Log(GeneratePassword(8))
	t.Log(GeneratePassword(100))
}

func TestGetPublicIp(t *testing.T) {
	ip, err := GetPublicIp()
	if err != nil {
		t.Log(err)
	}
	t.Log(ip)
}
