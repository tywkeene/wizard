package item

import (
	"github.com/tywkeene/wizard/entity"
	"github.com/tywkeene/wizard/position"
	"log"
)

type ItemApply func(e entity.Entity)

type Item struct {
	Pos    *position.Position
	Name   string
	Symbol rune
	Apply  ItemApply
}

var (
	AllCoins = []*Item{
		&Item{nil, "Copper coin", '$', func(e entity.Entity) {
			log.Println("You find no real use for the copper coin")
		}},
	}
	AllScrolls = []*Item{
		&Item{nil, "Scroll of teleportation", '?', func(e entity.Entity) {
			log.Println("You read the scroll of teleportation")
		}},
	}
)

func NewItem(pos *position.Position, name string, symbol rune, apply ItemApply) *Item {
	return &Item{pos, name, symbol, apply}
}

func (i *Item) GetPosition() *position.Position {
	return i.Pos
}

func (i *Item) GetSymbol() rune {
	return i.Symbol
}

func (i *Item) GetName() string {
	return i.Name
}

func (i *Item) Use(e entity.Entity) {
	i.Apply(e)
}
