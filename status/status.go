package status

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/position"
)

type StatusLine struct {
	Position *position.Position
}

func MakeStatusLine(x int, y int, width int) *StatusLine {
	return &StatusLine{
		Position: position.NewPosition(x, y, x, y, width, 1),
	}
}

func (s *StatusLine) Printf(messageFmt string, args ...interface{}) {
	var message string
	fmt.Sprintf(message, messageFmt, args)
	for x, b := range message {
		termbox.SetCell(x, s.Position.Y, rune(b), termbox.ColorWhite, termbox.ColorBlack)
	}
}

func (s *StatusLine) Println(message string) {
	var x int = 0
	for _, b := range message {
		termbox.SetCell(x, s.Position.Y, rune(b), termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
}

func (s *StatusLine) Clear() {
	for x := 0; x < s.Position.Width; x++ {
		termbox.SetCell(x, s.Position.Y, rune(' '), termbox.ColorWhite, termbox.ColorBlack)
	}
}
