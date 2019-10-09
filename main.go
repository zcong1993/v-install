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
	var (
		install bool
		update  bool
		help    bool
		bbr     bool
	)
	flag.BoolVar(&help, "help", false, "If need show help.")
	flag.BoolVar(&install, "install", false, "If install v2ray.")
	flag.BoolVar(&update, "update", false, "If update config.")
	flag.BoolVar(&bbr, "bbr", false, "If enable bbr.")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if !core.IsLinux() {
		core.End("Only support Linux.")
	}

	if bbr {
		core.EnableBbr()
		os.Exit(0)
	}

	serviceExists := core.CheckServiceExists(ServiceName)
	serviceRunning := core.IsServiceRunning(ServiceName)

	fmt.Printf("service exists: %v\n", serviceExists)
	fmt.Printf("service running: %v\n", serviceRunning)

	if serviceRunning {
		if update {
			core.SetupConfig(ConfigPath)
			fmt.Println("Restarting service...")
			core.RestartService(ServiceName)
			core.PrintByPath(ConfigPath)
			os.Exit(0)
		}
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
		core.SetupConfig(ConfigPath)
		fmt.Println("Starting service...")
		core.StartService(ServiceName)
	}

	core.PrintByPath(ConfigPath)
}
