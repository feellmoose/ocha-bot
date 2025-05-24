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
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, context telebot.Context) {
			log.Printf("Bot error(Default OnError): %v (in context: %v)", err, context.Text())
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Panicf("Error create bot: %v", err)
	}

	helper.BotName = bot.Me.Username
	helper.BotID = bot.Me.ID

	repoLanguage := helper.NewFileRepo(home, "language")
	repoMine := helper.NewFileRepo(home, "mine")

	langRepo := helper.NewLanguageRepo(repoLanguage)

	menu := command.NewMenuCommandExec(langRepo)
	mine := command.NewMineCommandExec(repoMine, langRepo, menu)
	help := command.NewHelpCommandExec(langRepo)
	lang := command.NewLanguageCommandExec(langRepo, menu)

	bot.Use(middleware.Recover(func(err error, c telebot.Context) {
		log.Printf("Bot error: %v (in context: %v)", err, c.Text())
		if err = c.Send(helper.Messages[langRepo.Context(c)]["error"].Execute(map[string]string{
			"Username": c.Sender().Username,
			"Message":  err.Error(),
		})); err != nil {
			log.Printf("Bot err Sent failed: %v", err)
		}
	}))

	bot.Handle("/mine", mine.Mine)
	bot.Handle("\fmine", mine.Mine)
	bot.Handle("\fflag", mine.Flag)
	bot.Handle("\fback", mine.Rollback)
	bot.Handle("\fquit", mine.Quit)
	bot.Handle("\fclick", mine.Click)
	bot.Handle("\fchange", mine.Change)

	bot.Handle("/help", help.Help)

	bot.Handle("/menu", menu.Menu)
	bot.Handle("\fmenu", menu.Menu)

	bot.Handle("/lang", lang.Language)
	bot.Handle("/lang_chat", lang.LanguageChat)
	bot.Handle("\flang", lang.Language)
	bot.Handle("\flang_chat", lang.LanguageChat)

	log.Println("Bot started")
	bot.Start()
}
