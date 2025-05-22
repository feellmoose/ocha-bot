package mine

import (
	"time"
)

type TelegramMineGame struct {
	data Serialized
	info Additional
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
	if game.status == UnInit {
		return t
	}

	box := Box{game.boxes[pos.X][pos.Y]}

	if game.win || box.IsClicked() || box.IsFlagged() || game.status == End {
		return t
	}

	clicked := 1
	newBoxes := CloneBoxes(game.boxes)
	newBoxes[pos.X][pos.Y] = box.Clicked().Value

	now := time.Now()
	newHistory := append(game.histories, History{
		Pos:     pos,
		Option:  Click,
		Updated: now,
	})

	if box.IsMine() {
		newHistory[len(newHistory)-1].Option = Boom
		return TelegramMineGame{
			data: Serialized{
				steps:     game.steps + clicked,
				histories: newHistory,
				boxes:     newBoxes,
				status:    End,
				update:    now,
				end:       now,
				win:       false,
				infos:     t.info.ToMap(),
				id:        game.id,
				user:      game.user,
				mines:     game.mines,
				width:     game.width,
				height:    game.height,
				start:     game.start,
				create:    game.create,
			},
			info: t.info,
		}
	}

	if box.Num() == 0 {
		related := clickedZero(game.width, game.height, pos, newBoxes)
		clicked += len(related)
		newHistory[len(newHistory)-1].Related = related

		for _, h := range related {
			p := h.Pos
			newBoxes[p.X][p.Y] = Box{game.boxes[p.X][p.Y]}.Clicked().Value
		}
	}

	if game.steps+clicked+game.mines == game.width*game.height {
		return TelegramMineGame{
			data: Serialized{
				steps:     game.steps + clicked,
				histories: newHistory,
				boxes:     newBoxes,
				status:    End,
				update:    now,
				end:       now,
				win:       true,
				infos:     t.info.ToMap(),
				id:        game.id,
				user:      game.user,
				mines:     game.mines,
				width:     game.width,
				height:    game.height,
				start:     game.start,
				create:    game.create,
			},
			info: t.info,
		}
	}

	return TelegramMineGame{
		data: Serialized{
			steps:     game.steps + clicked,
			histories: newHistory,
			boxes:     newBoxes,
			status:    Running,
			update:    now,
			end:       time.Time{},
			win:       false,
			infos:     t.info.ToMap(),
			id:        game.id,
			user:      game.user,
			mines:     game.mines,
			width:     game.width,
			height:    game.height,
			start:     game.start,
			create:    game.create,
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
	if !pos.InBounds(game.width, game.height) || game.win || game.status != Running {
		return t
	}

	box := Box{game.boxes[pos.X][pos.Y]}

	if box.IsClicked() {
		return t
	}

	now := time.Now()

	newHistory := append(game.histories, History{
		Pos:     pos,
		Option:  Flag,
		Updated: now,
	})

	newBoxes := CloneBoxes(game.boxes)
	newBoxes[pos.X][pos.Y] = box.Flagged().Value
	return TelegramMineGame{
		data: Serialized{
			steps:     game.steps,
			histories: newHistory,
			boxes:     newBoxes,
			status:    Running,
			update:    now,
			end:       time.Time{},
			win:       false,
			infos:     t.info.ToMap(),
			id:        game.id,
			user:      game.user,
			mines:     game.mines,
			width:     game.width,
			height:    game.height,
			start:     game.start,
			create:    game.create,
		},
		info: t.info,
	}
}

func (t TelegramMineGame) OnRollback(s int) Mine {
	game := t.data
	if game.status == UnInit {
		return t
	}

	now := time.Now()

	var newHistory []History
	if len(game.histories) > s {
		newHistory = game.histories[:len(game.histories)-s]
	} else {
		newHistory = []History{}
	}

	newBoxes := CloneBoxes(game.boxes)
	step := 0

	for i := len(game.histories) - 1; i >= len(game.histories)-s && i >= 0; i-- {
		h := game.histories[i]
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
			steps:     game.steps - step,
			histories: newHistory,
			boxes:     newBoxes,
			status:    Running,
			update:    now,
			end:       time.Time{},
			win:       false,
			infos:     t.info.ToMap(),
			id:        game.id,
			user:      game.user,
			mines:     game.mines,
			width:     game.width,
			height:    game.height,
			start:     game.start,
			create:    game.create,
		},
		info: t.info,
	}
}

func (t TelegramMineGame) Serialize() Serialized {
	game := &t.data
	return Serialized{
		id:        game.id,
		user:      game.user,
		infos:     t.info.ToMap(),
		steps:     game.steps,
		mines:     game.mines,
		width:     game.width,
		height:    game.height,
		boxes:     game.boxes,
		histories: game.histories,
		status:    game.status,
		create:    game.create,
		update:    game.update,
		start:     game.start,
		end:       game.end,
		win:       game.win,
	}
}

func (t TelegramMineGame) ID() string {
	return t.data.id
}

func (t TelegramMineGame) UserID() int64 {
	return t.data.user
}

func (t TelegramMineGame) Infos() Additional {
	return t.info
}

func (t TelegramMineGame) Steps() int {
	return t.data.steps
}

func (t TelegramMineGame) Mines() int {
	return t.data.mines
}

func (t TelegramMineGame) Width() int {
	return t.data.width
}

func (t TelegramMineGame) Height() int {
	return t.data.height
}

func (t TelegramMineGame) Boxes() [][]Box {
	res := make([][]Box, t.Width())
	for i := range res {
		res[i] = make([]Box, t.Height())
	}

	for i, row := range t.data.boxes {
		for j, val := range row {
			res[i][j] = Box{val}
		}
	}

	return res
}

func (t TelegramMineGame) History() []History {
	return t.data.histories
}

func (t TelegramMineGame) Status() GameStatus {
	return t.data.status
}

func (t TelegramMineGame) Duration() time.Duration {
	game := t.data
	switch game.status {
	case UnInit, Init:
		return 0
	case Running:
		start := game.start
		if start.IsZero() {
			start = time.Now()
		}
		return time.Since(start)
	case End:
		start := game.start
		if start.IsZero() {
			start = time.Now()
		}
		end := game.end
		if end.IsZero() {
			end = time.Now()
		}
		return end.Sub(start)
	default:
		return 0
	}
}

func (t TelegramMineGame) Win() bool {
	return t.data.win
}
