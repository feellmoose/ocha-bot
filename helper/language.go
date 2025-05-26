package helper

import (
	"errors"
	"gopkg.in/telebot.v4"
	"strconv"
)

type LanguageRepoFunc interface {
	Lang(lang string) string
	Context(c telebot.Context) string
	SetChatLanguageIfAdminByContext(c telebot.Context, lang string) error
	SetUserLanguageByContext(c telebot.Context, lang string) error
}

type LanguageRepo struct {
	repo Repo[string]
}

func NewLanguageRepo(repo Repo[string]) *LanguageRepo {
	return &LanguageRepo{
		repo: repo,
	}
}

func (r LanguageRepo) Lang(lang string) string {
	switch lang {
	case "en":
		return "en"
	case "cxg":
		return "cxg"
	case "zh":
		return "zh"
	default:
		return "en"
	}
}

func (r LanguageRepo) Context(c telebot.Context) string {
	user := strconv.FormatInt(c.Sender().ID, 10)
	if lang, ok := r.repo.Get("us_" + user); ok {
		return lang
	}
	chat := strconv.FormatInt(c.Chat().ID, 10)
	if lang, ok := r.repo.Get("ch_" + chat); ok {
		return lang
	}
	return r.Lang(c.Sender().LanguageCode)
}

func (r LanguageRepo) SetChatLanguageIfAdminByContext(c telebot.Context, lang string) error {
	switch c.Chat().Type {
	case telebot.ChatChannel, telebot.ChatSuperGroup, telebot.ChatGroup:
		members, err := c.Bot().AdminsOf(c.Chat())
		if err != nil {
			return err
		}
		admin := false
		for _, member := range members {
			if member.User.ID == c.Sender().ID {
				admin = true
				break
			}
		}
		if !admin {
			return errors.New("only admin can set chat language")
		}
	}
	if !r.setChatLanguage(strconv.FormatInt(c.Chat().ID, 10), r.Lang(lang)) {
		return errors.New("admin chat language set failed")
	}
	return nil
}

func (r LanguageRepo) SetUserLanguageByContext(c telebot.Context, lang string) error {
	if !r.setUserLanguage(strconv.FormatInt(c.Sender().ID, 10), r.Lang(lang)) {
		return errors.New("user language set failed")
	}
	return nil
}

func (r LanguageRepo) userLanguage(user string) string {
	if lang, ok := r.repo.Get("us_" + user); ok {
		return lang
	}
	return "en"
}

func (r LanguageRepo) setUserLanguage(user string, lang string) bool {
	return r.repo.Put("us_"+user, lang)
}

func (r LanguageRepo) chatLanguage(chat string) string {
	if lang, ok := r.repo.Get("ch_" + chat); ok {
		return lang
	}
	return "en"
}

func (r LanguageRepo) setChatLanguage(chat string, lang string) bool {
	return r.repo.Put("ch_"+chat, lang)
}
