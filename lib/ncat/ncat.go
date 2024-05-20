package ncat

import (
	"flag"
	"fmt"
	"io"
	"log"
	"myapp/lib/sysinfo"
	"net"
	"os"
)

var ip, err = sysinfo.GetIntIP()
var newip = ip.String()

func StartServer(host string, port int) {
	log.SetFlags(0)
	addr := fmt.Sprintf("%s:%d", newip, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("[+] Listening for connections on %s", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[+] Error accepting connection from client: %s", err)
		} else {
			go processClient(conn)
		}
	}
}

func processClient(conn net.Conn) {
	// _, err := io.Copy(os.Stdout, conn)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// conn.Close()
	defer conn.Close()
	file, err := os.Create("received_data.txt")
	if err != nil {
		log.Printf("[+] Error creating file: %s", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Printf("[+] Error writing data to file: %s", err)
	}
}

var (
	listen = flag.Bool("l", true, "Listen")
	host   = flag.String("h", newip, "Host")
	port   = flag.Int("p", 1304, "Port")
)

func StartNcat() {
	flag.Parse()
	if *listen {
		fmt.Println("[+] Starting ncat server...")
		go StartServer(*host, *port)
		fmt.Println("[+] Successfully started ncat server.")
	}
}
