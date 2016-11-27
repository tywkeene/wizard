package level

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
	"github.com/tywkeene/wizard/entity"
	"github.com/tywkeene/wizard/monster"
	"github.com/tywkeene/wizard/position"
	"github.com/tywkeene/wizard/room"
	"log"
)

type Path struct {
	From      *position.Position
	To        *position.Position
	Direction int
}

type Tile struct {
	Symbol   rune
	Passable bool
	Occupied entity.Entity
}

type Level struct {
	Width    int
	Height   int
	Player   *monster.Monster
	Entities []entity.Entity
	Map      [][]*Tile
	Rooms    []*room.Room
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
)

const (
	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3

	MinMapWidth  = 3
	MinMapHeight = 3

	MinRoomWidth  = 4
	MinRoomHeight = 4

	SmallRoomMax = 8

	MediumRoomMax = 12

	LargeRoomMax = 16
)

func RandomRoomSize() int {
	size := dice.MakeDie(1, 4).Roll()
	switch size {
	case 1:
		log.Println("Small room")
		return SmallRoomMax
		break
	case 2:
		log.Println("Medium room")
		return MediumRoomMax
		break
	case 3:
		log.Println("Large room")
		return LargeRoomMax
		break
	}
	return -1
}

func DirectionToString(direction int) string {
	var directionStrings = map[int]string{
		DirNorth: "North",
		DirEast:  "East",
		DirSouth: "South",
		DirWest:  "West",
	}
	return directionStrings[direction]
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

func (l *Level) IsPositionInsideLevel(p *position.Position) bool {
	if p.X < l.Width && p.Y < l.Height && p.X > 0 && p.Y > 0 {
		return true
	}
	return false
}

func (l *Level) IsRoomInsideLevel(r *room.Room) bool {
	if (r.Pos.X+r.Pos.Width) < l.Width-2 && (r.Pos.Y+r.Pos.Height) < l.Height-2 && r.Pos.X >= 0 && r.Pos.Y >= 0 {
		return true
	}
	return false
}

func (l *Level) IsInsideRoom(p *position.Position) bool {
	for _, r2 := range l.Rooms {
		if p.X < (r2.Pos.X+r2.Pos.Width+2) &&
			(p.X+p.Width+2) > r2.Pos.X &&
			p.Y < (r2.Pos.Y+r2.Pos.Height+2) &&
			(p.Y+p.Height+2) > r2.Pos.Y {
			return true
		}
	}
	return false
}

func (l *Level) PosToRoom(p *position.Position) *room.Room {
	for _, r := range l.Rooms {
		if l.IsInsideRoom(r.Pos) == true {
			return r
		}
	}
	return nil
}

func (l *Level) PlaceRoomWalls(r *room.Room) {
	topLeft := r.TopLeftCornerPos()
	l.Map[topLeft.X][topLeft.Y] = &TileTopLeftCorner

	topRight := r.TopRightCornerPos()
	l.Map[topRight.X][topRight.Y] = &TileTopRightCorner

	bottomLeft := r.BottomLeftCornerPos()
	l.Map[bottomLeft.X][bottomLeft.Y] = &TileBottomLeftCorner

	bottomRight := r.BottomRightCornerPos()
	l.Map[bottomRight.X][bottomRight.Y] = &TileBottomRightCorner
	//North
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		l.Map[x][r.Pos.Y-1] = &TileTopWall
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
		l.Map[r.Pos.X-1][y] = &TileSideWall
	}
}

func GetRandomDirection() int {
	return dice.MakeDie(DirNorth, DirWest).Roll()
}

func (l *Level) GetRandomRoom(r *room.Room) *room.Room {
	return l.Rooms[dice.MakeDie(0, len(l.Rooms)).Roll()]
}

func (l *Level) DoesPosHaveWall(p *position.Position, direction int) bool {
	wallTypes := []*Tile{&TileTopWall, &TileSideWall, &TileTopLeftCorner,
		&TileTopRightCorner, &TileBottomLeftCorner, &TileBottomRightCorner}
	switch direction {
	case DirNorth:
		p.Y--
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y] == wall ||
				l.IsPositionInsideLevel(p) == false {
				return true
			}
		}
		p.Y++
		break
	case DirEast:
		p.X++
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y] == wall ||
				l.IsPositionInsideLevel(p) == false {
				return true
			}
		}
		p.X--
		break
	case DirSouth:
		p.Y++
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y] == wall ||
				l.IsPositionInsideLevel(p) == false {
				return true
			}
		}
		p.Y--
		break
	case DirWest:
		p.X--
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y] == wall ||
				l.IsPositionInsideLevel(p) == false {
				return true
			}
		}
		p.X++
		break
	}
	return false
}

func (l *Level) PlaceRoomDoor(r *room.Room) {
	var direction int
	var p *position.Position
	for {
		direction = GetRandomDirection()
		p = r.DirectionToWallPosition(direction)
		if l.DoesPosHaveWall(p, direction) == false {
			break
		}
	}
	log.Printf("Placing door on %s wall @ [X:%d/Y:%d]", DirectionToString(direction), p.X, p.Y)
	l.Map[p.X][p.Y] = &TileDoor
}

func (l *Level) PlaceRoomFloor(r *room.Room) {
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
			l.Map[x][y] = &TileFloor
		}
	}
}

func (l *Level) GenerateRandomRoom(x int, y int) *room.Room {
	roomSize := RandomRoomSize()
	pos := &position.Position{X: x,
		Y:      y,
		Width:  dice.MakeDie(MinRoomWidth, roomSize).RollEven(),
		Height: dice.MakeDie(MinRoomHeight, roomSize).RollEven()}
	r := &room.Room{Pos: pos}
	wontFit := 0
	for {
		if l.IsInsideRoom(r.Pos) == false && l.IsRoomInsideLevel(r) == true {
			break
		}
		r.Pos.X = dice.MakeDie(MinMapWidth, l.Width).Roll()
		r.Pos.Y = dice.MakeDie(MinMapHeight, l.Height).Roll()
		wontFit++
		if wontFit == 1000 {
			return nil
		}
	}
	l.PlaceRoomFloor(r)
	l.PlaceRoomWalls(r)
	log.Printf("New room @ [X:%d/Y:%d] Size:[%dx%d]", r.Pos.X, r.Pos.Y, pos.Width, pos.Height)
	return r
}

func MakeLevel(maxRooms int, width int, height int) *Level {
	l := &Level{Width: width, Height: height, Player: nil,
		Entities: nil, Map: nil, Rooms: nil}
	l.Map = MakeFloor(width, height, &TileNil)
	for i := 0; i < maxRooms; i++ {
		r := l.GenerateRandomRoom(dice.MakeDie(2, l.Width).Roll(),
			dice.MakeDie(2, l.Height).Roll())
		if r != nil {
			l.Rooms = append(l.Rooms, r)
		}
	}
	for _, r := range l.Rooms {
		l.PlaceRoomDoor(r)
	}
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if l.Map[x][y] == &TileNil {
				l.Map[x][y] = &TileFloor
			}
		}
	}
	return l
}

func (l *Level) AddEntity(e entity.Entity) {
	l.Entities = append(l.Entities, e)
}

func (l *Level) CheckCollision(p *position.Position) bool {
	if l.IsPositionInsideLevel(p) == false {
		return true
	}
	if tile := l.Map[p.X][p.Y]; tile.Passable == false {
		return true
	}
	return false
}

func (l *Level) CheckEntityCollisions() {
	for _, e := range l.Entities {
		entityPos := e.GetPosition()
		if l.CheckCollision(entityPos) == true {
			entityPos.X = entityPos.PrevX
			entityPos.Y = entityPos.PrevY
		}
	}
}

func (l *Level) DrawEntities() {
	for _, e := range l.Entities {
		pos := e.GetPosition()
		symbol := e.GetSymbol()
		termbox.SetCell(pos.X, pos.Y, symbol, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func (l *Level) DrawMap() {
	for x := 1; x < l.Width; x++ {
		for y := 1; y < l.Height; y++ {
			tile := l.Map[x][y]
			termbox.SetCell(x, y, tile.Symbol, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func (l *Level) UpdateMap() {
	l.DrawMap()
	l.CheckEntityCollisions()
	l.DrawEntities()
}
