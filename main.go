package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/avvvet/comicloud/internal/app"
)

type Config struct {
	HttpPort uint
}

func main() {
	config := Config{}
	ctx, cancel := context.WithCancel(context.Background())
	flag.UintVar(&config.HttpPort, "httpPort", 0, "http port for comicloud")
	flag.Parse()

	app.Art()

	/*
	  http server
	*/
	http := app.NewApp(ctx, config.HttpPort)
	go http.Run()
	run(http, cancel)
}

func run(h *app.HttpServer, cancel func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Printf("\r👋️ stopped...\n")

	cancel()

	// if err := h.(); err != nil {
	// 	panic(err)
	// }
	os.Exit(0)
}
