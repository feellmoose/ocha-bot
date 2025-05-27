package mine

import (
	"time"
)

type TelegramMineGame struct {
	data Serialized
	info Additional
}
type TelegramMineGameScore struct {
	Username string  `json:"username,omitempty"`
	Time     string  `json:"time,omitempty"`
	Duration int64   `json:"duration,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Steps    int     `json:"steps,omitempty"`
	Mines    int     `json:"mines,omitempty"`
	Width    int     `json:"width,omitempty"`
	Height   int     `json:"height,omitempty"`
}

func (t TelegramMineGame) Score() TelegramMineGameScore {
	var (
		width    = t.data.Width
		height   = t.data.Height
		mines    = t.data.Mines
		steps    = t.data.Steps
		duration = t.Duration().Milliseconds()
		safe     = width*height - mines
		baseTime = float64(safe) * 0.5
	)

	totalCells := float64(width * height)
	diff := float64(mines) / totalCells
	const kd = 0.4
	diffScore := diff / (diff + kd)

	T := t.Duration().Seconds()
	timeScore := baseTime / (baseTime + T)

	return TelegramMineGameScore{
		Username: t.info.Username,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Duration: duration,
		Score:    100 * (diffScore + timeScore) / 2,
		Steps:    steps,
		Mines:    mines,
		Width:    width,
		Height:   height,
	}
}

func (t TelegramMineGame) OnInfoChanged(additional Additional) Mine {
	return TelegramMineGame{
		data: t.data,
		info: additional,
	}
}

func (t TelegramMineGame) OnClicked(pos Position) Mine {
	game := &t.data
	if !pos.InBounds(t.Width(), t.Height()) {
		return t
	}
	if game.Status == UnInit {
		return t
	}

	box := Box{game.Boxes[pos.X][pos.Y]}

	if game.Win || box.IsClicked() || box.IsFlagged() || game.Status == End {
		return t
	}

	clicked := 1
	newBoxes := CloneBoxes(game.Boxes)
	newBoxes[pos.X][pos.Y] = box.Clicked().Value

	now := time.Now()
	newHistory := append(game.Histories, History{
		Pos:     pos,
		Option:  Click,
		Updated: now,
	})

	if box.IsMine() {
		newHistory[len(newHistory)-1].Option = Boom
		return TelegramMineGame{
			data: Serialized{
				Steps:     game.Steps + clicked,
				Histories: newHistory,
				Boxes:     newBoxes,
				Status:    End,
				Update:    now,
				End:       now,
				Win:       false,
				Infos:     t.info.ToMap(),
				ID:        game.ID,
				User:      game.User,
				Mines:     game.Mines,
				Width:     game.Width,
				Height:    game.Height,
				Start:     game.Start,
				Create:    game.Create,
			},
			info: t.info,
		}
	}

	if box.Num() == 0 {
		related := clickedZero(game.Width, game.Height, pos, newBoxes)
		clicked += len(related)
		newHistory[len(newHistory)-1].Related = related

		for _, h := range related {
			p := h.Pos
			newBoxes[p.X][p.Y] = Box{game.Boxes[p.X][p.Y]}.Clicked().Value
		}
	}

	if game.Steps+clicked+game.Mines == game.Width*game.Height {
		return TelegramMineGame{
			data: Serialized{
				Steps:     game.Steps + clicked,
				Histories: newHistory,
				Boxes:     newBoxes,
				Status:    End,
				Update:    now,
				End:       now,
				Win:       true,
				Infos:     t.info.ToMap(),
				ID:        game.ID,
				User:      game.User,
				Mines:     game.Mines,
				Width:     game.Width,
				Height:    game.Height,
				Start:     game.Start,
				Create:    game.Create,
			},
			info: t.info,
		}
	}

	return TelegramMineGame{
		data: Serialized{
			Steps:     game.Steps + clicked,
			Histories: newHistory,
			Boxes:     newBoxes,
			Status:    Running,
			Update:    now,
			End:       time.Time{},
			Win:       false,
			Infos:     t.info.ToMap(),
			ID:        game.ID,
			User:      game.User,
			Mines:     game.Mines,
			Width:     game.Width,
			Height:    game.Height,
			Start:     game.Start,
			Create:    game.Create,
		},
		info: t.info,
	}
}

func neighbors(p Position) []Position {
	return []Position{
		{p.X - 1, p.Y - 1},
		{p.X, p.Y - 1},
		{p.X + 1, p.Y - 1},
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X - 1, p.Y + 1},
		{p.X, p.Y + 1},
		{p.X + 1, p.Y + 1},
	}
}

func filterPositions(positions []Position, width, height int) []Position {
	var result []Position
	for _, p := range positions {
		if p.InBounds(width, height) {
			result = append(result, p)
		}
	}
	return result
}

func CloneBoxes(boxes [][]int) [][]int {
	width := len(boxes)
	height := len(boxes[0])
	clone := make([][]int, width)
	for i := 0; i < width; i++ {
		clone[i] = make([]int, height)
		copy(clone[i], boxes[i])
	}
	return clone
}

func clickedZero(width, height int, from Position, boxes [][]int) []History {
	visited := map[Position]struct{}{from: {}}
	var history []History
	require := neighbors(from)

	filtered := filterPositions(require, width, height)
	for len(filtered) > 0 {
		var next []Position
		for _, p := range filtered {
			if _, ok := visited[p]; ok {
				continue
			}
			b := Box{boxes[p.X][p.Y]}
			if b.IsClicked() || b.IsFlagged() || b.IsMine() {
				continue
			}
			history = append(history, History{Pos: p, Option: Click})
			visited[p] = struct{}{}
			if b.Num() == 0 {
				next = append(next, neighbors(p)...)
			}
		}
		filtered = filterPositions(next, width, height)
		for i := 0; i < len(filtered); {
			if _, ok := visited[filtered[i]]; ok {
				filtered = append(filtered[:i], filtered[i+1:]...)
			} else {
				i++
			}
		}
	}
	return history
}

func (t TelegramMineGame) OnFlagged(pos Position) Mine {

	game := t.data
	if !pos.InBounds(game.Width, game.Height) || game.Win || game.Status != Running {
		return t
	}

	box := Box{game.Boxes[pos.X][pos.Y]}

	if box.IsClicked() {
		return t
	}

	now := time.Now()

	newHistory := append(game.Histories, History{
		Pos:     pos,
		Option:  Flag,
		Updated: now,
	})

	newBoxes := CloneBoxes(game.Boxes)
	newBoxes[pos.X][pos.Y] = box.Flagged().Value
	return TelegramMineGame{
		data: Serialized{
			Steps:     game.Steps,
			Histories: newHistory,
			Boxes:     newBoxes,
			Status:    Running,
			Update:    now,
			End:       time.Time{},
			Win:       false,
			Infos:     t.info.ToMap(),
			ID:        game.ID,
			User:      game.User,
			Mines:     game.Mines,
			Width:     game.Width,
			Height:    game.Height,
			Start:     game.Start,
			Create:    game.Create,
		},
		info: t.info,
	}
}

func (t TelegramMineGame) OnRollback(s int) Mine {
	game := t.data
	if game.Status == UnInit {
		return t
	}

	now := time.Now()

	var newHistory []History
	if len(game.Histories) > s {
		newHistory = game.Histories[:len(game.Histories)-s]
	} else {
		newHistory = []History{}
	}

	newBoxes := CloneBoxes(game.Boxes)
	step := 0

	for i := len(game.Histories) - 1; i >= len(game.Histories)-s && i >= 0; i-- {
		h := game.Histories[i]
		pos := h.Pos
		switch h.Option {
		case Flag:
			box := Box{newBoxes[pos.X][pos.Y]}
			newBoxes[pos.X][pos.Y] = box.Flagged().Value

		case Click:
			box := Box{newBoxes[pos.X][pos.Y]}
			newBoxes[pos.X][pos.Y] = NumBox(box.Num()).Value
			step++
			if h.Related != nil {
				step += len(h.Related)
				for _, rel := range h.Related {
					r := rel.Pos
					b := Box{newBoxes[r.X][r.Y]}
					newBoxes[r.X][r.Y] = NumBox(b.Num()).Value
				}
			}

		case Boom:
			step++
			newBoxes[pos.X][pos.Y] = MineBox().Value
		}
	}
	return TelegramMineGame{
		data: Serialized{
			Steps:     game.Steps - step,
			Histories: newHistory,
			Boxes:     newBoxes,
			Status:    Running,
			Update:    now,
			End:       time.Time{},
			Win:       false,
			Infos:     t.info.ToMap(),
			ID:        game.ID,
			User:      game.User,
			Mines:     game.Mines,
			Width:     game.Width,
			Height:    game.Height,
			Start:     game.Start,
			Create:    game.Create,
		},
		info: t.info,
	}
}

func (t TelegramMineGame) Serialize() Serialized {
	game := &t.data
	return Serialized{
		ID:        game.ID,
		User:      game.User,
		Infos:     t.info.ToMap(),
		Steps:     game.Steps,
		Mines:     game.Mines,
		Width:     game.Width,
		Height:    game.Height,
		Boxes:     game.Boxes,
		Histories: game.Histories,
		Status:    game.Status,
		Create:    game.Create,
		Update:    game.Update,
		Start:     game.Start,
		End:       game.End,
		Win:       game.Win,
	}
}

func (t TelegramMineGame) ID() string {
	return t.data.ID
}

func (t TelegramMineGame) UserID() int64 {
	return t.data.User
}

func (t TelegramMineGame) Infos() Additional {
	return t.info
}

func (t TelegramMineGame) Steps() int {
	return t.data.Steps
}

func (t TelegramMineGame) Mines() int {
	return t.data.Mines
}

func (t TelegramMineGame) Width() int {
	return t.data.Width
}

func (t TelegramMineGame) Height() int {
	return t.data.Height
}

func (t TelegramMineGame) Boxes() [][]Box {
	res := make([][]Box, t.Width())
	for i := range res {
		res[i] = make([]Box, t.Height())
	}

	for i, row := range t.data.Boxes {
		for j, val := range row {
			res[i][j] = Box{val}
		}
	}

	return res
}

func (t TelegramMineGame) History() []History {
	return t.data.Histories
}

func (t TelegramMineGame) Status() GameStatus {
	return t.data.Status
}

func (t TelegramMineGame) Duration() time.Duration {
	game := t.data
	switch game.Status {
	case UnInit, Init:
		return 0
	case Running:
		start := game.Start
		if start.IsZero() {
			start = time.Now()
		}
		return time.Since(start)
	case End:
		start := game.Start
		if start.IsZero() {
			start = time.Now()
		}
		end := game.End
		if end.IsZero() {
			end = time.Now()
		}
		return end.Sub(start)
	default:
		return 0
	}
}

func (t TelegramMineGame) Win() bool {
	return t.data.Win
}
