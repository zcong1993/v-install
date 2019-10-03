package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func ExecCmd(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	return out, nil
}

func Failed(err error, out []byte) {
	if err != nil {
		fmt.Println(string(out))
		log.Fatal(err)
	}
}

func CheckServiceExists(name string) bool {
	out, _ := ExecCmd("service", name, "status")
	isNotFound := strings.Contains(string(out), "could not be found.")
	if isNotFound {
		return false
	}
	return !strings.Contains(string(out), "unrecognized service")
}

func IsServiceRunning(name string) bool {
	out, _ := ExecCmd("service", name, "status")
	return strings.Contains(string(out), "Active: active (running)")
}

func StartService(name string) {
	_, err := ExecCmd("service", name, "start")
	Failed(err, nil)
}

func End(message string) {
	fmt.Println(message)
	os.Exit(0)
}

func GeneratePassword(l uint) string {
	id := uuid.New().String()
	if l < 36 {
		return id[:l]
	}
	return id
}
