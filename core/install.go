package core

import "fmt"

func InstallV2ray() {
	fmt.Println("Download install script...")
	out, err := ExecCmd("curl", "-sLO", "https://install.direct/go.sh")
	Failed(err, out)
	fmt.Println("Install v2ray...")
	out, err = ExecCmd("bash", "go.sh")
	Failed(err, out)
}

func EnableBbr() {
	fmt.Println("Enable bbr...")
	out, err := ExecCmd("bash", "-c", BbrScript)
	Failed(err, out)
	fmt.Println(string(out))
}
