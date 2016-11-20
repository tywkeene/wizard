package level

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
	"github.com/tywkeene/wizard/monster"
	"log"
)

type Position struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Tile struct {
	Symbol   rune
	Passable bool
	Occupied *monster.Monster
}

type Room struct {
	Pos *Position
}

type Level struct {
	Width    int
	Height   int
	Player   *monster.Monster
	Monsters []*monster.Monster
	Map      [][]*Tile
	Rooms    []*Room
}

var (
	TileFloor      = Tile{'.', true, nil}
	TileDoor       = Tile{'+', true, nil}
	TileStairsUp   = Tile{'<', true, nil}
	TileStairsDown = Tile{'>', true, nil}

	TileNil = Tile{' ', false, nil}

	TileSideWall          = Tile{'│', false, nil}
	TileTopWall           = Tile{'─', false, nil}
	TileTopLeftCorner     = Tile{'┌', false, nil}
	TileTopRightCorner    = Tile{'┐', false, nil}
	TileBottomLeftCorner  = Tile{'└', false, nil}
	TileBottomRightCorner = Tile{'┘', false, nil}

	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3
)

//Corner positions
func (r *Room) TopLeftCornerPos() *Position {
	return &Position{r.Pos.X, r.Pos.Y, 1, 1}
}

func (r *Room) TopRightCornerPos() *Position {
	return &Position{(r.Pos.X + r.Pos.Width), r.Pos.Y, 1, 1}
}

func (r *Room) BottomLeftCornerPos() *Position {
	return &Position{r.Pos.X, (r.Pos.Y + r.Pos.Height), 1, 1}
}

func (r *Room) BottomRightCornerPos() *Position {
	return &Position{(r.Pos.X + r.Pos.Width), (r.Pos.Y + r.Pos.Height), 1, 1}
}

//Wall positions
func (r *Room) MiddleWallNorthPos() *Position {
	return &Position{(r.Pos.X + (r.Pos.Width / 2)), r.Pos.Y, 1, 1}
}

func (r *Room) MiddleWallEastPos() *Position {
	return &Position{(r.Pos.X + r.Pos.Width), (r.Pos.Y + (r.Pos.Y / 2)), 1, 1}
}

func (r *Room) MiddleWallSouthPos() *Position {
	return &Position{(r.Pos.X + (r.Pos.Width / 2)), (r.Pos.Y + r.Pos.Height), 1, 1}
}

func (r *Room) MiddleWallWestPos() *Position {
	return &Position{r.Pos.X, r.Pos.Y + (r.Pos.Height / 2), 1, 1}
}

func MakeFloor(width int, height int, defaultTile *Tile) [][]*Tile {
	floor := make([][]*Tile, width)
	for i := range floor {
		floor[i] = make([]*Tile, height)
	}
	for x, _ := range floor {
		for y, _ := range floor[x] {
			floor[x][y] = defaultTile
		}
	}
	return floor
}

func (l *Level) IsPosInsideLevel(p *Position) bool {
	if (p.X+p.Width) > l.Width-2 || (p.Y+p.Height) > l.Height-2 || p.X <= 2 || p.Y <= 2 {
		return false
	}
	return true
}

func (l *Level) IsInsideRoom(r *Room) bool {
	for _, r2 := range l.Rooms {
		if r.Pos.X < (r2.Pos.X+r2.Pos.Width+2) &&
			(r.Pos.X+r.Pos.Width+2) > r2.Pos.X &&
			r.Pos.Y < (r2.Pos.Y+r2.Pos.Height+2) &&
			(r.Pos.Y+r.Pos.Height+2) > r2.Pos.Y {
			return true
		}
	}
	return false
}

func (l *Level) PosToRoom(p *Position) *Room {
	for _, r := range l.Rooms {
		if l.IsInsideRoom(r) == true {
			return r
		}
	}
	return nil
}

func (l *Level) CoordinatesToRoom(x int, y int) *Room {
	p := &Position{x, y, 1, 1}
	return l.PosToRoom(p)
}

func (l *Level) PlaceRoomWalls(r *Room) {
	//North
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		l.Map[x][r.Pos.Y] = &TileTopWall
	}
	//East
	for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
		l.Map[(r.Pos.X + r.Pos.Width)][y] = &TileSideWall
	}
	//South
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		l.Map[x][(r.Pos.Y + r.Pos.Height)] = &TileTopWall
	}
	//West
	for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
		l.Map[r.Pos.X][y] = &TileSideWall
	}

	topLeft := r.TopLeftCornerPos()
	l.Map[topLeft.X][topLeft.Y] = &TileTopLeftCorner

	topRight := r.TopRightCornerPos()
	l.Map[topRight.X][topRight.Y] = &TileTopRightCorner

	bottomLeft := r.BottomLeftCornerPos()
	l.Map[bottomLeft.X][bottomLeft.Y] = &TileBottomLeftCorner

	bottomRight := r.BottomRightCornerPos()
	l.Map[bottomRight.X][bottomRight.Y] = &TileBottomRightCorner
}

func (l Level) PlaceRoomDoor(r *Room) {
	walls := []*Position{r.MiddleWallNorthPos(), r.MiddleWallEastPos(), r.MiddleWallSouthPos(), r.MiddleWallWestPos()}
	side := walls[RandomDirection()]
	l.Map[side.X][side.Y] = &TileDoor
}

func RandomDirection() int {
	return dice.MakeDie(0, DirWest).Roll()
}

func (l *Level) GetRandomRoom(r *Room) *Room {
	return l.Rooms[dice.MakeDie(0, len(l.Rooms)).Roll()]
}

func (l *Level) PlaceRoomFloor(r *Room) {
	for i := r.Pos.X; i < (r.Pos.X + r.Pos.Width); i++ {
		for j := r.Pos.Y; j < (r.Pos.Y + r.Pos.Height); j++ {
			l.Map[i][j] = &TileFloor
		}
	}
}
func (l *Level) GenerateRandomRoom() *Room {
	pos := &Position{X: dice.MakeDie(2, l.Width).Roll(),
		Y:      dice.MakeDie(2, l.Height).Roll(),
		Width:  dice.MakeDie(4, 10).RollEven(),
		Height: dice.MakeDie(4, 10).RollEven()}
	r := &Room{Pos: pos}
	for {
		if l.IsPosInsideLevel(r.Pos) == true && l.IsInsideRoom(r) == false {
			log.Printf("Found position for room [X:%d/Y:%d] %d/%d", r.Pos.X, r.Pos.Y, r.Pos.Width, r.Pos.Height)
			break
		} else {
			r.Pos.X = dice.MakeDie(2, l.Width).Roll()
			r.Pos.Y = dice.MakeDie(2, l.Height).Roll()
			r.Pos.Width = dice.MakeDie(6, 10).Roll()
			r.Pos.Height = dice.MakeDie(6, 10).Roll()
		}
	}
	l.PlaceRoomFloor(r)
	l.PlaceRoomWalls(r)
	l.PlaceRoomDoor(r)
	return r
}

func MakeLevel(maxRooms int, width int, height int) *Level {
	l := &Level{Width: width,
		Height:   height,
		Player:   nil,
		Monsters: nil,
		Map:      nil,
		Rooms:    nil}
	l.Map = MakeFloor(width, height, &TileNil)
	var wontFit int = 0
	for i := 0; i < maxRooms; i++ {
		r := l.GenerateRandomRoom()
		if r != nil {
			l.Rooms = append(l.Rooms, r)
		} else {
			wontFit++
			if wontFit == 10 {
				break
			}
		}
	}
	return l
}

func (l *Level) DrawMap() {
	for x := 1; x < l.Width; x++ {
		for y := 1; y < l.Height; y++ {
			tile := l.Map[x][y]
			termbox.SetCell(x, y, tile.Symbol, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}
