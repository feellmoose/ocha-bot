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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	repoLanguage := helper.NewFileRepo(home, "language")
	repoMine := helper.NewFileRepo(home, "mine")
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Panicf("Error create bot: %v", err)
	}

	bot.Use(middleware.Recover())

	langRepo := helper.NewLanguageRepo(repoLanguage)

	menu := command.NewMenuCommandExec(langRepo)
	mine := command.NewMineCommandExec(repoMine, langRepo, menu)
	help := command.NewHelpCommandExec(langRepo)
	lang := command.NewLanguageCommandExec(langRepo)

	bot.Handle("/mine", mine.Mine)
	bot.Handle("\fmine", mine.Mine)
	bot.Handle("\fflag", mine.Flag)
	bot.Handle("\fback", mine.Rollback)
	bot.Handle("\fquit", mine.Quit)
	bot.Handle("\fclick", mine.Click)
	bot.Handle("\fchange", mine.Change)
	bot.Handle(telebot.OnCallback, func(context telebot.Context) error {
		log.Printf("data=%v,args=%v", context.Callback().Data, context.Args())
		return nil
	})

	bot.Handle("/help", help.Help)

	bot.Handle("/menu", menu.Menu)
	bot.Handle("\fmenu", menu.Menu)

	bot.Handle("/language", lang.Language)
	bot.Handle("/language_chat", lang.LanguageChat)

	log.Println("Bot started")
	bot.Start()
}

func OnError(err error, context telebot.Context) {
	log.Printf("Bot error: %v (in context: %v)", err, context.Text())
}
