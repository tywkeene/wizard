package room

import "github.com/tywkeene/wizard/position"

type Room struct {
	Pos     *position.Position
	Nearest *position.Position
}

func (r *Room) DirectionToWallPosition(direction int) *position.Position {
	walls := map[int]*position.Position{
		0: r.MiddleWallNorthPos(),
		1: r.MiddleWallEastPos(),
		2: r.MiddleWallSouthPos(),
		3: r.MiddleWallWestPos(),
	}
	return walls[direction]
}

//Corner positions
func (r *Room) TopLeftCornerPos() *position.Position {
	return &position.Position{-1, -1, r.Pos.X - 1, r.Pos.Y - 1, 1, 1}
}

func (r *Room) TopRightCornerPos() *position.Position {
	return &position.Position{-1, -1, (r.Pos.X + r.Pos.Width), r.Pos.Y - 1, 1, 1}
}

func (r *Room) BottomLeftCornerPos() *position.Position {
	return &position.Position{-1, -1, r.Pos.X - 1, (r.Pos.Y + r.Pos.Height), 1, 1}
}

func (r *Room) BottomRightCornerPos() *position.Position {
	return &position.Position{-1, -1, (r.Pos.X + r.Pos.Width), (r.Pos.Y + r.Pos.Height), 1, 1}
}

//Wall positions
func (r *Room) MiddleWallNorthPos() *position.Position {
	return &position.Position{-1, -1, (r.Pos.X + (r.Pos.Width / 2)), r.Pos.Y - 1, 1, 1}
}

func (r *Room) MiddleWallEastPos() *position.Position {
	return &position.Position{-1, -1, (r.Pos.X + r.Pos.Width), (r.Pos.Y + (r.Pos.Height / 2)), 1, 1}
}

func (r *Room) MiddleWallSouthPos() *position.Position {
	return &position.Position{-1, -1, (r.Pos.X + (r.Pos.Width / 2)), (r.Pos.Y + r.Pos.Height), 1, 1}
}

func (r *Room) MiddleWallWestPos() *position.Position {
	return &position.Position{-1, -1, r.Pos.X - 1, r.Pos.Y + (r.Pos.Height / 2), 1, 1}
}

//Room Center position.Position
func (r *Room) MiddleOfRoom() *position.Position {
	return &position.Position{-1, -1, r.Pos.X + (r.Pos.Width / 2), r.Pos.Y + (r.Pos.Height / 2), 1, 1}
}
