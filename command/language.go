package command

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
)

type LanguageCommandFunc interface {
	Language(c telebot.Context) error
	LanguageChat(c telebot.Context) error
}

/*
/language      [zh|en|cxg]
/language_chat [zh|en|cxg]
*/

type LanguageCommandExec struct {
	repo *helper.LanguageRepo
}

func NewLanguageCommandExec(repo *helper.LanguageRepo) *LanguageCommandExec {
	return &LanguageCommandExec{repo: repo}
}

func (l LanguageCommandExec) Language(c telebot.Context) error {
	if len(c.Args()) == 1 {
		return l.repo.SetUserLanguageByContext(c, c.Args()[0])
	} else {
		return l.repo.SetUserLanguageByContext(c, "en")
	}
}

func (l LanguageCommandExec) LanguageChat(c telebot.Context) error {
	if len(c.Args()) == 1 {
		return l.repo.SetChatLanguageIfAdminByContext(c, c.Args()[0])
	} else {
		return l.repo.SetChatLanguageIfAdminByContext(c, "en")
	}
}
