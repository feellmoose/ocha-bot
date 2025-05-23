package mine

import (
	"gopkg.in/telebot.v4"
	"log"
	"ocha_server_bot/helper"
	"strconv"
)

type Display interface {
	Display(c telebot.Context) error
}

func (t TelegramMineGame) Display(c telebot.Context) error {

	var (
		boxes   = t.Boxes()
		info    = t.Infos()
		buttons [][]telebot.InlineButton
		text    string
		err     error
	)

	log.Printf("%v", info)

	switch t.Status() {
	case Init, UnInit:
		buttons = t.emptyButton()
		text, err = helper.Messages[info.Locale]["mine.game.Start.note"].Execute(map[string]string{
			"Username": c.Sender().Username,
			"Width":    strconv.Itoa(t.Width()),
			"Height":   strconv.Itoa(t.Height()),
			"Mines":    strconv.Itoa(t.Mines()),
		})
		var change string
		if info.Button == BFlag {
			change = "mine.game.opt.click"
		} else {
			change = "mine.game.opt.flag"
		}
		buttons = append(buttons, []telebot.InlineButton{
			{
				Unique: "change",
				Text:   helper.Messages[info.Locale][change].String(),
				Data:   t.ID(),
			},
			{
				Unique: "quit",
				Text:   helper.Messages[info.Locale]["mine.game.opt.quit"].String(),
				Data:   t.ID(),
			},
		})
		if err != nil {
			return err
		}
		_, err = c.Bot().EditReplyMarkup(telebot.StoredMessage{
			MessageID: strconv.Itoa(info.Message),
			ChatID:    info.Chat,
		}, &telebot.ReplyMarkup{InlineKeyboard: buttons})
	case Running:
		buttons = t.runningButton(boxes)
		var change string
		if info.Button == BFlag {
			change = "mine.game.opt.click"
		} else {
			change = "mine.game.opt.flag"
		}
		buttons = append(buttons, []telebot.InlineButton{
			{
				Unique: "change",
				Text:   helper.Messages[info.Locale][change].String(),
				Data:   t.ID(),
			},
			{
				Unique: "quit",
				Text:   helper.Messages[info.Locale]["mine.game.opt.quit"].String(),
				Data:   t.ID(),
			},
		})
		_, err = c.Bot().EditReplyMarkup(telebot.StoredMessage{
			MessageID: strconv.Itoa(info.Message),
			ChatID:    info.Chat,
		}, &telebot.ReplyMarkup{InlineKeyboard: buttons})
	case End:
		buttons = t.endedButton(boxes, t.Win())
		if t.Win() {
			text, err = helper.Messages[info.Locale]["mine.game.Win.note"].Execute(map[string]string{
				"Username": c.Sender().Username,
				"Width":    strconv.Itoa(t.Width()),
				"Height":   strconv.Itoa(t.Height()),
				"Mines":    strconv.Itoa(t.Mines()),
				"Seconds":  strconv.FormatFloat(t.Duration().Seconds(), 'f', 3, 64),
			})
		} else {
			text, err = helper.Messages[info.Locale]["mine.game.lose.note"].Execute(map[string]string{
				"Username": c.Sender().Username,
				"Width":    strconv.Itoa(t.Width()),
				"Height":   strconv.Itoa(t.Height()),
				"Mines":    strconv.Itoa(t.Mines()),
				"Seconds":  strconv.FormatFloat(t.Duration().Seconds(), 'f', 3, 64),
			})
			buttons = append(buttons, []telebot.InlineButton{
				{
					Unique: "mine",
					Text:   helper.Messages[info.Locale]["mine.game.lose.button"].String(),
					Data:   "8|8|10|" + strconv.FormatInt(t.UserID(), 10) + "|" + strconv.Itoa(info.Topic),
				},
			})
		}
		if err != nil {
			return err
		}

		_, err = c.Bot().Edit(telebot.StoredMessage{
			MessageID: strconv.Itoa(info.Message),
			ChatID:    info.Chat,
		}, text, &telebot.ReplyMarkup{InlineKeyboard: buttons})
	}

	return err
}

func (t TelegramMineGame) endedButton(boxes [][]Box, win bool) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, len(boxes))
	for i, row := range boxes {
		buttons[i] = make([]telebot.InlineButton, len(boxes[i]))
		for j, box := range row {
			if box.IsMine() && box.IsClicked() {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   "💥",
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			} else if box.IsMine() && (box.IsFlagged() || win) {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   "✅",
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			} else if box.IsMine() {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   "💣",
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			} else if box.IsFlagged() {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   "🚩",
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			} else if box.IsClicked() {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   strconv.Itoa(box.Num()),
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			} else {
				buttons[i][j] = telebot.InlineButton{
					Unique: "empty",
					Text:   " ",
					Data:   t.ID() + strconv.Itoa(i) + strconv.Itoa(j),
				}
			}
		}
	}
	return buttons
}

func (t TelegramMineGame) runningButton(boxes [][]Box) [][]telebot.InlineButton {
	var action string
	if t.Infos().Button == BFlag {
		action = "flag"
	} else {
		action = "click"
	}

	buttons := make([][]telebot.InlineButton, len(boxes))
	for i, row := range boxes {
		buttons[i] = make([]telebot.InlineButton, len(boxes[i]))
		for j, box := range row {
			if box.IsFlagged() {
				buttons[i][j] = telebot.InlineButton{
					Unique: action,
					Text:   "🚩",
					Data:   t.ID() + "|" + strconv.Itoa(i) + "|" + strconv.Itoa(j),
				}
			} else if box.IsClicked() {
				buttons[i][j] = telebot.InlineButton{
					Unique: action,
					Text:   strconv.Itoa(box.Num()),
					Data:   t.ID() + "|" + strconv.Itoa(i) + "|" + strconv.Itoa(j),
				}
			} else {
				buttons[i][j] = telebot.InlineButton{
					Unique: action,
					Text:   " ",
					Data:   t.ID() + "|" + strconv.Itoa(i) + "|" + strconv.Itoa(j),
				}
			}
		}
	}
	return buttons
}

func (t TelegramMineGame) emptyButton() [][]telebot.InlineButton {
	var action string
	if t.Infos().Button == BFlag {
		action = "flag"
	} else {
		action = "click"
	}

	buttons := make([][]telebot.InlineButton, t.Width())
	for i := 0; i < t.Width(); i++ {
		buttons[i] = make([]telebot.InlineButton, t.Height())
		for j := 0; j < t.Height(); j++ {
			buttons[i][j] = telebot.InlineButton{
				Unique: action,
				Text:   " ",
				Data:   t.ID() + "|" + strconv.Itoa(i) + "|" + strconv.Itoa(j),
			}
		}
	}
	return buttons
}
