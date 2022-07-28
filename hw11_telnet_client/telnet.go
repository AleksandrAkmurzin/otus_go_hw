package main

import (
	"io"
	"net"
	"time"
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return TelnetClient{address: address, timeout: timeout, in: in, out: out}
}

type TelnetClient struct {
	address    string
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
	timeout    time.Duration
}

func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	t.connection = conn
	return err
}

func (t *TelnetClient) Close() error {
	errInClose := t.in.Close()

	errConnClose := t.connection.Close()
	if errConnClose != nil {
		return errConnClose
	}

	return errInClose
}

func (t *TelnetClient) Send() error {
	_, err := io.Copy(t.connection, t.in)
	return err
}

func (t *TelnetClient) Receive() error {
	_, err := io.Copy(t.out, t.connection)
	return err
}
