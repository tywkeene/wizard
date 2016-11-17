package status

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type StatusLine struct {
	x     int
	y     int
	width int
}

func MakeStatusLine(x int, y int, width int) *StatusLine {
	return &StatusLine{x, y, width}
}

func (sl *StatusLine) Printf(messageFmt string, args ...interface{}) {
	var message string
	fmt.Sprintf(message, messageFmt, args)
	for x, b := range message {
		termbox.SetCell(x, sl.y, rune(b), termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (sl *StatusLine) Println(message string) {
	for x, b := range message {
		termbox.SetCell(x, sl.y, rune(b), termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (sl *StatusLine) Clear() {
	for x := 0; x < sl.width; x++ {
		termbox.SetCell(x, sl.y, rune(' '), termbox.ColorBlack, termbox.ColorBlack)
	}
}
