package main

import (
	"bytes"
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

		result, err := eval(input)
		if err != nil {
			fmt.Errorf("%w", err)
		} else {
			fmt.Println(result)
		}
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
			eval_result, err := eval(tokens[i][1:(len(tokens[i]) - 1)])
			if err != nil {
				return "", err
			}
			parsed_result, err := parse(eval_result)
			if err != nil {
				return "", err
			}
			tokens = append(append(tokens[0:i], parsed_result...), tokens[i+1:]...)
			if err != nil {
				return "", err
			}
		}
	}

	// create job queue
	var job_queue []*exec.Cmd
	var argv_buffer []string
	for i := 0; i < len(tokens); i++ {
		// FIX THIS
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
	// assemble pipes
	for i := 1; i < len(job_queue); i++ {
		p0, p1 := io.Pipe()
		job_queue[i-1].Stdout = p1
		job_queue[i].Stdin = p0
	}
	var buf bytes.Buffer
	job_queue[len(job_queue)-1].Stdout = &buf
	// execute job queue
	for i := 0; i < len(job_queue); i++ {
		err := job_queue[i].Start()
		if err != nil {
			return "", err
		}
	}
	for i := 0; i < len(job_queue); i++ {
		job_queue[i].Wait()
	}
	return buf.String(), nil
}

func parse(input string) ([]string, error) {
	var output []string = make([]string, 0, 32)
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
			if next == -1 {
				output = append(output, input[i:])
				break
			} else {
				output = append(output, input[i:next])
				i = next
			}
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
