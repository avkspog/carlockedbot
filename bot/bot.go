package bot

import (
	"fmt"
	"log"
	"os"
)

const (
	APIURL  = "https://api.telegram.org/bot%s"
	LogFile = "botlog.txt"
)

type Bot struct {
	Token        string
	APIURL       string
	shutdownChan chan<- struct{}
	LogInfo      *log.Logger
	LogError     *log.Logger
}

func NewBot(token string) (*Bot, error) {
	linfo, lerror, err := newLogger()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Token:        token,
		APIURL:       fmt.Sprintf(APIURL, token),
		shutdownChan: make(chan struct{}),
		LogInfo:      linfo,
		LogError:     lerror,
	}

	return bot, nil
}

func (b *Bot) Start() {
	b.LogInfo.Printf("%s started\n", "<botname>") //TODO getMe
}

func (b *Bot) Shutdown() {
	close(b.shutdownChan)
	b.LogInfo.Printf("%s has stopped", "<botname>") //TODO getME
}

func newLogger() (linfo, lerror *log.Logger, err error) {
	file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	linfo = log.New(file, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	lerror = log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	return
}
