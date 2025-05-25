package command

import "gopkg.in/telebot.v4"

type RollCommandFunc interface {
	Roll(c telebot.Context) error
	Join(c telebot.Context) error
}
