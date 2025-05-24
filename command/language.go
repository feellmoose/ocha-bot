package command

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
	"strconv"
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
	menu *MenuCommandExec
}

func NewLanguageCommandExec(repo *helper.LanguageRepo, menu *MenuCommandExec) *LanguageCommandExec {
	return &LanguageCommandExec{
		repo: repo,
		menu: menu,
	}
}

func (l LanguageCommandExec) Language(c telebot.Context) error {
	args := c.Args()
	switch len(args) {
	case 0:
		return l.menu.RedirectTo(c, "lang")
	case 2:
		user, _ := strconv.ParseInt(args[1], 10, 64)
		if c.Sender().ID != user {
			return nil
		}
	}
	if err := l.repo.SetUserLanguageByContext(c, args[0]); err != nil {
		return err
	}
	lang := l.repo.Context(c)
	text, err := helper.Messages[lang]["lang.note"].Execute(map[string]string{"Username": c.Sender().Username})
	if err != nil {
		return err
	}
	return c.Send(text)
}

func (l LanguageCommandExec) LanguageChat(c telebot.Context) error {
	args := c.Args()
	switch len(args) {
	case 0:
		return l.menu.RedirectTo(c, "lang_chat")
	case 2:
		user, _ := strconv.ParseInt(args[1], 10, 64)
		if c.Sender().ID != user {
			return nil
		}
	}
	if err := l.repo.SetChatLanguageIfAdminByContext(c, args[0]); err != nil {
		return err
	}
	lang := l.repo.Context(c)
	text, err := helper.Messages[lang]["lang.chat.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"ChatName": c.Chat().Username,
	})
	if err != nil {
		return err
	}
	return c.Send(text)
}
