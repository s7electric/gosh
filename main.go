package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	tokens, err := parse(input)
	if err != nil {
		return "", err
	}
	// pass over for sub-evaluations
	for i := 0; i < len(tokens); i++ {
		if tokens[i][0] == '$' {
			tokens[i], err = eval(tokens[i][1:(len(tokens[i]) - 1)])
			if err != nil {
				return "", err
			}
		}
	}

	// create job queue
	var job_queue []*exec.Cmd
	var argv_buffer []string
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "|" {
			job := exec.Command(argv_buffer[0], argv_buffer...)
			job_queue = append(job_queue, job)
			clear(argv_buffer)
		} else if tokens[i] == "<" {

		} else if tokens[i] == ">" {

		} else {
			argv_buffer = append(argv_buffer, tokens[i])
		}
	}
	// execute job queue
	if len(job_queue) == 1 {
		job_queue[0].Start()
	} else {
		var pipe io.PipeReader
		pipe, err = job_queue[0].StdoutPipe()
		if err != nil {
			return "", err
		}
		for i := 1; i < len(job_queue)-1; i++ {
			job_queue[i].StdinPipe()
		}
	}
}

func parse(input string) ([]string, error) {
	var output []string
	for i := 0; i < len(input); {
		for input[i] == ' ' {
			i++
		}
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
