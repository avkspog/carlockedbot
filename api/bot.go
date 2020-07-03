package bot

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
	APIURL  = "https://api.telegram.org/bot%s/%s"
	LogFile = "botlog.txt"
)

type Bot struct {
	LogInfo      *log.Logger
	LogError     *log.Logger
	token        string
	Me           *User
	httpClient   *http.Client
	shutdownChan chan<- struct{}
}

func NewBot(token string) (*Bot, error) {
	linfo, lerror, err := newLogger()
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		LogInfo:      linfo,
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

func (b *Bot) Start() {
	b.LogInfo.Printf("%s started\n", b.Me.FirstName)
}

func (b *Bot) Shutdown() {
	close(b.shutdownChan)
	b.LogInfo.Printf("%s has stopped", b.Me.FirstName)
}

func (b *Bot) Request(api ApiMethod) (Response, error) {
	endpoint := fmt.Sprintf(APIURL, b.token, api.method())
	paramsReader := strings.NewReader(api.requestParam().Encode())

	req, err := http.NewRequest("POST", endpoint, paramsReader)
	if err != nil {
		return Response{}, err
	}

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

func newLogger() (linfo, lerror *log.Logger, err error) {
	file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	linfo = log.New(file, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	lerror = log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	return
}

func (b *Bot) GetMe() (User, error) {
	response, err := b.Request(GetMe{})
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(response.Result, &user)

	return user, nil
}
