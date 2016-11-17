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

	TileNil  = Tile{' ', false, nil}
	TileWall = Tile{'#', false, nil}

	DirNorth = 1
	DirEast  = 2
	DirSouth = 3
	DirWest  = 4
)

func NewPosition(x int, y int, width int, height int) *Position {
	return &Position{X: x, Y: y, Width: width, Height: height}
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

func (l *Level) EnsureInsideLevel(r *Room) bool {
	if ((r.Pos.X + r.Pos.Width) > l.Width-1) || ((r.Pos.Y + r.Pos.Height) > l.Height-1) {
		return false
	}
	return true
}

func (l *Level) RoomIsOverlapping(r *Room) bool {
	for x := r.Pos.X; x < r.Pos.Width; x++ {
		for y := r.Pos.Y; y < r.Pos.Height; y++ {
			if l.Map[x][y] != &TileNil {
				return true
			}
		}
	}
	return false
}

func (l *Level) RandomDoor(r *Room) {
	direction := dice.MakeDie(1, 4).Roll()
	switch direction {
	case 1: //North
		maxX := (r.Pos.X + r.Pos.Width)
		x := dice.MakeDie(r.Pos.X+1, maxX-1).Roll()
		l.Map[x][r.Pos.Y] = &TileDoor
		break
	case 2: //East
		maxY := (r.Pos.Y + r.Pos.Height)
		y := dice.MakeDie(r.Pos.Y+1, maxY-1).Roll()
		l.Map[r.Pos.X][y] = &TileDoor
		break
	case 3: //South
		maxX := (r.Pos.X + r.Pos.Width)
		x := dice.MakeDie(r.Pos.X+1, maxX-1).Roll()
		l.Map[x][r.Pos.Y+r.Pos.Height] = &TileDoor
		break
	case 4: //West
		maxY := (r.Pos.Y + r.Pos.Height)
		y := dice.MakeDie(r.Pos.Y+1, maxY-1).Roll()
		l.Map[r.Pos.X][y] = &TileDoor
		break
	}
}

func (l *Level) SetWalls(r *Room) {
	//North
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		l.Map[x][r.Pos.Y] = &TileWall
	}
	//East
	for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
		l.Map[r.Pos.X][y] = &TileWall
	}
	//South
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width + 1); x++ {
		l.Map[x][(r.Pos.Y + r.Pos.Height)] = &TileWall
	}
	//West
	for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
		l.Map[(r.Pos.X + r.Pos.Width)][y] = &TileWall
	}
}

func (l *Level) PosToRoom(x int, y int) *Room {
	for _, r := range l.Rooms {
		for i := r.Pos.X; i < (r.Pos.X + r.Pos.Width); i++ {
			for j := r.Pos.Y; j < (r.Pos.Y + r.Pos.Height); j++ {
				if x >= r.Pos.X && x < (r.Pos.X+r.Pos.Width) && y >= r.Pos.Y && y < (r.Pos.Y+r.Pos.Height) {
					return r
				}
			}
		}
	}
	return nil
}

func (l *Level) PosHasNeighbors(x int, y int) bool {
	if l.Map[x+1][y] != &TileNil || l.Map[x][y+1] != &TileNil ||
		l.Map[x+1][y+1] != &TileNil {
		return true
	}
	if l.Map[x-1][y] != &TileNil || l.Map[x][y-1] != &TileNil ||
		l.Map[x-1][y-1] != &TileNil {
		return true
	}
	return false
}

func (l *Level) DoesRoomFit(r *Room) bool {
	if l.EnsureInsideLevel(r) == false {
		return false
	}
	for i := r.Pos.X; i < (r.Pos.X + r.Pos.Width); i++ {
		for j := r.Pos.Y; j < (r.Pos.Y + r.Pos.Height); j++ {
			if l.Map[i][j] != &TileNil {
				return false
			}
		}
	}
	return true
}

func RandomDirection() int {
	return dice.MakeDie(DirNorth, DirWest).Roll()
}

func (l *Level) GetRandomRoom(r *Room) *Room {
	return l.Rooms[dice.MakeDie(0, len(l.Rooms)).Roll()]
}

func (l *Level) MakeRandomRoom() *Room {
	pos := NewPosition(dice.MakeDie(2, l.Width).Roll(),
		dice.MakeDie(2, l.Height).Roll(),
		dice.MakeDie(3, 10).Roll(),
		dice.MakeDie(3, 10).Roll())
	r := &Room{Pos: pos}
	for {
		if l.RoomIsOverlapping(r) == false && l.EnsureInsideLevel(r) == true && l.DoesRoomFit(r) == true {
			log.Printf("Found position for room [X:%d/Y:%d] %d/%d", r.Pos.X, r.Pos.Y, r.Pos.Width, r.Pos.Height)
			break
		} else {
			r.Pos.X = dice.MakeDie(2, l.Width-2).Roll()
			r.Pos.Y = dice.MakeDie(2, l.Height-2).Roll()
			r.Pos.Width = dice.MakeDie(4, 8).Roll()
			r.Pos.Height = dice.MakeDie(4, 8).Roll()
		}
	}
	for i := r.Pos.X; i < (r.Pos.X + r.Pos.Width); i++ {
		for j := r.Pos.Y; j < (r.Pos.Y + r.Pos.Height); j++ {
			l.Map[i][j] = &TileFloor
		}
	}
	l.SetWalls(r)
	l.RandomDoor(r)
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
	for i := 0; i < maxRooms; i++ {
		l.Rooms = append(l.Rooms, l.MakeRandomRoom())
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
