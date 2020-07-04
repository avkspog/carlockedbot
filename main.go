package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avkspog/carlockedbot/api"
)

func main() {
	token := flag.String("token", "", "-token=<token>")
	flag.Parse()

	signalCh := make(chan os.Signal)
	done := make(chan struct{})

	bot, err := api.NewBot(*token)
	if err != nil {
		log.Printf("Fatal error: %+v\n", err)
		os.Exit(1)
	}

	bot.Start()

	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalCh
		bot.LogInfo.Printf("os signal detected: %+v\n", sig)
		done <- struct{}{}
	}()

	<-done

	bot.Shutdown()
}
