package command

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/command/mine"
	"ocha_server_bot/helper"
	"strconv"
	"time"
)

type StatusCommandFunc interface {
	Status(c telebot.Context) error
	StatusMine(c telebot.Context) error
}

type StatusCommandExec struct {
	repos []helper.Repo
	lang  helper.LanguageRepoFunc
}

func NewStatusCommandExec(repos []helper.Repo, lang helper.LanguageRepoFunc) *StatusCommandExec {
	return &StatusCommandExec{
		repos: repos,
		lang:  lang,
	}
}

func (s *StatusCommandExec) StatusMine(c telebot.Context) error {
	active, running, total := s.analysisMineGame()
	text, err := helper.Messages[s.lang.Context(c)]["stat.game.mine.note"].Execute(map[string]string{
		"Active":  strconv.Itoa(active),
		"Running": strconv.Itoa(running),
		"Total":   strconv.Itoa(total),
	})
	if err != nil {
		return err
	}
	return c.Send(text)
}

// Status Command shows Repo info
func (s *StatusCommandExec) Status(c telebot.Context) error {
	repos := ""
	la := s.lang.Context(c)
	for _, repo := range s.repos {
		data := "NaN"
		if file, ok := repo.(*helper.FileRepo); ok {
			data = strconv.FormatInt(file.DataSize(), 10) + " bytes"
		}
		r, err := helper.Messages[la]["stat.repo.note"].Execute(map[string]string{
			"Name":     repo.Name(),
			"Type":     repo.Type(),
			"DataSize": data,
			"ObjsSize": strconv.Itoa(repo.Size()),
		})
		if err != nil {
			return err
		}
		repos = repos + r
	}
	active, running, total := s.analysisMineGame()
	m, err := helper.Messages[la]["stat.game.mine.note"].Execute(map[string]string{
		"Active":  strconv.Itoa(active),
		"Running": strconv.Itoa(running),
		"Total":   strconv.Itoa(total),
	})
	if err != nil {
		return err
	}
	text, err := helper.Messages[la]["stat.all.note"].Execute(map[string]string{
		"BotID":    strconv.FormatInt(helper.BotID, 10),
		"BotName":  helper.BotName,
		"Version":  helper.Version,
		"Update":   helper.Update,
		"RepoSize": strconv.Itoa(len(s.repos)),
		"Repos":    repos,
		"Mine":     m,
		"Now":      time.Now().String(),
	})
	if err != nil {
		return err
	}
	return c.Send(text, telebot.ModeHTML)
}

func (s *StatusCommandExec) analysisMineGame() (active, running, total int) {
	active = 0
	running = 0
	total = 0
	for _, repo := range s.repos {
		if repo.Name() == "mine" {
			after := time.Now().Add(-2 * time.Minute)
			repo.Range(func(key, value any) bool {
				if game, ok := value.(mine.Serialized); ok {
					total++
					if game.Status == mine.Running {
						running++
						if game.Update.After(after) {
							active++
						}
					}
				}
				return true
			})
		}
	}
	return active, running, total
}
