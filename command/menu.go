package command

import "gopkg.in/telebot.v4"

// MenuCommandFunc support commands:
type MenuCommandFunc interface {
	Menu(c telebot.Context) error
}

/*
/menu menu_id [jump|create] user topic
*/

type MenuCommandExec struct {
}
