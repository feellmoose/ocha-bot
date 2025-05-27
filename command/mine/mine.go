package mine

import (
	"gopkg.in/telebot.v4"
	"ocha_server_bot/helper"
	"strconv"
	"time"
)

type Mine interface {
	ID() string
	UserID() int64
	Steps() int
	Mines() int
	Width() int
	Height() int
	Boxes() [][]Box
	History() []History
	Status() GameStatus
	Duration() time.Duration
	Infos() Additional
	Win() bool

	OnClicked(pos Position) Mine
	OnFlagged(pos Position) Mine
	OnRollback(steps int) Mine
	OnInfoChanged(additional Additional) Mine

	Serialize() Serialized
	Display(c telebot.Context) error
	RankDisplay(c telebot.Context, ranker helper.Ranker[TelegramMineGameScore]) error
}

type Serialized struct {
	ID        string            `json:"id,omitempty"`
	User      int64             `json:"user,omitempty"`
	Infos     map[string]string `json:"infos,omitempty"`
	Steps     int               `json:"steps,omitempty"`
	Mines     int               `json:"mines,omitempty"`
	Width     int               `json:"width,omitempty"`
	Height    int               `json:"height,omitempty"`
	Boxes     [][]int           `json:"boxes,omitempty"`
	Histories []History         `json:"histories,omitempty"`
	Status    GameStatus        `json:"status,omitempty"`
	Create    time.Time         `json:"create,omitempty"`
	Update    time.Time         `json:"update,omitempty"`
	Start     time.Time         `json:"start,omitempty"`
	End       time.Time         `json:"end,omitempty"`
	Win       bool              `json:"win,omitempty"`
}

func (s Serialized) Deserialize() Mine {
	infos, _ := FromMap(s.Infos)
	return TelegramMineGame{
		data: s,
		info: infos,
	}
}

// Position for each Steps
type Position struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

func (p Position) InBounds(width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

// History for rollback
type History struct {
	Pos     Position   `json:"pos,omitempty"`
	Option  GameOption `json:"option,omitempty"`
	Updated time.Time  `json:"updated,omitempty"`
	Related []History  `json:"related,omitempty"`
}

// GameStatus for Steps and Win check
type GameStatus int

const (
	UnInit GameStatus = iota
	Init
	Running
	End
)

// GameOption for history
type GameOption int

const (
	Click GameOption = iota
	Flag
	Boom
)

// Box is a no Status mine unit
type Box struct {
	Value int
}

func NewBox(value int) Box {
	return Box{Value: value}
}

func NumBox(value int) Box {
	return NewBox(value)
}

func MineBox() Box {
	return NewBox(0x10000)
}

func (b Box) Num() int {
	return b.Value & 0xFF
}

func (b Box) IsFlagged() bool {
	return b.Value&0x1000 != 0
}

func (b Box) IsClicked() bool {
	return b.Value&0x100000 != 0
}

func (b Box) IsMine() bool {
	return b.Value&0x10000 != 0
}

func (b Box) Flagged() Box {
	if b.IsFlagged() {
		return Box{Value: b.Value &^ 0x1000}
	}
	return Box{Value: b.Value | 0x1000}
}

func (b Box) Clicked() Box {
	if b.IsClicked() {
		return Box{Value: b.Value &^ 0x100000}
	}
	return Box{Value: b.Value | 0x100000}
}

type GameType string

const (
	ClassicBottom GameType = "Classic_Bottom"
	Classic       GameType = "c"
	Rank          GameType = "r"
)

type Button string

const (
	BClick Button = "Click"
	BFlag  Button = "Flag"
)

type Additional struct {
	Type     GameType
	Button   Button
	Locale   string
	Username string
	Topic    int
	Chat     int64
	Message  int
}

func (a Additional) ToMap() map[string]string {
	res := map[string]string{
		"type":   string(a.Type),
		"button": string(a.Button),
		"locale": a.Locale,
	}
	if a.Username != "" {
		res["username"] = a.Username
	}
	if a.Topic != 0 {
		res["topic"] = strconv.Itoa(a.Topic)
	}
	if a.Chat != 0 {
		res["chat"] = strconv.FormatInt(a.Chat, 10)
	}
	if a.Message != 0 {
		res["message"] = strconv.Itoa(a.Message)
	}
	return res
}

func FromMap(m map[string]string) (Additional, error) {
	chat, topic, message := int64(0), 0, 0

	if val, ok := m["topic"]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			topic = v
		}
	}
	if val, ok := m["chat"]; ok {
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			chat = v
		}
	}
	if val, ok := m["message"]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			message = v
		}
	}
	username := m["username"]
	return Additional{
		Type:     GameType(m["type"]),
		Button:   Button(m["button"]),
		Locale:   m["locale"],
		Topic:    topic,
		Chat:     chat,
		Message:  message,
		Username: username,
	}, nil
}
