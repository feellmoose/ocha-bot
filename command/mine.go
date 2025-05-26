package command

import (
	"errors"
	"gopkg.in/telebot.v4"
	"ocha_server_bot/command/mine"
	"ocha_server_bot/helper"
	"strconv"
	"time"
)

// MineCommandFunc support commands:
type MineCommandFunc interface {
	Mine(c telebot.Context) error
	MineR(c telebot.Context) error
	MineRank(c telebot.Context) error
	Click(c telebot.Context) error
	Flag(c telebot.Context) error
	Change(c telebot.Context) error
	Rollback(c telebot.Context) error
	Quit(c telebot.Context) error
}

/*
/mine [][][]      user topic {length = 4,6}
/click  game [][]
/flag   game [][]
/back   game
/change game
/quit   game
*/

type MineCommandExec struct {
	repo     helper.Repo
	langRepo helper.LanguageRepoFunc
	id       helper.GenID
	factory  mine.Factory
	menu     MenuCommandFunc
	rank     helper.Ranker
}

func NewMineCommandExec(
	repo helper.Repo,
	rank helper.Ranker,
	langRepo helper.LanguageRepoFunc,
	menu MenuCommandFunc,
) *MineCommandExec {
	return &MineCommandExec{
		repo:     repo,
		langRepo: langRepo,
		factory:  mine.Factory{},
		id:       helper.NewGenRandomRepoShortID(4, 16, 5, repo),
		rank:     rank,
		menu:     menu,
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
	if c.Callback() != nil {

		args := c.Args()
		if len(args) != 5 {
			return errors.New("mine callback args len != 5")
		}

		message = c.Callback().Message.ID
		chat = c.Callback().Message.Chat.ID
		width, _ = strconv.Atoi(args[0])
		height, _ = strconv.Atoi(args[1])
		mines, _ = strconv.Atoi(args[2])
		user, _ = strconv.ParseInt(args[3], 10, 64)
		topic, _ = strconv.Atoi(args[4])

		if user != c.Sender().ID {
			return nil
		}
	} else if c.Message() != nil {
		args := c.Args()
		switch len(args) {
		case 0:
			return m.menu.RedirectTo(c, "mine")
		case 3:
			width, _ = strconv.Atoi(args[0])
			height, _ = strconv.Atoi(args[1])
			mines, _ = strconv.Atoi(args[2])
			return m.menu.RedirectToButtonClassic(width, height, mines, c)
		}
	}
	return m.mine(
		width, height, mines,
		message,
		topic,
		user,
		chat,
		m.langRepo.Context(c),
		mine.Classic,
		c)
}

func (m *MineCommandExec) MineRank(c telebot.Context) error {
	l := m.langRepo.Context(c)
	lines := ""
	for _, rank := range m.rank.Items() {
		score := rank.Item.(mine.TelegramMineGameScore)
		text, err := helper.Messages[l]["mine.game.rank.line.note"].Execute(map[string]string{
			"Index":    strconv.Itoa(rank.Index),
			"Score":    strconv.FormatFloat(rank.Score, 'f', 2, 64),
			"Duration": strconv.FormatInt(score.Duration, 64) + "ms",
			"Username": score.Username,
		})
		if err != nil {
			return err
		}
		lines = lines + text
	}
	text, err := helper.Messages[l]["mine.game.rank.res.note"].Execute(map[string]string{
		"Username":  c.Sender().Username,
		"RankLines": lines,
		"Update":    time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return err
	}
	return c.Send(text, telebot.ModeHTML)
}

func (m *MineCommandExec) MineR(c telebot.Context) error {
	var (
		width   int
		height  int
		mines   int
		message int
		topic   int
		user    int64
		chat    int64
	)
	if c.Callback() != nil {

		args := c.Args()
		if len(args) != 5 {
			return errors.New("mine callback args len != 5")
		}

		message = c.Callback().Message.ID
		chat = c.Callback().Message.Chat.ID
		width, _ = strconv.Atoi(args[0])
		height, _ = strconv.Atoi(args[1])
		mines, _ = strconv.Atoi(args[2])
		user, _ = strconv.ParseInt(args[3], 10, 64)
		topic, _ = strconv.Atoi(args[4])

		if user != c.Sender().ID {
			return nil
		}
	} else if c.Message() != nil {
		return nil
	}
	return m.mine(
		width, height, mines,
		message,
		topic,
		user,
		chat,
		m.langRepo.Context(c),
		mine.Rank,
		c)
}

func (m *MineCommandExec) Click(c telebot.Context) error {
	args := c.Args()
	x, _ := strconv.Atoi(args[1])
	y, _ := strconv.Atoi(args[2])
	return m.click(args[0], c.Sender().ID, x, y, c)
}

func (m *MineCommandExec) Flag(c telebot.Context) error {
	args := c.Args()
	x, _ := strconv.Atoi(args[1])
	y, _ := strconv.Atoi(args[2])
	return m.flag(args[0], c.Sender().ID, x, y, c)
}

func (m *MineCommandExec) Change(c telebot.Context) error {
	return m.change(c.Args()[0], c.Sender().ID, c)
}

func (m *MineCommandExec) Rollback(c telebot.Context) error {
	return m.rollback(c.Args()[0], c.Sender().ID, c)
}

func (m *MineCommandExec) Quit(c telebot.Context) error {
	return m.quit(c.Args()[0], c.Sender().ID, c)
}

func (m *MineCommandExec) mine(width, height, mines, message, topic int, user, chat int64, locale string, t mine.GameType, c telebot.Context) error {
	return m.id.WithID(func(id string) error {
		game, err := m.factory.Empty(id, user, mine.Additional{
			Type:     t,
			Button:   mine.BClick,
			Locale:   locale,
			Topic:    topic,
			Chat:     chat,
			Message:  message,
			Username: c.Sender().Username,
		}, width, height, mines)
		if err != nil {
			return err
		}
		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		switch game.Infos().Type {
		case mine.ClassicBottom, mine.Classic:
			return game.Display(c)
		case mine.Rank:
			return game.RankDisplay(c, m.rank)
		default:
			return game.Display(c)
		}
	})
}

func (m *MineCommandExec) click(id string, user int64, x, y int, c telebot.Context) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		if game.Status() == mine.UnInit {
			var err error
			game, err = m.factory.Init(game.(mine.TelegramMineGame), x, y)
			if err != nil {
				return err
			}
		}
		game = game.OnClicked(mine.Position{X: x, Y: y})

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		switch game.Infos().Type {
		case mine.ClassicBottom, mine.Classic:
			return game.Display(c)
		case mine.Rank:
			return game.RankDisplay(c, m.rank)
		default:
			return game.Display(c)
		}
	}
	return nil
}

func (m *MineCommandExec) flag(id string, user int64, x, y int, c telebot.Context) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() || game.Status() == mine.UnInit {
			return nil
		}
		game = game.OnFlagged(mine.Position{X: x, Y: y})

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return game.Display(c)
	}
	return nil
}

func (m *MineCommandExec) change(id string, user int64, c telebot.Context) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()

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

		game = game.OnInfoChanged(mine.Additional{
			Type:    info.Type,
			Button:  button,
			Locale:  info.Locale,
			Topic:   info.Topic,
			Chat:    info.Chat,
			Message: info.Message,
		})

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return game.Display(c)
	}
	return nil
}

func (m *MineCommandExec) rollback(id string, user int64, c telebot.Context) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		game = game.OnRollback(1)

		if !m.repo.Put(id, game.Serialize()) {
			return errors.New("put repo failed")
		}
		return game.Display(c)
	}
	return nil
}

func (m *MineCommandExec) quit(id string, user int64, c telebot.Context) error {
	if data, ok := m.repo.Get(id); ok {
		game := data.(mine.Serialized).Deserialize()
		if user != game.UserID() {
			return nil
		}
		if !m.repo.Del(id) {
			return errors.New("put repo failed")
		}
		text, err := helper.Messages[m.langRepo.Context(c)]["mine.game.quit.note"].Execute(map[string]string{
			"Username": c.Sender().Username,
		})
		if err != nil {
			return err
		}
		_, err = c.Bot().Edit(telebot.StoredMessage{
			MessageID: strconv.Itoa(game.Infos().Message),
			ChatID:    game.Infos().Chat},
			text,
			&telebot.ReplyMarkup{InlineKeyboard: make([][]telebot.InlineButton, 0)},
		)
		return err
	}
	return nil
}
