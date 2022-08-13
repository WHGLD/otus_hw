package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

const defaultTimeout = 10

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("address should be like host port")
	}

	host := args[0]
	port := args[1]

	cli := NewTelnetClient(host+":"+port, timeout, os.Stdin, os.Stdout)
	if err := cli.Connect(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go sendMsg(cli, cancel)
	go receiveMsg(cli, cancel)

	<-ctx.Done()
}

func sendMsg(cli TelnetClient, cancel context.CancelFunc) {
	if err := cli.Send(); err != nil {
		log.Printf("send error: %v", err)
	}
	cancel()
}

func receiveMsg(cli TelnetClient, cancel context.CancelFunc) {
	if err := cli.Receive(); err != nil {
		log.Printf("receive error: %v", err)
	} else {
		log.Printf("connection was closed by peer")
	}
	cancel()
}
