package command

import (
	"gopkg.in/telebot.v4"
	"log"
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
}

var Version = "v0.2.0 (Go Rewrite)"
var Update = time.Now().String()

func (h HelpCommandExec) Help(c telebot.Context) error {
	lang := c.Sender().LanguageCode
	text, err := helper.Messages[lang]["mine.help"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Version":  Version,
		"Update":   Update,
	})
	if err != nil {
		return err
	}
	log.Println(text)
	_, err = c.Bot().Send(c.Message().Chat, text, telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
		ThreadID:  c.Message().ThreadID,
	})
	err = c.Send(text, telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})
	return err
}
