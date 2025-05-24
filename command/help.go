package command

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
)

/*

/menu menu_id [jump|create] user topic

/mine [][][]      user topic {length = 4,6}
/mine level    [] user topic {length = 3,5} level
/mine random      user topic {length = 2,4} random
/mine classic     user topic {length = 2,4} classic
/click  game [][]
/flag   game [][]
/back   game
/change game
/quit   game

/language [zh|en|cxg]
/language_chat [zh|en|cxg]

/help

*/

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

func (h HelpCommandExec) Help(c telebot.Context) error {
	lang := h.repo.Context(c)
	text, err := helper.Messages[lang]["help.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Version":  helper.Version,
		"Update":   helper.Update,
		"BotName":  helper.BotName,
	})
	if err != nil {
		return err
	}
	err = c.Send(text, telebot.ModeHTML)
	return err
}
