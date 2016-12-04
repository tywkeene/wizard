package status

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type StatusLine struct {
	X     int
	Y     int
	Width int
}

func MakeStatusLine(x int, y int, width int) *StatusLine {
	return &StatusLine{
		X:     x,
		Y:     y,
		Width: width,
	}
}

func (sl *StatusLine) Printf(messageFmt string, args ...interface{}) {
	var message string
	fmt.Sprintf(message, messageFmt, args)
	for x, b := range message {
		termbox.SetCell(x, sl.Y, rune(b), termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (sl *StatusLine) Println(message string) {
	var x int = 0
	for _, b := range message {
		termbox.SetCell(x, sl.Y, rune(b), termbox.ColorBlack, termbox.ColorWhite)
		x++
	}
	termbox.Flush()
}

func (sl *StatusLine) Clear() {
	for x := 0; x < sl.Width; x++ {
		termbox.SetCell(x, sl.Y, rune(' '), termbox.ColorBlack, termbox.ColorBlack)
	}
	termbox.Flush()
}
