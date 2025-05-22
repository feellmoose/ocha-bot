package mine

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestMineGameBox(t *testing.T) {
	assert.Equal(t, NumBox(1).IsMine(), false)
	assert.Equal(t, NumBox(1).IsFlagged(), false)
	assert.Equal(t, NumBox(1).IsClicked(), false)
	assert.Equal(t, NumBox(1).Clicked().IsMine(), false)
	assert.Equal(t, NumBox(1).Clicked().IsFlagged(), false)
	assert.Equal(t, NumBox(1).Clicked().IsClicked(), true)
	assert.Equal(t, NumBox(1).Clicked().Clicked().IsMine(), false)
	assert.Equal(t, NumBox(1).Clicked().Clicked().IsFlagged(), false)
	assert.Equal(t, NumBox(1).Clicked().Clicked().IsClicked(), false)
	assert.Equal(t, NumBox(1).Flagged().IsMine(), false)
	assert.Equal(t, NumBox(1).Flagged().IsFlagged(), true)
	assert.Equal(t, NumBox(1).Flagged().IsClicked(), false)
	assert.Equal(t, NumBox(1).Flagged().Flagged().IsMine(), false)
	assert.Equal(t, NumBox(1).Flagged().Flagged().IsFlagged(), false)
	assert.Equal(t, NumBox(1).Flagged().Flagged().IsClicked(), false)
	assert.Equal(t, NumBox(0).Num(), 0)
	assert.Equal(t, NumBox(1).Num(), 1)
	assert.Equal(t, NumBox(2).Num(), 2)
	assert.Equal(t, NumBox(3).Num(), 3)
	assert.Equal(t, NumBox(4).Num(), 4)
	assert.Equal(t, NumBox(5).Num(), 5)
	assert.Equal(t, NumBox(6).Num(), 6)
	assert.Equal(t, NumBox(7).Num(), 7)
	assert.Equal(t, NumBox(8).Num(), 8)
	assert.Equal(t, MineBox().IsMine(), true)
	assert.Equal(t, MineBox().IsFlagged(), false)
	assert.Equal(t, MineBox().IsClicked(), false)
}

func TestMineGameClick(t *testing.T) {
	f := Factory{}
	mine, err := f.Empty("test", 1, Additional{}, 8, 8, 10)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, mine.Status(), UnInit)
	mine, err = f.Init(mine, 4, 4)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, mine.Status(), Init)
	assert.Equal(t, mine.Boxes()[4][4].IsMine(), false)
	mine = mine.OnClicked(Position{X: 4, Y: 4}).(TelegramMineGame)
	assert.Equal(t, mine.Status(), Running)
	assert.Equal(t, mine.Boxes()[4][4].IsClicked(), true)
}
