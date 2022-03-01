package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	c.connection = conn
	return nil
}

func (c *Client) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

func (c *Client) Send() error {
	_, err := io.Copy(c.connection, c.in)
	fmt.Fprintln(os.Stderr, "...EOF")
	return err
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return err
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
