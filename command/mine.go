package command

import (
	"errors"
	"gopkg.in/telebot.v4"
	"log"
	"ocha_server_bot/command/mine"
	"ocha_server_bot/helper"
	"strconv"
)

// MineCommandFunc support commands:
type MineCommandFunc interface {
	Mine(c telebot.Context) error
	Click(c telebot.Context) error
	Flag(c telebot.Context) error
	Change(c telebot.Context) error
	Rollback(c telebot.Context) error
	Quit(c telebot.Context) error
}

/*
/mine [][][]      user topic {length = 4,6}
/mine level    [] user topic {length = 3,5} level
/mine random      user topic {length = 2,4} random
/mine classic     user topic {length = 2,4} classic
/click  game [][]
/flag   game [][]
/back   game
/change game
/quit   game
*/

type MineCommandExec struct {
	repo    *helper.MemRepo
	id      *helper.GenRandomRepoShortID
	factory mine.Factory
}

func NewMineCommandExec(repo *helper.MemRepo) *MineCommandExec {
	return &MineCommandExec{
		repo:    repo,
		factory: mine.Factory{},
		id:      helper.NewGenRandomRepoShortID(4, 16, 5, repo),
	}
}

func (m *MineCommandExec) Mine(c telebot.Context) error {
	var (
		width   int
		height  int
		mines   int
		message int
		topic   int
		user    int64
		chat    int64
	)
	log.Printf("%v", c.Args())
	if c.Message() != nil {
		width, height, mines, message, topic, user, chat = m.handleMineMessage(c)
	} else if c.Callback() != nil {
		width, height, mines, message, topic, user, chat = m.handleMineCallback(c)
	}
	return m.mine(
		width, height, mines,
		message,
		topic,
		user,
		chat,
		c.Sender().LanguageCode,
		mine.ClassicBottom)
}

func (m *MineCommandExec) handleMineMessage(c telebot.Context) (
	width,
	height,
	mines,
	message,
	topic int,
	user,
	chat int64) {

	args := c.Args()
	switch len(args) {
	case 1:
		action := args[0]
		switch action {
		case "classic":
			width = 8
			height = 8
			mines = 10
		case "random":
			width = helper.Num(3, 8)
			height = helper.Num(3, 8)
			density := helper.RandomDensity(0.15, 0.25, func(f float64) float64 {
				if f > 0.7 {
					return 0.3 * f
				} else if f < 0.2 {
					return 2.5 * f
				}
				return f
			})
			mines = int(float64(width*height) * density)
		}
	case 2:
		action := args[0]
		switch action {
		case "level":
			level := args[1]
			switch level {
			case "easy":
				width = 6
				height = 6
				mines = 5
			case "normal":
				width = 8
				height = 8
				mines = 10
			case "hard":
				width = 8
				height = 8
				mines = 13
			}
		}
	case 4:
		width, _ = strconv.Atoi(args[0])
		height, _ = strconv.Atoi(args[1])
		mines, _ = strconv.Atoi(args[2])
	}
	message = c.Message().ID
	topic = c.Message().ThreadID
	user = c.Message().Sender.ID
	chat = c.Message().Chat.ID
	return
}

func (m *MineCommandExec) handleMineCallback(c telebot.Context) (
	width,
	height,
	mines,
	message,
	topic int,
	user,
	chat int64) {

	args := c.Args()
	message = c.Callback().Message.ID
	chat = c.Callback().Message.Chat.ID

	switch len(args) {
	case 3:
		action := args[0]
		switch action {
		case "classic":
			width = 8
			height = 8
			mines = 10
		case "random":
			width = helper.Num(3, 8)
			height = helper.Num(3, 8)
			density := helper.RandomDensity(0.15, 0.25, func(f float64) float64 {
				if f > 0.7 {
					return 0.3 * f
				} else if f < 0.2 {
					return 2.5 * f
				}
				return f
			})
			mines = int(float64(width*height) * density)
		}
		user, _ = strconv.ParseInt(args[1], 10, 64)
		topic, _ = strconv.Atoi(args[2])
	case 4:
		action := args[0]
		switch action {
		case "level":
			level := args[1]
			switch level {
			case "easy":
				width = 6
				height = 6
				mines = 5
			case "normal":
				width = 8
				height = 8
				mines = 10
			case "hard":
				width = 8
				height = 8
				mines = 13
			}
		}
		user, _ = strconv.ParseInt(args[2], 10, 64)
		topic, _ = strconv.Atoi(args[3])
	case 5:
		width, _ = strconv.Atoi(args[0])
		height, _ = strconv.Atoi(args[1])
		mines, _ = strconv.Atoi(args[2])
		user, _ = strconv.ParseInt(args[3], 10, 64)
		topic, _ = strconv.Atoi(args[4])
	}
	return
}

func (m *MineCommandExec) Click(c telebot.Context) error {
	var (
		args = c.Args()
		user int64
	)
	if c.Message() == nil {
		user = c.Message().Sender.ID
	} else {
		user, _ = strconv.ParseInt(args[3], 10, 64)
	}
	id := args[0]
	x, _ := strconv.Atoi(args[1])
	y, _ := strconv.Atoi(args[2])
	return m.click(id, user, x, y)
}

func (m *MineCommandExec) Flag(c telebot.Context) error {
	var (
		args = c.Args()
		user int64
	)
	if c.Message() == nil {
		user = c.Message().Sender.ID
	} else {
		user, _ = strconv.ParseInt(args[3], 10, 64)
	}
	id := args[0]
	x, _ := strconv.Atoi(args[1])
	y, _ := strconv.Atoi(args[2])
	return m.flag(id, user, x, y)
}

func (m *MineCommandExec) Change(c telebot.Context) error {
	var (
		args = c.Args()
		user int64
	)
	if c.Message() == nil {
		user = c.Message().Sender.ID
	} else {
		user, _ = strconv.ParseInt(args[1], 10, 64)
	}
	id := args[0]
	return m.change(id, user)
}

func (m *MineCommandExec) Rollback(c telebot.Context) error {
	var (
		args = c.Args()
		user int64
	)
	if c.Message() == nil {
		user = c.Message().Sender.ID
	} else {
		user, _ = strconv.ParseInt(args[1], 10, 64)
	}
	id := args[0]
	return m.rollback(id, user)
}

func (m *MineCommandExec) Quit(c telebot.Context) error {
	var (
		args = c.Args()
		user int64
	)
	if c.Message() == nil {
		user = c.Message().Sender.ID
	} else {
		user, _ = strconv.ParseInt(args[1], 10, 64)
	}
	id := args[0]
	return m.quit(id, user)
}

func (m *MineCommandExec) mine(width, height, mines, message, topic int, user, chat int64, locale string, t mine.GameType) error {
	err := m.id.WithID(func(id string) error {
		game, err := m.factory.Empty(id, user, mine.Additional{
			Type:    t,
			Button:  mine.BClick,
			Locale:  locale,
			Topic:   topic,
			Chat:    chat,
			Message: message,
		}, width, height, mines)
		if err != nil {
			return err
		}
		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MineCommandExec) click(id string, user int64, x, y int) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		game = game.OnClicked(mine.Position{X: x, Y: y})

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) flag(id string, user int64, x, y int) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		game = game.OnFlagged(mine.Position{X: x, Y: y})

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) change(id string, user int64) error {
	if data, ok := m.repo.Get(id); ok {
		serialized := data.(mine.Serialized)

		game := serialized.Deserialize()

		if user != game.UserID() {
			return nil
		}

		info := game.Infos()
		var button mine.Button
		if info.Button == mine.BClick {
			button = mine.BFlag
		} else {
			button = mine.BClick
		}

		serialized = game.OnInfoChanged(mine.Additional{
			Type:    info.Type,
			Button:  button,
			Locale:  info.Locale,
			Topic:   info.Topic,
			Chat:    info.Chat,
			Message: info.Message,
		}).Serialize()

		if !m.repo.Put(id, serialized) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) rollback(id string, user int64) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		game = game.OnRollback(1)

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) quit(id string, user int64) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		if !m.repo.Del(id) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}
