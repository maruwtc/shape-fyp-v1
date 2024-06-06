package functions

import (
	"fmt"
	"os/exec"
	"strings"
)

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
