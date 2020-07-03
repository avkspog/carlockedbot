package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avkspog/carlockedbot/bot"
)

func main() {
	token := flag.String("token", "", "-token=<token>")
	flag.Parse()

	signalCh := make(chan os.Signal)
	done := make(chan struct{})

	bot, err := bot.NewBot(*token)
	if err != nil {
		log.Printf("Fatal error: %+v\n", err)
		os.Exit(1)
	}

	bot.LogInfo.Println(bot.Me)

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
