package command

import (
	"errors"
	"gopkg.in/telebot.v4"
	"ocha_server_bot/src/command/mine"
	"ocha_server_bot/src/helper"
	"strconv"
)

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

func (m *MineCommandExec) MineMessage(args []string, message telebot.Message) error {
	width, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	height, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	mines, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	return m.Mine(
		width, height, mines,
		message.ID,
		message.ThreadID,
		message.Sender.ID,
		message.Chat.ID,
		message.Sender.LanguageCode,
		mine.ClassicBottom)
}

func (m *MineCommandExec) Mine(width, height, mines, message, topic int, user, chat int64, locale string, t mine.GameType) error {
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

func (m *MineCommandExec) Dig(id string, x, y int) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized)

		game = game.
			Deserialize().
			OnClicked(mine.Position{X: x, Y: y}).
			Serialize()

		if !m.repo.Put(id, game) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) Flag(id string, x, y int) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized)

		game = game.
			Deserialize().
			OnFlagged(mine.Position{X: x, Y: y}).
			Serialize()

		if !m.repo.Put(id, game) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}

func (m *MineCommandExec) Change(id string) error {
	if data, ok := m.repo.Get(id); ok {
		serialized := data.(mine.Serialized)

		game := serialized.Deserialize()

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

func (m *MineCommandExec) Rollback(id string) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized)

		game = game.
			Deserialize().
			OnRollback(1).
			Serialize()

		if !m.repo.Put(id, game) {
			return errors.New("put repo failed")
		}
		return nil
	}
	return nil
}
