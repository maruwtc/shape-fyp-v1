package sysinfo

import (
	"fmt"
	"net"

	externalip "github.com/glendc/go-external-ip"
)

func GetIntIP() (net.IP, error) {
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

func GetExtIP() (string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func ListInfo() {
	ip, err := GetIntIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Internal IP:", ip)
	ipString, err := GetExtIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ip = net.ParseIP(ipString)
	fmt.Println("External IP:", ip)
}

func TargetInfo() (string, string) {
	targetip := "192.168.0.132"
	targetport := "8080"
	return targetip, targetport
}
