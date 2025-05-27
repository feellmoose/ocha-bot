package command

import (
	"errors"
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
	"strconv"
)

// MenuCommandFunc support commands:
type MenuCommandFunc interface {
	Menu(c telebot.Context) error
	RedirectTo(c telebot.Context, args ...string) error
	RedirectToButtonClassic(width, height, mines int, c telebot.Context) error
}

/*
/menu menu_id [jump|create] user topic
*/

func NewMenuCommandExec(repo helper.LanguageRepoFunc) *MenuCommandExec {
	return &MenuCommandExec{repo: repo}
}

type MenuCommandExec struct {
	repo helper.LanguageRepoFunc
}

// RedirectTo menu_id [jump|create] [user] [topic]
func (m MenuCommandExec) RedirectTo(c telebot.Context, args ...string) error {
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
		topic = c.Message().ThreadID
	case 2:
		menu, opt = args[0], args[1]
		user = c.Sender().ID
		topic = c.Message().ThreadID
	case 4:
		menu, opt = args[0], args[1]
		user, _ = strconv.ParseInt(args[2], 10, 64)
		topic, _ = strconv.Atoi(args[3])
		if user != c.Sender().ID {
			return nil
		}
	default:
		return errors.New("menu args len error")
	}

	lang = m.repo.Context(c)

	switch menu {
	case "cancel":
		return c.Delete()
	case "lang":
		text, reply, err = language(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "lang_chat":
		text, reply, err = languageChat(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "mine_menu":
		text, reply, err = mineMenu(user, topic, lang, c)
		if err != nil {
			return err
		}
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
	case "mine_r":
		text, reply, err = mineRank(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "mine_random_r":
		text, reply, err = buttonRandomRank(user, topic, lang, c)
		if err != nil {
			return err
		}
	case "mine_easy_r":
		text, reply, err = buttonRank(user, topic, lang, 6, 6, 5, c)
		if err != nil {
			return err
		}
	case "mine_normal_r":
		text, reply, err = buttonRank(user, topic, lang, 8, 8, 10, c)
		if err != nil {
			return err
		}
	case "mine_hard_r":
		text, reply, err = buttonRank(user, topic, lang, 8, 8, 13, c)
		if err != nil {
			return err
		}
	case "mine_nightmare_r":
		text, reply, err = buttonRank(user, topic, lang, 8, 8, 17, c)
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
	return m.RedirectTo(c, c.Args()...)
}

func mineMenu(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["mine.game.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.menu.rank.button"].String(),
			"menu",
			"mine",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
		reply.Data(
			helper.Messages[lang]["mine.game.menu.classic.button"].String(),
			"menu",
			"mine_r",
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

func mineClassic(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["mine.game.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(
		reply.Row(
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
				helper.Messages[lang]["menu.back.button"].String(),
				"menu",
				"mine_menu",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
		),
	)
	return text, reply, nil
}

func mineRank(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["mine.game.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(
		reply.Row(
			reply.Data(
				helper.Messages[lang]["mine.game.menu.easy.button"].String(),
				"menu",
				"mine_easy_r",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["mine.game.menu.normal.button"].String(),
				"menu",
				"mine_normal_r",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["mine.game.menu.hard.button"].String(),
				"menu",
				"mine_hard_r",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
		),
		reply.Row(
			reply.Data(
				helper.Messages[lang]["mine.game.menu.nightmare.button"].String(),
				"menu",
				"mine_nightmare_r",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["mine.game.menu.random.button"].String(),
				"menu",
				"mine_random_r",
				"jump",
				strconv.FormatInt(user, 10),
				strconv.Itoa(topic),
			),
			reply.Data(
				helper.Messages[lang]["menu.back.button"].String(),
				"menu",
				"mine_menu",
				"jump",
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

func buttonRandomRank(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
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
	return buttonRank(user, topic, lang, width, height, mines, c)
}

func buttonRank(user int64, topic int, lang string, width, height, mines int, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	reply := &telebot.ReplyMarkup{}
	text, _ := helper.Messages[lang]["mine.game.rank.start.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"Width":    strconv.Itoa(width),
		"Height":   strconv.Itoa(height),
		"Mines":    strconv.Itoa(mines),
	})
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["mine.game.start.button"].String(),
			"mine_r",
			strconv.Itoa(width),
			strconv.Itoa(height),
			strconv.Itoa(mines),
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
		reply.Data(
			helper.Messages[lang]["menu.back.button"].String(),
			"menu",
			"mine_r",
			"jump",
			strconv.FormatInt(user, 10),
			strconv.Itoa(topic),
		),
	),
	)
	return text, reply, nil
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

func language(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["lang.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["lang.en.button"].String(),
			"lang",
			"en",
			strconv.FormatInt(user, 10),
		),
		reply.Data(
			helper.Messages[lang]["lang.zh.button"].String(),
			"lang",
			"zh",
			strconv.FormatInt(user, 10),
		),
		reply.Data(
			helper.Messages[lang]["lang.cxg.button"].String(),
			"lang",
			"cxg",
			strconv.FormatInt(user, 10),
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

func languageChat(user int64, topic int, lang string, c telebot.Context) (string, *telebot.ReplyMarkup, error) {
	text, _ := helper.Messages[lang]["lang.chat.menu.note"].Execute(map[string]string{
		"Username": c.Sender().Username,
		"ChatName": c.Chat().Username,
	})
	reply := &telebot.ReplyMarkup{}
	reply.Inline(reply.Row(
		reply.Data(
			helper.Messages[lang]["lang.en.button"].String(),
			"lang_chat",
			"en",
			strconv.FormatInt(user, 10),
		),
		reply.Data(
			helper.Messages[lang]["lang.zh.button"].String(),
			"lang_chat",
			"zh",
			strconv.FormatInt(user, 10),
		),
		reply.Data(
			helper.Messages[lang]["lang.cxg.button"].String(),
			"lang_chat",
			"cxg",
			strconv.FormatInt(user, 10),
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

func (m MenuCommandExec) RedirectToButtonClassic(width, height, mines int, c telebot.Context) error {
	lang := m.repo.Context(c)
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
