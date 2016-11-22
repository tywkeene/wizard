package position

type Position struct {
	PrevX  int
	PrevY  int
	X      int
	Y      int
	Width  int
	Height int
}

func NewPosition(prevX int, prevY int, x int, y int, width int, height int) *Position {
	return &Position{prevX, prevY, x, y, width, height}
}
