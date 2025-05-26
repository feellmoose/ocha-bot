package main

import (
	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
	"log"
	"ocha_server_bot/command"
	"ocha_server_bot/command/mine"
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
	repoRank := helper.NewFileRepo(home, "mine_rank")

	langRepo := helper.NewLanguageRepo(repoLanguage)

	rank := helper.NewQueueRank(repoRank, 100, func(a any) float64 {
		if s, ok := a.(mine.TelegramMineGameScore); ok {
			return s.Score
		}
		return 0
	})

	menu := command.NewMenuCommandExec(langRepo)
	mi := command.NewMineCommandExec(repoMine, rank, langRepo, menu)
	help := command.NewHelpCommandExec(langRepo)
	lang := command.NewLanguageCommandExec(langRepo, menu)
	stat := command.NewStatusCommandExec([]helper.Repo{repoLanguage, repoMine}, langRepo)

	bot.Use(middleware.Recover(func(err error, c telebot.Context) {
		log.Printf("Bot error: %v (in context: %v)", err, c.Text())
		if err = c.Send(helper.Messages[langRepo.Context(c)]["error"].Execute(map[string]string{
			"Username": c.Sender().Username,
			"Message":  err.Error(),
		})); err != nil {
			log.Printf("Bot err Sent failed: %v", err)
		}
	}))

	bot.Handle("/mine", mi.Mine)
	bot.Handle("\fmine", mi.Mine)
	bot.Handle("\fflag", mi.Flag)
	bot.Handle("\fback", mi.Rollback)
	bot.Handle("\fquit", mi.Quit)
	bot.Handle("\fclick", mi.Click)
	bot.Handle("\fchange", mi.Change)

	bot.Handle("\fmine_r", mi.MineRank)

	bot.Handle("/help", help.Help)

	bot.Handle("/menu", menu.Menu)
	bot.Handle("\fmenu", menu.Menu)

	bot.Handle("/lang", lang.Language)
	bot.Handle("/lang_chat", lang.LanguageChat)
	bot.Handle("\flang", lang.Language)
	bot.Handle("\flang_chat", lang.LanguageChat)

	bot.Handle("/stat", stat.Status)
	bot.Handle("/stat_mine", stat.StatusMine)

	log.Println("Bot started")
	bot.Start()
}
