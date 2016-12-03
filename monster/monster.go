package monster

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/item"
	"github.com/tywkeene/wizard/position"
)

type Inventory struct {
	List map[*item.Item]int
}

type Monster struct {
	Name     string
	ID       int
	Position *position.Position
	Items    *Inventory
	Symbol   rune
	Passable bool
	Type     int
}

func (m *Monster) PickupItem(i *item.Item) {
	m.Items.List[i]++
}

func (m *Monster) GetName() string {
	return m.Name
}

func (m *Monster) SetID(id int) {
	m.ID = id
}

func (m *Monster) GetID() int {
	return m.ID
}

func (m *Monster) GetPosition() *position.Position {
	return m.Position
}

func (m *Monster) GetSymbol() rune {
	return m.Symbol
}

func (m *Monster) IsPassable() bool {
	return m.Passable
}

func (m *Monster) GetType() int {
	return m.Type
}

func (m *Monster) Move(x int, y int) {
	m.Position.PrevX = m.Position.X
	m.Position.PrevY = m.Position.Y
	m.Position.X = x
	m.Position.Y = y
}

func (m *Monster) Draw() {
	p := m.GetPosition()
	termbox.SetCell(p.X, p.Y, m.GetSymbol(), termbox.ColorWhite, termbox.ColorBlack)
}

func NewInventory() *Inventory {
	return &Inventory{make(map[*item.Item]int)}
}

func MakeMonster(name string, symbol rune, monsterType int) *Monster {
	return &Monster{
		Name:     name,
		ID:       -1,
		Position: position.NewPosition(-1, -1, -1, -1, 1, 1),
		Items:    NewInventory(),
		Symbol:   symbol,
		Passable: false,
		Type:     monsterType,
	}
}
