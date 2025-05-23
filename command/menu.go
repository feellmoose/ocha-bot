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

	log.Printf("menu:args=%v", args)

	switch len(args) {
	case 1:
		menu = args[0]
		user = c.Sender().ID
		topic = c.Message().ThreadID
	case 2:
		menu, opt = args[0], args[1]
		user = c.Sender().ID
		topic = c.Message().ThreadID
	case 4:
		menu, opt = args[0], args[1]
		user, _ = strconv.ParseInt(args[2], 10, 64)
		topic, _ = strconv.Atoi(args[3])
	default:
		return errors.New("menu args len error")
	}

	lang = c.Sender().LanguageCode

	switch menu {
	case "cancel":
		return c.Delete()
	case "mine":
		text, reply, err = mineClassic(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "mine_random":
		text, reply, err = buttonRandom(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "mine_easy":
		text, reply, err = buttonClassic(user, topic, lang, 6, 6, 5, c)
		if err != nil {
			return err
		}
	case "mine_normal":
		text, reply, err = buttonClassic(user, topic, lang, 8, 8, 10, c)
		if err != nil {
			return err
		}
	case "mine_hard":
		text, reply, err = buttonClassic(user, topic, lang, 8, 8, 13, c)
		if err != nil {
			return err
		}
	case "mine_nightmare":
		text, reply, err = buttonClassic(user, topic, lang, 8, 8, 17, c)
		if err != nil {
			return err
		}
	}
	switch opt {
	case "create", "":
		_, err = c.Bot().Send(c.Chat(),
			text,
			&telebot.SendOptions{
				ThreadID:    topic,
				ReplyMarkup: reply,
			})
		return err
	case "jump":
		return c.Edit(
			text,
			&telebot.SendOptions{
				ThreadID:    topic,
				ReplyMarkup: reply,
			})
	case "cancel":
		return c.Delete()
	}
	return errors.New("menu opt error (opt=" + opt + ")")
}

func (m MenuCommandExec) Menu(c telebot.Context) error {
	return RedirectTo(c, c.Args()...)
}

func mineClassic(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["mine.game.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.menu.easy.button"].String(),
			"menu",
			"mine_easy",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
		reply.Data(
			helper.Messages[lang]["mine.game.menu.normal.button"].String(),
			"menu",
			"mine_normal",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
		reply.Data(
			helper.Messages[lang]["mine.game.menu.hard.button"].String(),
			"menu",
			"mine_hard",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
	),
		reply.Row(
			reply.Data(
				helper.Messages[lang]["mine.game.menu.nightmare.button"].String(),
				"menu",
				"mine_nightmare",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["mine.game.menu.random.button"].String(),
				"menu",
				"mine_random",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["menu.cancel.button"].String(),
				"menu",
				"cancel",
				"cancel",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
		),
	)
	return text, reply, nil
}

func buttonRandom(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	width := helper.RandomNum(3, 8)
	height := helper.RandomNum(3, 8)
	density := helper.RandomDensity(0.15, 0.25, func(f float64) float64 {
		if f > 0.7 {
			return 0.3 * f
		} else if f < 0.2 {
			return 2.5 * f
		}
		return f
	})
	mines := int(float64(width*height) * density)
	return buttonClassic(user, topic, lang, width, height, mines, c)
}

func buttonClassic(user int64, topic int, lang string, width, height, mines int, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	reply := &telebot.ReplyMarkup{}
	text, _ := helper.Messages[lang]["mine.game.start.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Width":    strconv.Itoa(width),
		"Height":   strconv.Itoa(height),
		"Mines":    strconv.Itoa(mines),
	})
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.start.button"].String(),
			"mine",
			strconv.Itoa(width),
			strconv.Itoa(height),
			strconv.Itoa(mines),
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
		reply.Data(
			helper.Messages[lang]["menu.back.button"].String(),
			"menu",
			"mine",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
	),
	)
	return text, reply, nil
}

func RedirectToButtonClassic(width, height, mines int, c telebot.Context) error {
	lang := c.Sender().LanguageCode
	reply := &telebot.ReplyMarkup{}
	text, _ := helper.Messages[lang]["mine.game.start.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Width":    strconv.Itoa(width),
		"Height":   strconv.Itoa(height),
		"Mines":    strconv.Itoa(mines),
	})
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.start.button"].String(),
			"mine",
			strconv.Itoa(width),
			strconv.Itoa(height),
			strconv.Itoa(mines),
			strconv.FormatInt(c.Sender().ID, 10),
			strconv.Itoa(c.Message().ThreadID),
		),
	),
	)
	return c.Send(
		text,
		&telebot.SendOptions{
			ThreadID:    c.Message().ThreadID,
			ReplyMarkup: reply,
		})
}
