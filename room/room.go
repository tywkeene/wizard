package room

import "github.com/tywkeene/wizard/position"

type Room struct {
	Pos     *position.Position
	Nearest *position.Position
}

func (r *Room) DirectionToWallPosition(direction int) *position.Position {
	walls := map[int]*position.Position{
		0: r.Pos.MiddleTopPosition(),
		1: r.Pos.MiddleRightPosition(),
		2: r.Pos.MiddleBottomPosition(),
		3: r.Pos.MiddleLeftPosition(),
	}
	return walls[direction]
}
