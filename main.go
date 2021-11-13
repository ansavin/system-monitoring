package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	la, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Println("cant read procfs")
	}
	laOneMin := strings.Split(string(la), " ")[0]
	fmt.Println("la:", laOneMin)
}
