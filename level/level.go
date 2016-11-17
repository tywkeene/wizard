package level

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
	"github.com/tywkeene/wizard/monster"
	"log"
)

type Tile struct {
	Symbol   rune
	Passable bool
	Occupied *monster.Monster
}

type Room struct {
	X      int
	Y      int
	Width  int
	Height int
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
	if ((r.X + r.Width) > l.Width-1) || ((r.Y + r.Height) > l.Height-1) {
		return false
	}
	return true
}

func (l *Level) RoomIsOverlapping(r *Room) bool {
	for x := r.X; x < r.Width; x++ {
		for y := r.Y; y < r.Height; y++ {
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
		maxX := (r.X + r.Width)
		x := dice.MakeDie(r.X, maxX-1).Roll()
		l.Map[x][r.Y] = &TileDoor
		break
	case 2: //East
		maxY := (r.Y + r.Height)
		y := dice.MakeDie(r.Y, maxY-1).Roll()
		l.Map[r.X][y] = &TileDoor
		break
	case 3: //South
		maxX := (r.X + r.Width)
		x := dice.MakeDie(r.X, maxX-1).Roll()
		l.Map[x][r.Y+r.Height] = &TileDoor
		break
	case 4: //West
		maxY := (r.Y + r.Height)
		y := dice.MakeDie(r.Y, maxY-1).Roll()
		l.Map[r.X][y] = &TileDoor
		break
	}
}

func (l *Level) SetWalls(r *Room) {
	//North
	for x := r.X; x < (r.X + r.Width); x++ {
		l.Map[x][r.Y] = &TileWall
	}
	//East
	for y := r.Y; y < (r.Y + r.Height); y++ {
		l.Map[r.X][y] = &TileWall
	}
	//South
	for x := r.X; x < (r.X + r.Width + 1); x++ {
		l.Map[x][(r.Y + r.Height)] = &TileWall
	}
	//West
	for y := r.Y; y < (r.Y + r.Height); y++ {
		l.Map[(r.X + r.Width)][y] = &TileWall
	}
}

func (l *Level) PosToRoom(x int, y int) *Room {
	for _, r := range l.Rooms {
		for i := r.X; i < (r.X + r.Width); i++ {
			for j := r.Y; j < (r.Y + r.Height); j++ {
				if x >= r.X && x < (r.X+r.Width) && y >= r.Y && y < (r.Y+r.Height) {
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
	if (r.X-1) < 2 || (r.X+1) > (l.Width-1) || (r.Y-1) < 1 || (r.Y+1) > (l.Height-1) {
		return false
	}
	for i := r.X; i < (r.X + r.Width); i++ {
		for j := r.Y; j < (r.Y + r.Height); j++ {
			if l.Map[i][j] != &TileNil {
				return false
			}
		}
	}
	//North
	for x := r.X; x < (r.X + r.Width); x++ {
		if l.Map[x][r.Y-1] != &TileNil {
			return false
		}
	}
	//East
	for y := r.Y; y < (r.Y + r.Height); y++ {
		if l.Map[(r.X+r.Width)+1][y] != &TileNil {
			return false
		}
	}
	//South
	for x := r.X; x < (r.X + r.Width + 1); x++ {
		if l.Map[x][(r.Y+r.Height)+1] != &TileNil {
			return false
		}
	}
	//West
	for y := r.Y; y < (r.Y + r.Height); y++ {
		if l.Map[r.X-1][y] != &TileNil {
			return false
		}
	}
	return true
}

func (l *Level) MakeRandomRoom() *Room {
	r := &Room{X: dice.MakeDie(2, l.Width-2).Roll(),
		Y:      dice.MakeDie(2, l.Height-2).Roll(),
		Width:  dice.MakeDie(3, 10).Roll(),
		Height: dice.MakeDie(3, 10).Roll()}
	var noFit int = 0
	for noFit < 10 {
		if l.RoomIsOverlapping(r) == false && l.EnsureInsideLevel(r) == true && l.DoesRoomFit(r) == true {
			log.Printf("Found position for room [X:%d/Y:%d] %d/%d", r.X, r.Y, r.Width, r.Height)
			break
		} else {
			noFit++
			r.X = dice.MakeDie(2, l.Width-2).Roll()
			r.Y = dice.MakeDie(2, l.Height-2).Roll()
			r.Width = dice.MakeDie(4, 8).Roll()
			r.Height = dice.MakeDie(4, 8).Roll()
		}
	}
	for i := r.X; i < (r.X + r.Width); i++ {
		for j := r.Y; j < (r.Y + r.Height); j++ {
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
	/*
		for i := 0; i < l.Width; i++ {
			for j := 0; j < l.Height; j++ {
				if l.Map[i][j] != &TileWall && l.Map[i][j] != &TileDoor {
					l.Map[i][j] = &TileFloor
				}
			}
		}
	*/
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
