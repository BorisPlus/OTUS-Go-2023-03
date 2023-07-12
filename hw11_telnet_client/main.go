package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strconv"
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

func argParse(args []string) (timeout time.Duration, host, port string, err error) {
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
	// HostPortConcat = fmt.Sprintf("tcp://%s:%s", hostArg, portArg)
	HostPortConcat := fmt.Sprintf("%s:%s", hostArg, portArg)
	var AddrPort netip.AddrPort
	AddrPort, err = netip.ParseAddrPort(HostPortConcat)
	if err == nil {
		if strings.HasSuffix(AddrPort.Addr().String(), ".0") {
			err = fmt.Errorf("HOST (ip-address): is gateway")
			return
		}
		host = AddrPort.Addr().String()
		port = fmt.Sprint(AddrPort.Port()) // TODO: parsed PORT is UINT16
		return
	}
	// if HOST is Domain-address
	//
	// REMARK:
	// _ "net"
	// _ "net/url"
	// Functions `url.Parse` and `net.SplitHostPort` allow invalid domain names, bzzzz Ж(.
	//
	// urlStruct, err := url.Parse(HostPortConcat)
	// if err != nil {
	// 	return
	// }
	// //TODO: wtf with invalid `HostPortConcat` `urlStruct.Host`
	// //TODO: parsed PORT is STRING
	// host, port, err = net.SplitHostPort(urlStruct.Host)
	// if err != nil {
	// 	return
	// }
	// if port != portArg {
	// 	return
	// }
	//
	// So I get foreign code base
	// var errD error
	err = checkDomain(hostArg)
	if err != nil {
		// err = errD
		return
	}
	host = hostArg
	var portN uint64
	portN, err = strconv.ParseUint(portArg, 10, 16)
	if err != nil {
		// err = err
		return
	}
	if portN == 0 {
		err = fmt.Errorf("PORT must be not zero")
		return
	}
	port = portArg
	return
}

func main() {
	timeout, host, port, err := argParse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		usage()
		return
	}

	in := &bytes.Buffer{}
	out := os.Stdout

	client := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), timeout, io.NopCloser(in), out)
	defer client.Close()

	err = client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())

	go func(waitGroup *sync.WaitGroup, ctx *context.Context, telnetClient *TelnetClient) {
		errorChannel := make(chan error)
		defer func() {
			waitGroup.Done()
		}()
		for {
			select {
			case <-(*ctx).Done():
				return
			case errorChannel <- (*telnetClient).Receive():
				errReceive := <-errorChannel
				if errReceive != nil {
					fmt.Println(errReceive)
					return
				}
			}
		}
	}(&wg, &ctx, &client)

	go func(waitGroup *sync.WaitGroup, canceler context.CancelFunc, telnetClient *TelnetClient) {
		defer func() {
			waitGroup.Done()
			canceler()
			// channel `close(errReceive)` not help, need below
			client.Close()
		}()
		reader := bufio.NewReader(os.Stdin)
		for {
			response, errSend := reader.ReadString('\n')
			if errSend != nil {
				fmt.Println(errSend)
				return
			}
			_, errSend = in.WriteString(fmt.Sprintf("\r%s", response))
			if errSend != nil {
				fmt.Println(errSend)
				return
			}
			errSend = (*telnetClient).Send()
			if errSend != nil {
				fmt.Println(errSend)
				return
			}
		}
	}(&wg, cancel, &client)

	wg.Wait()
}
