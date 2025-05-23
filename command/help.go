package command

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
	"time"
)

// HelpCommandFunc support commands:
type HelpCommandFunc interface {
	Help(c telebot.Context) error
}

/*
/help
*/

type HelpCommandExec struct {
	repo *helper.LanguageRepo
}

func NewHelpCommandExec(repo *helper.LanguageRepo) *HelpCommandExec {
	return &HelpCommandExec{repo: repo}
}

var Version = "v0.2.0 (Go Rewrite)"
var Update = time.Now().String()

func (h HelpCommandExec) Help(c telebot.Context) error {
	lang := h.repo.Context(c)
	text, err := helper.Messages[lang]["mine.help"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Version":  Version,
		"Update":   Update,
	})
	if err != nil {
		return err
	}
	err = c.Send(text, telebot.ModeHTML)
	return err
}
