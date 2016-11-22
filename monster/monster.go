package monster

import (
	"github.com/tywkeene/wizard/position"
	"log"
)

type Monster struct {
	Pos    *position.Position
	Name   string
	Symbol rune
	Index  int
}

func MakeMonster(x int, y int, name string, symbol rune, index int) *Monster {
	log.Printf("Initialized monster at position [%d:%d]: (%s) (%s)", x, y, name, string(symbol))
	pos := position.NewPosition(-1, -1, x, y, 1, 1)
	return &Monster{Pos: pos, Name: name, Symbol: symbol, Index: index}
}

func (m *Monster) Move(newX int, newY int) {
	m.Pos.PrevX = m.Pos.X
	m.Pos.PrevY = m.Pos.Y
	m.Pos.X = newX
	m.Pos.Y = newY
}

func (m *Monster) GetPosition() *position.Position {
	return m.Pos
}

func (m *Monster) GetSymbol() rune {
	return m.Symbol
}

func (m *Monster) GetName() string {
	return m.Name
}

func (m *Monster) GetIndex() int {
	return m.Index
}
