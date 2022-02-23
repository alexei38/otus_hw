package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	ErrHostNotFound = errors.New("host not found")
	ErrPortNotFound = errors.New("port not found")
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func receive(client TelnetClient, cancel context.CancelFunc) {
	defer cancel()
	if err := client.Receive(); err != nil {
		// Выполняем cancel перед exit 1
		cancel()
		log.Fatalln(err) // nolint: gocritic
	}
}

func send(client TelnetClient, cancel context.CancelFunc) {
	defer cancel()
	if err := client.Send(); err != nil {
		// Выполняем cancel перед exit 1
		cancel()
		log.Fatalln(err) // nolint: gocritic
	}
}

func main() {
	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		log.Fatalln(ErrHostNotFound)
	}
	port := flag.Arg(1)
	if port == "" {
		log.Fatalln(ErrPortNotFound)
	}
	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	ctx, cancel := context.WithCancel(context.Background())

	go receive(client, cancel)
	go send(client, cancel)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
