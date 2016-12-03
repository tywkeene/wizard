package menu

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/position"
)

type OptionHandler func(*Menu) int
type InputHandler func(chan termbox.Event, *Menu) int

type MenuOption struct {
	Description string
	Do          OptionHandler
}

type Menu struct {
	Pos     *position.Position
	Options []*MenuOption
	Handle  InputHandler
}

const (
	OptionExit = iota
	OptionStartGame
)

func NewMenu(x int, y int, width int, height int, handle InputHandler) *Menu {
	pos := position.NewPosition(-1, -1, x, y, width, height)
	return &Menu{pos, make([]*MenuOption, 0), handle}
}

func NewMenuOption(description string, handle OptionHandler) *MenuOption {
	return &MenuOption{description, handle}
}

func (m *Menu) AddOption(option string, handle OptionHandler) {
	o := NewMenuOption(option, handle)
	m.Options = append(m.Options, o)
}

func (m *Menu) Draw(hilight int) {
	cursor := position.NewPosition(-1, -1, m.Pos.X, m.Pos.Y, 1, 1)
	foreground := termbox.ColorWhite
	background := termbox.ColorBlack
	for _, message := range m.Options {
		for x, b := range message.Description {
			if cursor.Y == hilight {
				foreground = termbox.ColorBlack
				background = termbox.ColorWhite
			}
			termbox.SetCell(x, cursor.Y, rune(b), foreground, background)
			termbox.Flush()
		}
		foreground = termbox.ColorWhite
		background = termbox.ColorBlack
		cursor.Y++
	}
}

func (m *Menu) Execute(events chan termbox.Event) int {
	m.Draw(1)
	return m.Handle(events, m)
}

func EmptyOptionHandle(m *Menu) int {
	return 0
}

//Start Menu Handles
func StartMenuStartGame(m *Menu) int {
	return OptionStartGame
}

func StartMenuExitGame(m *Menu) int {
	return OptionExit
}

func StartMenuInputHandle(events chan termbox.Event, m *Menu) int {
	var leaveMenu bool = false
	var ret int = OptionExit
	var cursorPos int = 1

	for leaveMenu == false {
		select {
		case ev := <-events:
			switch {
			case ev.Ch == 'k': //up
				if cursorPos > 1 {
					cursorPos--
				}
				break
			case ev.Ch == 'j': //down
				if cursorPos < len(m.Options)-1 {
					cursorPos++
				}
				break
			case ev.Key == termbox.KeyEnter:
				ret = m.Options[cursorPos].Do(m)
				leaveMenu = true
				break
			}
		}
		m.Draw(cursorPos)
	}
	return ret
}
