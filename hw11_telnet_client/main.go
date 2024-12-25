package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

var (
	ErrConnTimeout            = errors.New("connection timed out")
	ErrReceiveTerminateSignal = errors.New("received terminate signal, connection closed")
)

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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
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

	wg := sync.WaitGroup{}
	wg.Add(2)

	// sender
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					fmt.Fprintln(os.Stderr, ErrConnTimeout)
				} else {
					fmt.Fprintln(os.Stderr, ErrReceiveTerminateSignal)
				}

				return
			case err = <-sendCh:
				fmt.Fprintln(os.Stderr, err)
				client.Close()
				return
			}
		}
	}()

	// receiver
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case err = <-receiveCh:
				if !errors.Is(err, ErrConnClosed) {
					fmt.Fprintln(os.Stderr, err)
				}
				return
			}
		}
	}()

	wg.Wait()
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
