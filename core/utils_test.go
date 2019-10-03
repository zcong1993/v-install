package core

import "testing"

func TestGeneratePassword(t *testing.T) {
	t.Log(GeneratePassword(8))
	t.Log(GeneratePassword(100))
}
