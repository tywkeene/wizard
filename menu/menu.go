package menu

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/position"
	"github.com/tywkeene/wizard/state"
	"log"
	"os"
)

type OptionHandler func(*state.GameState, *Menu) int
type InputHandler func(*state.GameState, *Menu)

type MenuOption struct {
	Key         rune
	Description string
	Do          OptionHandler
}

type Menu struct {
	Pos     *position.Position
	Options map[rune]*MenuOption
	Handle  InputHandler
}

const (
	OptionExit = iota
	OptionStartGame
)

var alpha = map[int]rune{
	0: 'a', 1: 'b', 2: 'c', 3: 'd',
	4: 'e', 5: 'f', 6: 'g', 7: 'h',
	8: 'i', 9: 'j', 10: 'k', 11: 'l',
	12: 'm', 13: 'n', 14: 'o', 15: 'p',
	16: 'q', 17: 'r', 18: 's', 19: 't',
	20: 'u', 21: 'v', 22: 'w', 23: 'x',
	24: 'y', 25: 'z',
}

func NewMenu(x int, y int, width int, height int, handle InputHandler) *Menu {
	pos := position.NewPosition(-1, -1, x, y, width, height)
	m := &Menu{pos, make(map[rune]*MenuOption), handle}
	return m
}

func NewMenuOption(key rune, description string, handle OptionHandler) *MenuOption {
	return &MenuOption{key, description, handle}
}

func (m *Menu) AddOption(description string, handle OptionHandler) {
	o := NewMenuOption(alpha[len(m.Options)], description, handle)
	m.Options[o.Key] = o
	log.Printf("Added option with key '%c' and description '%s'", o.Key, description)
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
	topLeft := m.Pos.TopLeftCornerPosition()
	termbox.SetCell(topLeft.X, topLeft.Y, TopLeftBorder, termbox.ColorWhite, termbox.ColorBlack)

	topRight := m.Pos.TopRightCornerPosition()
	termbox.SetCell(topRight.X, topRight.Y, TopRightBorder, termbox.ColorWhite, termbox.ColorBlack)

	bottomLeft := m.Pos.BottomLeftCornerPosition()
	termbox.SetCell(bottomLeft.X, bottomLeft.Y, BottomLeftBorder, termbox.ColorWhite, termbox.ColorBlack)

	bottomRight := m.Pos.BottomRightCornerPosition()
	termbox.SetCell(bottomRight.X, bottomRight.Y, BottomRightBorder, termbox.ColorWhite, termbox.ColorBlack)

	//Top
	for x := m.Pos.X; x < (m.Pos.X + m.Pos.Width); x++ {
		termbox.SetCell(x, m.Pos.Y-1, TopBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Right
	for y := m.Pos.Y; y < (m.Pos.Y + m.Pos.Height); y++ {
		termbox.SetCell(m.Pos.X+m.Pos.Width, y, SideBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Bottom
	for x := m.Pos.X; x < (m.Pos.X + m.Pos.Width); x++ {
		termbox.SetCell(x, m.Pos.Y+m.Pos.Height, TopBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
	//Left
	for y := m.Pos.Y; y < (m.Pos.Y + m.Pos.Height); y++ {
		termbox.SetCell(m.Pos.X-1, y, SideBorder, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func (m *Menu) Draw() {
	cursor := position.NewPosition(-1, -1, m.Pos.X, m.Pos.Y, 1, 1)
	var x int = cursor.X
	m.DrawBorder()
	for _, option := range m.Options {
		termbox.SetCell(x, cursor.Y, option.Key, termbox.ColorBlack, termbox.ColorWhite)
		x++
		termbox.SetCell(x, cursor.Y, ':', termbox.ColorBlack, termbox.ColorWhite)
		x++
		termbox.SetCell(x, cursor.Y, ' ', termbox.ColorBlack, termbox.ColorWhite)
		x++
		for _, b := range option.Description {
			termbox.SetCell(x, cursor.Y, b, termbox.ColorBlack, termbox.ColorWhite)
			x++
			termbox.Flush()
		}
		cursor.Y++
		x = cursor.X
	}
}

func (m *Menu) Execute(s *state.GameState) {
	m.Draw()
	m.Handle(s, m)
}

func EmptyOptionHandle(s *state.GameState, m *Menu) int {
	return 0
}

//Start Menu Handles
func StartMenuStartGame(s *state.GameState, m *Menu) int {
	return OptionStartGame
}

func StartMenuExitGame(s *state.GameState, m *Menu) int {
	return OptionExit
}

func getOptionIndex(opt rune) int {
	for i := 0; i < len(alpha); i++ {
		if alpha[i] == opt {
			return i
		}
	}
	return -1
}

func StartMenuInputHandle(s *state.GameState, m *Menu) {
	for {
		select {
		case ev := <-s.Events:
			optIndex := getOptionIndex(ev.Ch)
			log.Printf("Got option index %d", optIndex)
			if optIndex > len(m.Options) || optIndex == -1 {
				s.MessageLine.Println("Invalid menu option")
				continue
			} else {
				ret := m.Options[ev.Ch].Do(s, m)
				switch ret {
				case OptionStartGame:
					return
				case OptionExit:
					termbox.Close()
					os.Exit(0)
				}
			}
		}
	}
	m.Draw()
}
