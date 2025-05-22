package main

import (
	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
	"log"
	"ocha_server_bot/command"
	"ocha_server_bot/helper"
	"os"
	"time"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	log.Printf("Bot starting with Token(token=%s)", token)
	pref := telebot.Settings{
		Token:   token,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: OnError,
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Panicf("Error create bot: %v", err)
	}

	bot.Use(middleware.Recover())

	mine := command.NewMineCommandExec(helper.NewMemRepo())
	help := command.HelpCommandExec{}
	menu := command.MenuCommandExec{}

	bot.Handle("/mine", mine.Mine)
	bot.Handle(telebot.InlineButton{Unique: "/flag"}, mine.Mine)
	bot.Handle(telebot.InlineButton{Unique: "/back"}, mine.Mine)
	bot.Handle(telebot.InlineButton{Unique: "/quit"}, mine.Mine)
	bot.Handle(telebot.InlineButton{Unique: "/click"}, mine.Mine)
	bot.Handle(telebot.InlineButton{Unique: "/change"}, mine.Mine)

	bot.Handle("/help", help.Help)

	bot.Handle("/menu", menu.Menu)

	log.Println("Bot started")
	bot.Start()
}

func OnError(err error, context telebot.Context) {
	log.Printf("Bot error: %v (in context: %v)", err, context.Text())
}
