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

var ErrorArgumentsMissing = errors.New("some arguments missing")

func parseArgs() (address string, timeout time.Duration, err error) {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		return "", timeout, ErrorArgumentsMissing
	}
	return net.JoinHostPort(args[0], args[1]), timeout, nil
}

func main() {
	address, timeout, err := parseArgs()
	if err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())

	if err = tc.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer func(tc TelnetClient) {
		if err := tc.Close(); err != nil {
			log.Fatalln(err)
		}
	}(tc)

	go func() {
		defer cancel()
		if err := tc.Receive(); err != nil {
			log.Println(err)
			return
		}
	}()

	go func() {
		defer cancel()
		if err := tc.Send(); err != nil {
			log.Println(err)
			return
		}
	}()

	signal.Notify(ch, os.Interrupt)
	select {
	case <-ch:
		cancel()
	case <-ctx.Done():
		log.Println("Connection was closed by peer")
		close(ch)
	}
}
