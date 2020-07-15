package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avkspog/carlockedbot/api"
)

type App struct {
	Bot      *api.Bot
	UpdateCh <-chan api.Update
}

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

	updateCh, err := bot.Start(60)
	if err != nil {
		log.Printf("Fatal error: %+v", err)
		os.Exit(1)
	}

	app := &App{
		Bot:      bot,
		UpdateCh: updateCh,
	}

	app.HandleMessages()

	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalCh
		bot.LogInfo.Printf("os signal detected: %+v\n", sig)
		done <- struct{}{}
	}()

	<-done

	bot.Shutdown()
}
