package mine

import (
	"errors"
	"math/rand/v2"
	"time"
)

type Factory struct {
}

func (f Factory) Create(id string, user int64, info Additional, width, height, mines int) (TelegramMineGame, error) {
	boxes := make([][]Box, width)
	for i := range boxes {
		boxes[i] = make([]Box, height)
	}

	total := width * height
	indices := rand.Perm(total)

	for i := 0; i < mines; i++ {
		idx := indices[i]
		bx, by := idx/height, idx%height
		boxes[bx][by] = MineBox()
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if boxes[i][j].IsMine() {
				continue
			}
			count := 0
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					nx, ny := i+dx, j+dy
					if nx == i && ny == j {
						continue
					}
					if nx >= 0 && ny >= 0 && nx < width && ny < height && boxes[nx][ny].IsMine() {
						count++
					}
				}
			}
			boxes[i][j] = NumBox(count)
		}
	}

	boxNum := make([][]int, width)
	for i := range boxNum {
		boxNum[i] = make([]int, height)
		for j, box := range boxes[i] {
			boxNum[i][j] = box.Value
		}
	}

	now := time.Now()

	return TelegramMineGame{
		data: Serialized{
			id:        id,
			user:      user,
			infos:     info.ToMap(),
			steps:     0,
			mines:     mines,
			width:     width,
			height:    height,
			boxes:     boxNum,
			histories: make([]History, 0),
			status:    UnInit,
			create:    now,
			update:    time.Time{},
			start:     now,
			end:       time.Time{},
			win:       false,
		},
		info: info,
	}, nil
}

func (f Factory) Empty(id string, user int64, info Additional, width, height, mines int) (TelegramMineGame, error) {
	if width <= 0 || height <= 0 {
		return TelegramMineGame{}, errors.New("width <= 0 || height <= 0")
	}
	if mines < 0 {
		return TelegramMineGame{}, errors.New("mines < 0")
	}
	if width*height <= mines {
		return TelegramMineGame{}, errors.New("width * height <= mines")
	}

	return TelegramMineGame{
		data: Serialized{
			id:        id,
			user:      user,
			infos:     info.ToMap(),
			steps:     0,
			mines:     mines,
			width:     width,
			height:    height,
			boxes:     nil,
			histories: nil,
			status:    UnInit,
			create:    time.Now(),
			update:    time.Time{},
			start:     time.Time{},
			end:       time.Time{},
			win:       false,
		},
		info: info,
	}, nil
}

func (f Factory) Init(empty TelegramMineGame, x, y int) (TelegramMineGame, error) {
	game := empty.data
	if game.status != UnInit {
		return empty, nil
	}

	width, height, mines := game.width, game.height, game.mines

	boxes := make([][]Box, width)
	for i := range boxes {
		boxes[i] = make([]Box, height)
	}

	total := width * height
	indices := rand.Perm(total)

	var safeSet map[int]struct{}
	if x >= 0 && y >= 0 && x < width && y < height {
		safeSet = make(map[int]struct{})
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				nx, ny := x+dx, y+dy
				if nx >= 0 && ny >= 0 && nx < width && ny < height {
					safeSet[nx*height+ny] = struct{}{}
				}
			}
		}
	}

	filtered := make([]int, 0, total)
	for _, idx := range indices {
		if safeSet == nil || safeSet[idx] == struct{}{} {
			continue
		}
		filtered = append(filtered, idx)
	}

	if len(filtered) < mines {
		return f.Create(game.id, game.user, empty.Infos(), game.width, game.height, game.mines)
	}

	for i := 0; i < mines; i++ {
		idx := filtered[i]
		bx, by := idx/height, idx%height
		boxes[bx][by] = MineBox()
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if boxes[i][j].IsMine() {
				continue
			}
			count := 0
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					nx, ny := i+dx, j+dy
					if nx == i && ny == j {
						continue
					}
					if nx >= 0 && ny >= 0 && nx < width && ny < height && boxes[nx][ny].IsMine() {
						count++
					}
				}
			}
			boxes[i][j] = NumBox(count)
		}
	}

	boxNum := make([][]int, width)
	for i := range boxNum {
		boxNum[i] = make([]int, height)
		for j, box := range boxes[i] {
			boxNum[i][j] = box.Value
		}
	}

	now := time.Now()

	return TelegramMineGame{
		data: Serialized{
			id:        game.id,
			user:      game.user,
			infos:     empty.info.ToMap(),
			steps:     0,
			mines:     mines,
			width:     width,
			height:    height,
			boxes:     boxNum,
			histories: make([]History, 0),
			status:    UnInit,
			create:    now,
			update:    time.Time{},
			start:     now,
			end:       time.Time{},
			win:       false,
		},
		info: empty.info,
	}, nil
}
