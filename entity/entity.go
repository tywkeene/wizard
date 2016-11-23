package entity

import "github.com/tywkeene/wizard/position"

type Entity interface {
	GetPosition() *position.Position
	GetSymbol() rune
	GetName() string
	Move(x int, y int)
}
