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

	"golang.org/x/sync/errgroup"
)

const defaultTimeout = 10 * time.Second

var ErrReceiveTerminateSignal = errors.New("received terminate signal, connection closed")

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "connected to %s\n", address)

	sendCh := sendRoutine(client)
	receiveCh := receiveRoutine(client)

	g, ctx := errgroup.WithContext(ctx)

	// sender
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ErrReceiveTerminateSignal
		case err = <-sendCh:
			return err
		}
	})

	// receiver
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ErrReceiveTerminateSignal
		case err = <-receiveCh:
			if errors.Is(err, ErrConnClosed) {
				return nil
			}
			return err
		}
	})

	if err = g.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func sendRoutine(client TelnetClient) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		for {
			err := client.Send()
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return errCh
}

func receiveRoutine(client TelnetClient) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		for {
			err := client.Receive()
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return errCh
}
