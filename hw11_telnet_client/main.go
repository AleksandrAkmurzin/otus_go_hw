package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func main() {
	if err := run(); err != nil {
		serviceOutput(err.Error())
	}
}

func run() error {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "--timeout=3s")
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)
	if host == "" || port == "" {
		return errors.New("host or port was not set, usage: go-telnet --timeout=5s localhost 4242")
	}

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("error connecting to %s: [%w]", addr, err)
	}
	defer client.Close()
	serviceOutput("Connected to " + addr)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	go func() {
		if err := client.Receive(); err != nil {
			serviceOutput(fmt.Sprintf("error receiving data: [%s]", err))
			return
		}

		serviceOutput("Connection was closed by peer")
		cancel()
	}()

	go func() {
		if err := client.Send(); err != nil {
			serviceOutput(fmt.Sprintf("error sending data: [%s]", err))
			return
		}

		serviceOutput("EOF")
		cancel()
	}()

	<-ctx.Done()

	return nil
}

func serviceOutput(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "...%s\n", msg)
}
