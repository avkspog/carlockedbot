package main

import "github.com/avkspog/carlockedbot/api"

func (app *App) HandleMessages() {
	go func() {
		for {
			msg, ok := <-app.UpdateCh
			if !ok {
				continue
			}

			if msg.Message.IsCommand() {
				app.handleCommands(msg)
			} else if msg.Message.IsCarNumber() {

			}
		}
	}()
}

func (app *App) handleCommands(update api.Update) {
	cmd := update.Message.Text

	app.Bot.LogDebug.Println(cmd)

	if cmd == "/start" {
		welcomeText := "Welcome to CarLockedBot\n ... " //TODO
		//TODO send message with keyboard
		message, err := app.Bot.SendMessage(update.Message.Chat.ID, welcomeText)
		if err != nil {
			app.Bot.LogError.Println(err)
		}
		app.Bot.LogDebug.Println("%+v", message)
	}
}
