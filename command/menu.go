package command

import (
	"errors"
	"gopkg.in/telebot.v4"
	"log"
	"ocha_server_bot/helper"
	"strconv"
)

// MenuCommandFunc support commands:
type MenuCommandFunc interface {
	Menu(c telebot.Context) error
}

/*
/menu menu_id [jump|create] user topic
*/

type MenuCommandExec struct {
}

// RedirectTo menu_id [jump|create] [user] [topic]
func RedirectTo(c telebot.Context, args ...string) error {
	var (
		menu, opt, text, lang string
		user                  int64
		topic                 int

		reply *telebot.ReplyMarkup
		err   error
	)

	switch len(args) {
	case 1:
		menu = args[0]
		user = c.Sender().ID
		topic = c.Topic().ThreadID
	case 2:
		menu, opt = args[0], args[1]
		user = c.Sender().ID
		topic = c.Topic().ThreadID
	case 4:
		menu, opt = args[0], args[1]
		user, _ = strconv.ParseInt(args[2], 10, 64)
		topic, _ = strconv.Atoi(args[3])
	default:
		return errors.New("menu args len error")
	}

	lang = c.Sender().LanguageCode

	log.Printf("%v", args)

	switch menu {
	case "mine", "mine_classic":
		text, reply, err = mineClassic(user, topic, lang, c)
		if err != nil {
			return err
		}
	}
	switch opt {
	case "create", "":
		_, err = c.Bot().Send(c.Chat(),
			text,
			telebot.SendOptions{
				ThreadID:    topic,
				ReplyMarkup: reply,
			})
		return err
	case "jump":
		return c.Edit(
			text,
			telebot.SendOptions{
				ThreadID:    topic,
				ReplyMarkup: reply,
			})
	}
	return errors.New("menu opt error (opt=" + opt + ")")
}

func (m MenuCommandExec) Menu(c telebot.Context) error {
	return RedirectTo(c, c.Args()...)
}

func mineClassic(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	reply := &telebot.ReplyMarkup{}
	text, _ := helper.Messages[lang]["mine.game.start.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Width":    "8",
		"Height":   "8",
		"Mines":    "10",
	})
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.start.button"].String(),
			"/create",
			"8", "8", "10", strconv.FormatInt(user, 10), strconv.Itoa(topic),
		),
	),
	)
	return text, reply, nil
}
