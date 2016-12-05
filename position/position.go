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

//Corner positions
func (p *Position) TopLeftCornerPosition() *Position {
	return &Position{-1, -1, p.X - 1, p.Y - 1, 1, 1}
}

func (p *Position) TopRightCornerPosition() *Position {
	return &Position{-1, -1, (p.X + p.Width), p.Y - 1, 1, 1}
}

func (p *Position) BottomLeftCornerPosition() *Position {
	return &Position{-1, -1, p.X - 1, (p.Y + p.Height), 1, 1}
}

func (p *Position) BottomRightCornerPosition() *Position {
	return &Position{-1, -1, (p.X + p.Width), (p.Y + p.Height), 1, 1}
}

//Wall positions
func (p *Position) MiddleTopPosition() *Position {
	return &Position{-1, -1, (p.X + (p.Width / 2)), p.Y - 1, 1, 1}
}

func (p *Position) MiddleLeftPosition() *Position {
	return &Position{-1, -1, (p.X + p.Width), (p.Y + (p.Height / 2)), 1, 1}
}

func (p *Position) MiddleBottomPosition() *Position {
	return &Position{-1, -1, (p.X + (p.Width / 2)), (p.Y + p.Height), 1, 1}
}

func (p *Position) MiddleRightPosition() *Position {
	return &Position{-1, -1, p.X - 1, p.Y + (p.Height / 2), 1, 1}
}

//Position Center position.Position
func (p *Position) MiddleOfPosition() *Position {
	return &Position{-1, -1, p.X + (p.Width / 2), p.Y + (p.Height / 2), 1, 1}
}
