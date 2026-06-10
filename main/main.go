package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func main() {
	var stop bool
	for !stop {
		working_dir, err := os.Getwd()
		if err != nil {
			fmt.Errorf("getcwd: %w", err)
		}
		user, err := user.LookupId(strconv.Itoa(os.Geteuid()))
		if err != nil {
			fmt.Errorf("userid: %w", err)
		}
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Errorf("hostname: %w", err)
		}
		fmt.Printf("%s@%s:%s?>", user.Name, hostname, working_dir)

		var input [512]byte
		fmt.Scanln(&input)
	}
}
