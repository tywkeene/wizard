package menu

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/position"
	"github.com/tywkeene/wizard/state"
	"log"
)

type OptionHandler func(*state.GameState, *Menu, ...interface{}) int
type InputHandler func(*state.GameState, *Menu) int

type MenuOption struct {
	Description string
	Do          OptionHandler
	ID          int
}

type Menu struct {
	Position *position.Position
	Options  []*MenuOption
	Handle   InputHandler
}

func NewMenu(x int, y int, width int, height int, handle InputHandler) *Menu {
	pos := position.NewPosition(-1, -1, x, y, width, height)
	m := &Menu{
		Position: pos,
		Options:  make([]*MenuOption, 0),
		Handle:   handle,
	}
	return m
}

func NewMenuOption(description string, do OptionHandler, id int) *MenuOption {
	return &MenuOption{
		Description: description,
		Do:          do,
		ID:          id,
	}
}

func (m *Menu) AddOption(description string, do OptionHandler) {
	o := NewMenuOption(description, do, len(m.Options)+1)
	m.Options = append(m.Options, o)
	log.Printf("Added option '%d' with description '%s'", len(m.Options), description)
}

func (m *Menu) DrawBorder() {
	const (
		SideBorder        = '│'
		TopBorder         = '─'
		TopLeftBorder     = '┌'
		TopRightBorder    = '┐'
		BottomLeftBorder  = '└'
		BottomRightBorder = '┘'
	)
	topLeft := m.Position.TopLeftCornerPosition()
	termbox.SetCell(topLeft.X, topLeft.Y, TopLeftBorder, termbox.ColorWhite, termbox.ColorBlack)

	topRight := m.Position.TopRightCornerPosition()
	termbox.SetCell(topRight.X, topRight.Y, TopRightBorder, termbox.ColorWhite, termbox.ColorBlack)

	bottomLeft := m.Position.BottomLeftCornerPosition()
	termbox.SetCell(bottomLeft.X, bottomLeft.Y, BottomLeftBorder, termbox.ColorWhite, termbox.ColorBlack)

	bottomRight := m.Position.BottomRightCornerPosition()
	termbox.SetCell(bottomRight.X, bottomRight.Y, BottomRightBorder, termbox.ColorWhite, termbox.ColorBlack)

	//Top
	for x := m.Position.X; x < (m.Position.X + m.Position.Width); x++ {
		termbox.SetCell(x, m.Position.Y-1, TopBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Right
	for y := m.Position.Y; y < (m.Position.Y + m.Position.Height); y++ {
		termbox.SetCell(m.Position.X+m.Position.Width, y, SideBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Bottom
	for x := m.Position.X; x < (m.Position.X + m.Position.Width); x++ {
		termbox.SetCell(x, m.Position.Y+m.Position.Height, TopBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Left
	for y := m.Position.Y; y < (m.Position.Y + m.Position.Height); y++ {
		termbox.SetCell(m.Position.X-1, y, SideBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}

func (m *Menu) Draw(hilite int) {
	cursor := position.NewPosition(-1, -1, m.Position.X, m.Position.Y, 1, 1)
	var x int = cursor.X
	var y int = cursor.Y
	foreground := termbox.ColorWhite
	background := termbox.ColorBlack

	log.Printf("Cursor Position: %d", hilite)
	for _, option := range m.Options {
		if y == hilite {
			foreground = termbox.ColorBlack
			background = termbox.ColorWhite
		} else {
			foreground = termbox.ColorWhite
			background = termbox.ColorBlack
		}
		for _, char := range option.Description {
			termbox.SetCell(x, y, rune(char), foreground, background)
			x++
		}
		y++
		x = cursor.X
	}
	termbox.Flush()
}

func (m *Menu) Execute(s *state.GameState) int {
	s.ClearTerminal()
	m.Draw(1)
	return m.Handle(s, m)
}

func (m *Menu) GetOptionHandle(index int) OptionHandler {
	for _, option := range m.Options {
		if option.ID == index {
			return option.Do
		}
	}
	return nil
}

func StartGameHandle(s *state.GameState, m *Menu, data ...interface{}) int {
	return 1
}

func ExitGameHandle(s *state.GameState, m *Menu, data ...interface{}) int {
	return 2
}

func StartMenuInputHandle(s *state.GameState, m *Menu) int {
	var cursorPos int = 1
	m.DrawBorder()
	s.MessageLine.Println("Welcome to Wizard!")
	for {
		select {
		case ev := <-s.Events:
			switch {
			case ev.Key == termbox.KeyEsc:
				return 2
			case ev.Key == termbox.KeyEnter:
				if handle := m.GetOptionHandle(cursorPos); handle != nil {
					return handle(s, m, nil)
				}
				break
			case ev.Ch == 'k':
				if cursorPos > 1 {
					cursorPos--
				}
				break
			case ev.Ch == 'j':
				if cursorPos < len(m.Options) {
					cursorPos++
				}
				break
			}
			m.Draw(cursorPos)
		}
	}
	return 2
}
