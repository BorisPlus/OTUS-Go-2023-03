package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address:    address,
		Timeout:    timeout,
		In:         in,
		Out:        out,
		connection: nil,
	}
}

type Client struct {
	Address    string
	Timeout    time.Duration
	In         io.ReadCloser
	Out        io.Writer
	connection net.Conn
}

var DEBUG = false

func (c *Client) Connect() error {
	if !DEBUG {
		log.SetOutput(io.Discard)
	}
	var dialer net.Dialer
	ctx := context.Background()
	var cancel context.CancelFunc
	if c.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}
	fmt.Printf("...Try connect to %s\n", c.Address)
	connection, err := dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		log.Printf("Failed to connect: %v\n", err)
		return err
	}
	fmt.Printf("...Connected to %s\n", c.Address)
	c.connection = connection
	return nil
}

func (c *Client) Close() error {
	if !DEBUG {
		log.SetOutput(io.Discard)
	}
	if c.connection == nil {
		return nil
	}
	err := c.connection.Close()
	if err != nil {
		log.Printf("Failed to close: %v\n", err)
		if c.connection != nil {
			c.connection = nil
			fmt.Printf("...Force disconnect from %s\n", c.Address)
		}
		return err
	}
	fmt.Printf("...Disconnected from %s\n", c.Address)
	return nil
}

func (c *Client) Send() error {
	if !DEBUG {
		log.SetOutput(io.Discard)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(c.In); err != nil {
		log.Println(err)
		return err
	}
	if _, err := c.connection.Write(buf.Bytes()); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) Receive() error {
	if !DEBUG {
		log.SetOutput(io.Discard)
	}
	if _, err := io.Copy(c.Out, c.connection); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
