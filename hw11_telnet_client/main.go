package main

import (
	"context"
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
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "--timeout=3s")
	flag.Parse()
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	client := NewTelnetClient(
		addr,
		timeout,
		os.Stdin,
		os.Stdout,
	)

	if err := client.Connect(); err != nil {
		exitWithError(fmt.Errorf("error connecting to %s: [%w]", addr, err))
	}
	serviceOutput("Connected to " + addr)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	go func() {
		err := client.Receive()
		if ctx.Err() != nil {
			os.Exit(0)
		}
		if err != nil {
			exitWithError(fmt.Errorf("error receiving data: [%w]", err))
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			exitWithError(fmt.Errorf("error sending data: [%w]", err))
		}

		serviceOutput("EOF")
		cancel()
	}()

	<-ctx.Done()
}

func exitWithError(err error) {
	serviceOutput(err.Error())
	os.Exit(1)
}

func serviceOutput(message string) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("...%s\n", message))
}
