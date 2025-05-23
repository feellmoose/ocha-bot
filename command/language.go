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
		if err := l.repo.SetUserLanguageByContext(c, c.Args()[0]); err != nil {
			return err
		}
	} else {
		if err := l.repo.SetUserLanguageByContext(c, "en"); err != nil {
			return err
		}
	}
	lang := l.repo.Context(c)
	text, err := helper.Messages[lang]["language.note"].Execute(map[string]string{"Username": c.Sender().Username})
	if err != nil {
		return err
	}
	return c.Send(text)
}

func (l LanguageCommandExec) LanguageChat(c telebot.Context) error {
	if len(c.Args()) == 1 {
		if err := l.repo.SetChatLanguageIfAdminByContext(c, c.Args()[0]); err != nil {
			return err
		}
	} else {
		if err := l.repo.SetChatLanguageIfAdminByContext(c, "en"); err != nil {
			return err
		}
	}
	lang := l.repo.Context(c)
	text, err := helper.Messages[lang]["language.note"].Execute(map[string]string{"Username": c.Sender().Username})
	if err != nil {
		return err
	}
	return c.Send(text)
}
