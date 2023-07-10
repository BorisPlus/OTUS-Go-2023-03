package main

import (
	"bytes"
	"context"
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

func (c *Client) Connect() error {
	var dialer net.Dialer
	ctx := context.Background()
	// var cancel context.CancelFunc
	// if c.Timeout > 0 {
	// 	ctx, cancel = context.WithTimeout(ctx, c.Timeout)
	// 	defer cancel()
	// }
	connection, err := dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		log.Printf("Failed to connect: %v\n", err)
		return err
	}
	c.connection = connection
	return nil
}

func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}
	return c.connection.Close()
}

func (c *Client) Send() error {
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
	if _, err := io.Copy(c.Out, c.connection); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
