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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientData{address: address, timeout: timeout, in: in, out: out}
}

type TelnetClientData struct {
	address    string
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
	timeout    time.Duration
}

func (t *TelnetClientData) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	t.connection = conn
	return err
}

func (t *TelnetClientData) Close() error {
	errInClose := t.in.Close()

	errConnClose := t.connection.Close()
	if errConnClose != nil {
		return errConnClose
	}

	return errInClose
}

func (t *TelnetClientData) Send() error {
	_, err := io.Copy(t.connection, t.in)
	return err
}

func (t *TelnetClientData) Receive() error {
	_, err := io.Copy(t.out, t.connection)
	return err
}
