package mine

import (
	"errors"
	"log"
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
			ID:        id,
			User:      user,
			Infos:     info.ToMap(),
			Steps:     0,
			Mines:     mines,
			Width:     width,
			Height:    height,
			Boxes:     boxNum,
			Histories: make([]History, 0),
			Status:    Init,
			Create:    now,
			Update:    time.Time{},
			Start:     now,
			End:       time.Time{},
			Win:       false,
		},
		info: info,
	}, nil
}

func (f Factory) Empty(id string, user int64, info Additional, width, height, mines int) (TelegramMineGame, error) {
	if width <= 0 || height <= 0 {
		return TelegramMineGame{}, errors.New("Width <= 0 || Height <= 0")
	}
	if mines < 0 {
		return TelegramMineGame{}, errors.New("Mines < 0")
	}
	if width*height <= mines {
		return TelegramMineGame{}, errors.New("Width * Height <= Mines")
	}

	log.Printf("Create mine game (...info=%v)", info)

	return TelegramMineGame{
		data: Serialized{
			ID:        id,
			User:      user,
			Infos:     info.ToMap(),
			Steps:     0,
			Mines:     mines,
			Width:     width,
			Height:    height,
			Boxes:     nil,
			Histories: nil,
			Status:    UnInit,
			Create:    time.Now(),
			Update:    time.Time{},
			Start:     time.Time{},
			End:       time.Time{},
			Win:       false,
		},
		info: info,
	}, nil
}

func (f Factory) Init(empty TelegramMineGame, x, y int) (TelegramMineGame, error) {
	game := empty.data
	if game.Status != UnInit {
		return empty, nil
	}

	width, height, mines := game.Width, game.Height, game.Mines

	boxes := make([][]Box, width)
	for i := range boxes {
		boxes[i] = make([]Box, height)
	}

	total := width * height
	indices := rand.Perm(total)

	safe := make(map[int]struct{})
	if x >= 0 && y >= 0 && x < width && y < height {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				nx, ny := x+dx, y+dy
				if nx >= 0 && ny >= 0 && nx < width && ny < height {
					safe[nx*height+ny] = struct{}{}
				}
			}
		}
	}

	filtered := make([]int, 0, total)
	for _, idx := range indices {
		if _, ok := safe[idx]; !ok {
			filtered = append(filtered, idx)
		}
	}

	if len(filtered) < mines {
		return f.Create(game.ID, game.User, empty.Infos(), game.Width, game.Height, game.Mines)
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
			ID:        game.ID,
			User:      game.User,
			Infos:     empty.info.ToMap(),
			Steps:     0,
			Mines:     mines,
			Width:     width,
			Height:    height,
			Boxes:     boxNum,
			Histories: make([]History, 0),
			Status:    Init,
			Create:    now,
			Update:    time.Time{},
			Start:     now,
			End:       time.Time{},
			Win:       false,
		},
		info: empty.info,
	}, nil
}
