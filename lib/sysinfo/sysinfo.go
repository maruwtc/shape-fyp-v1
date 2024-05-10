package sysinfo

import (
	"fmt"
	"net"
)

func GetIP() (net.IP, error) {
	var (
		ret    net.IP
		err    error
		ifaces []net.Interface
		addrs  []net.Addr
	)

	if ifaces, err = net.Interfaces(); err == nil {
		for _, i := range ifaces {
			if addrs, err = i.Addrs(); err == nil {
				for _, a := range addrs {
					if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipv4 := ipnet.IP.To4(); ipv4 != nil && ipv4.IsGlobalUnicast() {
							ret = ipv4
							// fmt.Println("IP:", ret.String())
							return ret, nil
						}
					}
				}
			}
		}
	}
	fmt.Println("Error:", err)
	return nil, err
}

// func ListIP() {
// 	ip, err := GetIP()
// 	if err == nil {
// 		fmt.Println("IP:", ip)
// 	} else {
// 		fmt.Println("Error:", err)
// 	}
// }
