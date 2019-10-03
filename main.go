package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zcong1993/v-installer/core"
)

const ServiceName = "v2ray"
const ConfigPath = "/etc/v2ray/config.json"

func main() {
	var install bool
	flag.BoolVar(&install, "install", false, "If install v2ray.")
	flag.Parse()

	if !core.IsLinux() {
		core.End("Only support Linux.")
	}

	serviceExists := core.CheckServiceExists(ServiceName)
	serviceRunning := core.IsServiceRunning(ServiceName)

	fmt.Printf("service exists: %v\n", serviceExists)
	fmt.Printf("service running: %v\n", serviceRunning)

	if serviceRunning {
		fmt.Println("Already running.")
		core.PrintByPath(ConfigPath)
		os.Exit(0)
	}

	if !install {
		if !serviceExists {
			core.End("Should run -install first.")
		}
		core.PrintByPath(ConfigPath)
		os.Exit(0)
	}

	if !serviceExists {
		core.InstallV2ray()
		config := core.GenerateDefaultConfig()
		cfg, err := core.BuildV2rayConfig(config)
		core.Failed(err, cfg)
		fmt.Println("Writing config...")
		core.PutConfig(ConfigPath, cfg)
		fmt.Println("Starting service...")
		core.StartService(ServiceName)
	}
	core.PrintByPath(ConfigPath)
}
