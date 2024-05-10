package main

import (
	"fmt"
	"myapp/lib/execjndi"
	"myapp/lib/payload"
)

func main() {
	banner()
	startPayload := make(chan bool)
	go func() {
		execjndi.ExecJNDI(startPayload)
	}()
	go func() {
		<-startPayload
		payload.PayloadInput()
	}()
	select {}
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func banner() {
	ban := `
	_______ _______      ___     _______ _______ _   ___ _______ __   __ _______ ___     ___     
	|       |       |    |   |   |       |       | | |   |       |  | |  |       |   |   |   |    
	|    ___|   _   |____|   |   |   _   |    ___| |_|   |  _____|  |_|  |    ___|   |   |   |    
	|   | __|  | |  |____|   |   |  | |  |   | __|       | |_____|       |   |___|   |   |   |    
	|   ||  |  |_|  |    |   |___|  |_|  |   ||  |___    |_____  |       |    ___|   |___|   |___ 
	|   |_| |       |    |       |       |   |_| |   |   |_____| |   _   |   |___|       |       |
	|_______|_______|    |_______|_______|_______|   |___|_______|__| |__|_______|_______|_______|   by Chris Wong
`
	fmt.Println(Cyan + ban + Reset)
	fmt.Println("This is a JNDI exploit tool.")
	fmt.Println("This tool will start a JNDI exploit server and generate a payload.")
	fmt.Println("Prefix with" + Red + " [+] " + Reset + "related to JNDI exploit server.")
	fmt.Println("------------------------------------------------------------------")
}
