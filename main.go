package main

import (
	"flag"
	"fmt"
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

	updateCh, err := bot.Start(60)
	if err != nil {
		log.Printf("Fatal error: %+v", err)
		os.Exit(1)
	}

	HandleMessage(updateCh)

	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalCh
		bot.LogInfo.Printf("os signal detected: %+v\n", sig)
		done <- struct{}{}
	}()

	<-done

	bot.Shutdown()
}

func HandleMessage(message <-chan api.Update) {
	go func() {
		for {
			msg, ok := <-message
			if !ok {
				continue
			}
			fmt.Println(msg.Message.Text)
		}
	}()
}
