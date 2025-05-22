package mine

import (
	"gopkg.in/telebot.v4"
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
}

type Serialized struct {
	id        string
	user      int64
	infos     map[string]string
	steps     int
	mines     int
	width     int
	height    int
	boxes     [][]int
	histories []History
	status    GameStatus
	create    time.Time
	update    time.Time
	start     time.Time
	end       time.Time
	win       bool
}

func (s Serialized) Deserialize() Mine {
	infos, _ := FromMap(s.infos)
	return TelegramMineGame{
		data: s,
		info: infos,
	}
}

// Position for each steps
type Position struct {
	X int
	Y int
}

func (p Position) InBounds(width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

// History for rollback
type History struct {
	Pos     Position
	Option  GameOption
	Updated time.Time
	Related []History
}

// GameStatus for steps and win check
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

// Box is a no status mine unit
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
	RandomBottom  GameType = "Random_Bottom"
	LevelBottom   GameType = "Level_Bottom"
)

type Button string

const (
	BClick Button = "Click"
	BFlag  Button = "Flag"
)

type Additional struct {
	Type    GameType
	Button  Button
	Locale  string
	Topic   int
	Chat    int64
	Message int
}

func (a Additional) ToMap() map[string]string {
	res := map[string]string{
		"type":   string(a.Type),
		"button": string(a.Button),
		"locale": a.Locale,
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

	return Additional{
		Type:    GameType(m["type"]),
		Button:  Button(m["button"]),
		Locale:  m["locale"],
		Topic:   topic,
		Chat:    chat,
		Message: message,
	}, nil
}
