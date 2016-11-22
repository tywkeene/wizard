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

	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3

	MinMapWidth  = 2
	MinMapHeight = 2

	MinRoomWidth  = 4
	MaxRoomWidth  = 10
	MinRoomHeight = 4
	MaxRoomHeight = 10
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

func (l *Level) IsPosInsideLevel(p *position.Position) bool {
	if (p.X+p.Width) < l.Width-2 && (p.Y+p.Height) < l.Height-2 && p.X >= 2 && p.Y >= 2 {
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

func (l *Level) CoordinatesToRoom(x int, y int) *room.Room {
	p := &position.Position{-1, -1, x, y, 1, 1}
	return l.PosToRoom(p)
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

func (l Level) PlaceRoomDoor(r *room.Room) {
	var side *position.Position
	walls := []*position.Position{r.MiddleWallNorthPos(), r.MiddleWallEastPos(), r.MiddleWallSouthPos(), r.MiddleWallWestPos()}
	for direction, side := range walls {
		if l.DoesPosHaveWall(side, direction) == false {
			break
		}
		log.Printf("Direction %d has wall @ Position: [x:%d/y:%d]", direction, side.X, side.Y)
	}
	l.Map[side.X][side.Y] = &TileDoor
}

func GetRandomDirection() int {
	return dice.MakeDie(0, DirWest).Roll()
}

func (l *Level) GetRandomRoom(r *room.Room) *room.Room {
	return l.Rooms[dice.MakeDie(0, len(l.Rooms)).Roll()]
}

func (l *Level) PlaceRoomFloor(r *room.Room) {
	for x := r.Pos.X; x < (r.Pos.X + r.Pos.Width); x++ {
		for y := r.Pos.Y; y < (r.Pos.Y + r.Pos.Height); y++ {
			l.Map[x][y] = &TileFloor
		}
	}
}

func (l *Level) DoesPosHaveWall(p *position.Position, direction int) bool {
	wallTypes := []*Tile{&TileTopWall, &TileSideWall, &TileTopLeftCorner,
		&TileTopRightCorner, &TileBottomLeftCorner, &TileBottomRightCorner}
	switch direction {
	case DirNorth:
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y-1] == wall {
				return true
			}
		}
		break
	case DirEast:
		for _, wall := range wallTypes {
			if l.Map[p.X+1][p.Y] == wall {
				return true
			}
		}
		break
	case DirSouth:
		for _, wall := range wallTypes {
			if l.Map[p.X][p.Y+1] == wall {
				return true
			}
		}
		break
	case DirWest:
		for _, wall := range wallTypes {
			if l.Map[p.X-1][p.Y] == wall {
				return true
			}
		}
		break
	}
	return false
}

func (l *Level) GenerateRandomRoom() *room.Room {
	pos := &position.Position{X: dice.MakeDie(MinMapWidth, l.Width).RollEven(),
		Y:      dice.MakeDie(MinMapHeight, l.Height).RollEven(),
		Width:  dice.MakeDie(MinRoomWidth, MaxRoomWidth).Roll(),
		Height: dice.MakeDie(MinRoomHeight, MaxRoomHeight).Roll()}
	r := &room.Room{Pos: pos}
	wontFit := 0
	for {
		if l.IsInsideRoom(r.Pos) == false && l.IsPosInsideLevel(r.Pos) == true {
			break
		}
		r.Pos.X = dice.MakeDie(MinMapWidth, l.Width).Roll()
		r.Pos.Y = dice.MakeDie(MinMapHeight, l.Height).Roll()
		wontFit++
		if wontFit == 100 {
			return nil
		}
	}
	log.Printf("New room @ [x:%d/y:%d] Size:[%dx%d]", r.Pos.X, r.Pos.Y, pos.Width, pos.Height)
	l.PlaceRoomWalls(r)
	l.PlaceRoomFloor(r)
	return r
}

func MakeLevel(maxRooms int, width int, height int) *Level {
	l := &Level{Width: width, Height: height, Player: nil,
		Entities: nil, Map: nil, Rooms: nil}
	l.Map = MakeFloor(width, height, &TileNil)
	for i := 0; i < maxRooms; i++ {
		if r := l.GenerateRandomRoom(); r != nil {
			l.Rooms = append(l.Rooms, r)
			l.PlaceRoomWalls(r)
			l.PlaceRoomFloor(r)
		}
	}
	return l
}

func (l *Level) AddEntity(e entity.Entity) {
	l.Entities = append(l.Entities, e)
}

func (l *Level) CheckEntityCollision() {
	for _, e := range l.Entities {
		entityPos := e.GetPosition()
		if l.IsPosInsideLevel(entityPos) == false || l.Map[entityPos.X][entityPos.Y] != &TileFloor {
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
	l.CheckEntityCollision()
	l.DrawEntities()
}
