package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
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
		fmt.Printf("%s@%s:%s?> ", user.Name, hostname, working_dir)

		var input string
		fmt.Scanln(&input)

	}
}

func eval(input string) (string, error) {

}

func parse(input string) ([]string, error) {
	var output []string
	for i := 0; i < len(input); {
		// 1. consume all spaces
		for input[i] == ' ' {
			i++
		}
		// 2. handle next token
		if input[i] == '$' {
			if c_bracket := find_matching_c_bracket(input[i+2:]); c_bracket == -1 || input[i+1] != '(' {
				return nil, errors.New("syntax error - incorrect usage of '$()'")
			} else {
				output = append(output, input[i:c_bracket+1])
				i = c_bracket + 1
			}
		} else {
			next := strings.IndexByte(input, ' ')
			output = append(output, input[i:next])
			i = next
		}
	}
	return output, nil
}

func find_matching_c_bracket(str string) int {
	var balance int = 1
	for i := 0; i < len(str); i++ {
		if str[i] == '(' {
			balance++
		} else if str[i] == ')' {
			balance--
		}
		if balance == 0 {
			return i
		}
	}
	return -1
}
