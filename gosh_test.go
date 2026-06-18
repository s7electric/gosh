package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{"ls", []string{"ls"}},
		{"ls -l", []string{"ls", "-l"}},
		{"ls -l | grep foo", []string{"ls", "-l", "|", "grep", "foo"}},
	}
	for _, test := range tests {
		output, err := parse(test.input)
		if err != nil {
			t.Errorf("parse(%s) = %v", test.input, output)
			return
		}
		for i := range output {
			if output[i] != test.output[i] {
				t.Errorf("parse(%s)[%d] = %s; want %s", test.input, i, output[i], test.output[i])
				return
			}
		}
	}
}
