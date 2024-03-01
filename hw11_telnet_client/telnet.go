package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	Address string
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	Timeout time.Duration
}

func (cl *Client) Connect() (err error) {
	cl.conn, err = net.Dial("tcp", cl.Address)
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) Close() (err error) {
	return cl.conn.Close()
}

func (cl *Client) Send() error {
	_, err := io.Copy(cl.conn, cl.in)
	return err
}

func (cl *Client) Receive() error {
	_, err := io.Copy(cl.out, cl.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address: address,
		Timeout: timeout,
		in:      in,
		out:     out,
	}
}
