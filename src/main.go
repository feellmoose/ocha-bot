package main

import (
	"gopkg.in/telebot.v4"
	"log"
	"ocha_server_bot/src/command"
	"ocha_server_bot/src/helper"
	"os"
	"time"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	log.Printf("Bot starting with Token(token=%s)", token)
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, context telebot.Context) {
			log.Printf("Bot error: %v (in context: %v)", err, context.Text())
		},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Panicf("Error create bot: %v", err)
	}
	mine := command.NewMineCommandExec(helper.NewMemRepo())
	help := command.HelpCommandExec{}
	bot.Handle("/mine", mine.Mine)
	bot.Handle("/flag", mine.Flag)
	bot.Handle("/back", mine.Rollback)
	bot.Handle("/quit", mine.Quit)
	bot.Handle("/click", mine.Click)
	bot.Handle("/change", mine.Change)

	bot.Handle("/help", help.Help)

	log.Println("Bot started")
	bot.Start()
}
