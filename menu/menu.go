package menu

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/position"
	"github.com/tywkeene/wizard/state"
	"log"
)

type OptionHandler func(*state.GameState, *Menu) int
type InputHandler func(*state.GameState, *Menu)

type MenuOption struct {
	Description string
	Do          OptionHandler
}

type Menu struct {
	Pos     *position.Position
	Options []*MenuOption
	Handle  InputHandler
}

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

func (m *Menu) Execute(s *state.GameState) {
	m.Draw(1)
	m.Handle(s, m)
}

//Start Menu Handles
func EmptyOptionHandle(s *state.GameState, m *Menu) int {
	return 0
}
func StartMenuStartGame(s *state.GameState, m *Menu) int {
	s.Running = true
	return 1
}

func StartMenuExitGame(s *state.GameState, m *Menu) int {
	s.Running = false
	return 1
}

func StartMenuInputHandle(s *state.GameState, m *Menu) {
	var leaveMenu = false
	var cursorPos = 1
	for leaveMenu == false {
		select {
		case ev := <-s.Events:
			switch {
			case ev.Ch == 'k': //up
				if cursorPos > 1 {
					cursorPos--
				}
				break
			case ev.Ch == 'j': //down
				if cursorPos < len(m.Options) {
					cursorPos++
				}
				break
			case ev.Key == termbox.KeyEnter:
				if ret := m.Options[cursorPos].Do(s, m); ret == 1 {
					log.Printf("Selected menu option %d", cursorPos)
					leaveMenu = true
				}
				break
			}
		}
		m.Draw(cursorPos)
	}
}
