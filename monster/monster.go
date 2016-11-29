package monster

import (
	"github.com/tywkeene/wizard/position"
	"log"
)

type Monster struct {
	Pos    *position.Position
	Name   string
	Symbol rune
}

func MakeMonster(name string, symbol rune) *Monster {
	log.Printf("Initialized monster: (%s) (%s)", name, string(symbol))
	pos := position.NewPosition(-1, -1, -1, -1, 1, 1)
	return &Monster{Pos: pos, Name: name, Symbol: symbol}
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
