package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	APIURL          = "https://api.telegram.org/bot%s/%s"
	LogFile         = "botlog.txt"
	CarNumberPrefix = "ZC"
)

type Bot struct {
	LogInfo      *log.Logger
	LogDebug     *log.Logger
	LogError     *log.Logger
	token        string
	Me           *User
	httpClient   *http.Client
	shutdownChan chan struct{}
}

func NewBot(token string) (*Bot, error) {
	linfo, ldebug, lerror, err := newLoggers()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		LogInfo:      linfo,
		LogDebug:     ldebug,
		LogError:     lerror,
		token:        token,
		shutdownChan: make(chan struct{}),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	me, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Me = &me

	return bot, nil
}

func (b *Bot) Start(timeout int) (<-chan Update, error) {
	b.LogInfo.Printf("%s started\n", b.Me.FirstName)

	method := GetUpdates{
		Limit:   100,
		Timeout: timeout,
	}

	ch := make(chan Update, 100)

	go func() {
		for {
			select {
			case <-b.shutdownChan:
				close(ch)
				return
			default:
			}

			updates, err := b.GetUpdates(method)
			if err != nil {
				b.LogError.Println(err)
				continue
			}

			for _, update := range updates {
				if update.UpdateID >= method.Offset {
					method.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch, nil
}

func (b *Bot) Shutdown() {
	b.shutdownChan <- struct{}{}
	b.LogInfo.Printf("%s has stopped\n", b.Me.FirstName)
}

func (b *Bot) ApiRequest(api ApiMethod) (Response, error) {
	endpoint := fmt.Sprintf(APIURL, b.token, api.method())
	paramsReader := strings.NewReader(api.requestParam().Encode())

	req, err := http.NewRequest("POST", endpoint, paramsReader)
	if err != nil {
		return Response{}, err
	}

	b.LogDebug.Printf("REQUEST: %+v\n", api)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	json.Unmarshal(data, &response)

	return response, nil
}

func newLoggers() (linfo, ldebug, lerror *log.Logger, err error) {
	file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	linfo = log.New(file, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	ldebug = log.New(file, "[DEBUG]: ", log.Ldate|log.Ltime|log.Lshortfile)
	lerror = log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	return
}

func (b *Bot) GetMe() (User, error) {
	resp, err := b.ApiRequest(GetMe{})
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (b *Bot) GetUpdates(method GetUpdates) ([]Update, error) {
	resp, err := b.ApiRequest(method)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)
	if err != nil {
		return []Update{}, err
	}

	return updates, nil
}

func (b *Bot) SendMessage(chatID int64, text string) (Message, error) {
	method := SendMessage{
		ChatID: chatID,
		Text:   text,
	}

	return b.sendMessage(method)
}

func (b *Bot) sendMessage(method SendMessage) (Message, error) {
	resp, err := b.ApiRequest(method)
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}
