package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"vpn_work/config"
	"vpn_work/router"
	"vpn_work/server"
)


var (
	configFile string
	cfg        *config.Config
)


func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	cfg = config.New().FromFile(configFile)
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}

	notify := make(chan struct{}, 1)
	cfg.Wacher(ctx, notify)
	r := router.New(ctx)
	s := server.Server{
		Config: cfg,
		Router: r,
		Notify: notify,
	}
	s.Run(ctx)

	<-sig
}
