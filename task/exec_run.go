package task

import (
	"os/exec"
	"strings"
)

func Exec(cmd string) (output []string, err error) {
	_output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	//time.Sleep(1 * time.Second)
	//fmt.Println("this :", _output)
	output = strings.Split(string(_output), "\n")
	return
}

func Run(cmd string) (err error) {
	err = exec.Command("/bin/bash", "-c", cmd).Run()
	return
}
