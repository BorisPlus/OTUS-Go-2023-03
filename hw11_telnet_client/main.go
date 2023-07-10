package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func usage() {
	fmt.Println(`Usage: 
  ./program [--timeout=TIMEOUT] HOST PORT
Example: 
  ./program server.tel.net 23
  ./program --timeout=10s 192.168.1.3 4242`)
}

func argParse() (timeout time.Duration, host, port string, err error) {
	args := os.Args[1:]
	if len(args) < 2 || len(args) > 3 {
		err = fmt.Errorf("command run format error")
		return
	}

	deltaArg := 0
	if len(args) == 3 {
		timeoutArg := args[0]
		if !strings.HasPrefix(timeoutArg, "--timeout=") {
			err = fmt.Errorf("TIMEOUT parsing error")
			return
		}
		timeout, err = time.ParseDuration(timeoutArg[10:])
		if err != nil {
			return
		}
		deltaArg = 1
	}
	hostArg := args[0+deltaArg]
	portArg := args[1+deltaArg]
	// if HOST is IP-address
	HostPortConcat := fmt.Sprintf("tcp://%s:%s", hostArg, portArg)
	AddrPort, err := netip.ParseAddrPort(HostPortConcat)
	if err == nil {
		host = AddrPort.Addr().String()
		//TODO: parsed PORT is UINT16
		port = fmt.Sprint(AddrPort.Port())
		return
	}
	// if HOST is Domain-address
	urlStruct, err := url.Parse(HostPortConcat)
	if err != nil {
		return
	}
	//TODO: wtf with invalid `HostPortConcat` `urlStruct.Host`
	//TODO: parsed PORT is STRING
	host, port, err = net.SplitHostPort(urlStruct.Host)
	if err != nil {
		return
	}
	if port != portArg {
		return
	}
	return
}

func main() {
	timeout, host, port, err := argParse()
	fmt.Println(timeout, host, port)
	if err != nil {
		fmt.Println(err)
		usage()
		return
	}

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	client := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), timeout, io.NopCloser(in), out)
	defer client.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			err = client.Receive()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = out.WriteString(response)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = client.Send()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	err = client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	wg.Wait()
}
