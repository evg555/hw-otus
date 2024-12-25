package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrConnClosed = errors.New("connection closed by peer")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}

func (c *telnetClient) Send() error {
	buf := make([]byte, 4096)

	n, err := c.in.Read(buf)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(buf[:n])
	if err != nil {
		return ErrConnClosed
	}

	return nil
}

func (c *telnetClient) Receive() error {
	buf := make([]byte, 4096)

	n, err := c.conn.Read(buf)
	if err != nil {
		return ErrConnClosed
	}

	_, err = c.out.Write(buf[:n])
	if err != nil {
		return err
	}

	return nil
}
