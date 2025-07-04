package command

import (
	"errors"
	"ocha_server_bot/helper"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
)

type TaskCommandFunc interface {
	Cron(c telebot.Context) error
	Remove(c telebot.Context) error
	List(c telebot.Context) error
	RecoverAll() error
}

type Task struct {
	ID       string    `json:"id,omitempty"`
	Editer   string    `json:"user,omitempty"`
	Chat     int64     `json:"chat,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Message  string    `json:"message,omitempty"`
	Cron     string    `json:"cron,omitempty"`
	Type     string    `json:"type,omitempty"`
	Language string    `json:"language,omitempty"`
	Create   time.Time `json:"create,omitempty"`
}

type TaskCommandExec struct {
	bot      *telebot.Bot
	repo     helper.Repo[Task]
	id       helper.GenID
	langRepo helper.LanguageRepoFunc
	lock     sync.Mutex
	tasks    map[string]cron.EntryID
	cron     *cron.Cron
}

func NewTaskCommandExec(
	bot *telebot.Bot,
	repo helper.Repo[Task],
	langRepo helper.LanguageRepoFunc,

) *TaskCommandExec {
	return &TaskCommandExec{
		bot:      bot,
		repo:     repo,
		id:       &helper.NanoTimeID{},
		langRepo: langRepo,
		lock:     sync.Mutex{},
		tasks:    make(map[string]cron.EntryID),
		cron:     cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))),
	}
}

func (t *TaskCommandExec) Cron(c telebot.Context) error {
	args := c.Args()
	lang := t.langRepo.Context(c)
	if len(args) > 5 {
		cron := strings.Join(args[0:5], " ")
		msg := strings.Join(args[5:len(args)-1], " ")
		id, err := t.id.NextID()
		if err != nil {
			return err
		}
		task := Task{
			ID:       id,
			Editer:   c.Sender().Username,
			Chat:     c.Chat().ID,
			Thread:   c.Message().ThreadID,
			Message:  msg[1 : len(msg)-1],
			Cron:     cron,
			Type:     "message",
			Language: lang,
			Create:   time.Now(),
		}
		if !t.repo.Put(task.ID, task) {
			return errors.New("repo put new task failed")
		}
		if err := t.Recover(task); err != nil {
			return err
		}
	} else {
		if msg, err := helper.Messages[lang]["cron.help.note"].Execute(map[string]string{"Username": c.Sender().Username}); err != nil {
			return err
		} else {
			c.Send(
				msg,
				telebot.ModeHTML,
			)
		}
	}
	return nil
}

func (t *TaskCommandExec) RecoverAll() error {
	t.repo.Range(func(key string, value Task) bool {
		if _, ok := t.tasks[key]; !ok {
			if err := t.Recover(value); err != nil {
				return false
			}
		}
		return true
	})
	return nil
}

func (t *TaskCommandExec) Recover(task Task) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.tasks[task.ID]; ok {
		return errors.New("task exists")
	}
	switch task.Type {
	case "message":
		id, err := t.cron.AddFunc(task.Cron, func() {
			t.bot.Send(
				&telebot.Chat{ID: task.Chat},
				task.Message,
				&telebot.Topic{ThreadID: task.Thread},
			)
		})
		if err != nil {
			if msg, e := helper.Messages[task.Language]["cron.help.note"].Execute(map[string]string{"Username": task.Editer}); e != nil {
				err = errors.Join(err, e)
			} else {
				t.bot.Send(
					&telebot.Chat{ID: task.Chat},
					msg,
					&telebot.Topic{ThreadID: task.Thread},
					telebot.ModeHTML,
				)
			}
			return err
		}
		t.tasks[task.ID] = id
	}
	return nil
}

func (t *TaskCommandExec) Remove(c telebot.Context) error {
	args := c.Args()
	if len(args) == 2 {
		id := args[1]
		if entry, ok := t.tasks[id]; ok {
			t.repo.Del(id)
			t.cron.Remove(entry)
		}
	}
	return nil
}

func (t *TaskCommandExec) List(c telebot.Context) error {
	lang := t.langRepo.Context(c)
	lines := ""
	t.repo.Range(func(key string, value Task) bool {
		lines = lines + "Task\n\t-ID " + value.ID + "\n\t-Cron " + value.Cron + "\n\t-Type " + value.Type + "\n\t-Editer " + value.Editer + "\n"
		if value.Type == "message" {
			last := 10
			len := len(value.Message)
			if len < 10 {
				last = len
			}
			lines = lines + "\t-Message " + value.Message[1:last] + ".." + strconv.Itoa(len) + "\n"
		}
		return true
	})
	msg, err := helper.Messages[lang]["cron.list.note"].Execute(map[string]string{
		"Username":  c.Sender().Username,
		"TaskLines": lines,
		"Update":    time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return err
	}
	return c.Send(msg, telebot.ModeHTML)
}
