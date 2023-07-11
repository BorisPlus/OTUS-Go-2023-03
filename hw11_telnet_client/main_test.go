package main

import (
	"fmt"
	"testing"
	"time"
)

func TestArgParsePositive(t *testing.T) {
	type expectedParams struct {
		timeout    time.Duration
		host, port string
	}
	testCases := []struct {
		args     []string
		expected expectedParams
	}{
		{
			args:     []string{"--timeout=10s", "localhost", "4242"},
			expected: expectedParams{timeout: 10000000000, host: "localhost", port: "4242"},
		},
		{
			args:     []string{"--timeout=5s", "localhost.com", "23"},
			expected: expectedParams{timeout: 5000000000, host: "localhost.com", port: "23"},
		},
		{
			args:     []string{"--timeout=11s", "telnet.localhost.com", "23"},
			expected: expectedParams{timeout: 11000000000, host: "telnet.localhost.com", port: "23"},
		},
		{
			args:     []string{"127.0.0.1", "23"},
			expected: expectedParams{timeout: 0, host: "127.0.0.1", port: "23"},
		},
		{
			args:     []string{"--timeout=10s", "1.1.1.1", "65535"},
			expected: expectedParams{timeout: 10000000000, host: "1.1.1.1", port: "65535"},
		},
	}

	for _, testCase := range testCases {
		timeout, host, port, _ := argParse(testCase.args)
		if timeout != testCase.expected.timeout || host != testCase.expected.host || port != testCase.expected.port {
			t.Errorf("%s return unexpected: %s, %s, %s ", testCase.args, timeout, host, port)
		} else {
			fmt.Printf("It's Ok. for args %v\n", testCase.args)
		}
	}
}

func TestArgParseNegative(t *testing.T) {
	testCases := []struct {
		args []string
	}{
		{args: []string{"--timeout=5s", ".40ca1host.com", "23"}},
		{args: []string{"--timeout=11s", "telnet.localhost.1", "23"}},
		{args: []string{"--timeout=11s", "telnet.net."}},
		{args: []string{"127.0.0.1", "--timeout=11s", "23"}},
		{args: []string{"--timeout=10s", "1.1.1.257", "65535"}},
		{args: []string{"--timeout=10s", "1.1.1.1", "65537"}},
		{args: []string{"--timeout=10s", "1.1.1.0", "1"}},
	}
	for _, testCase := range testCases {
		_, _, _, err := argParse(testCase.args)
		if err != nil {
			fmt.Printf("It's Ok. Get expected error %v\n", err)
			fmt.Printf("         for args %v\n", testCase.args)
		} else {
			t.Errorf("Expect error %v\n", testCase.args)
		}
	}
}
