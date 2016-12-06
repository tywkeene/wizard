package state

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
	"github.com/tywkeene/wizard/entity"
	"github.com/tywkeene/wizard/level"
	"github.com/tywkeene/wizard/logging"
	"github.com/tywkeene/wizard/monster"
	"github.com/tywkeene/wizard/status"
	"log"
)

type GameState struct {
	Running        bool
	Events         chan termbox.Event
	CurrentLevel   *level.Level
	Player         *monster.Monster
	MessageLine    *status.StatusLine
	PlayerStatus   *status.StatusLine
	TerminalWidth  int
	TerminalHeight int
}

func (s *GameState) Initialize() {
	logging.OpenLog("./wizard.log")
	log.Println("Initializng termbox...")
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	log.Println("Initializng game state...")
	s.StartEventRoutine()
	log.Println("Intitializing RNG...")
	dice.Init()

	width, height := termbox.Size()
	log.Printf("Got terminal size: [%dx%d]...\n", width, height)

	//Top message line
	log.Println("Initializing Message Line (Top)...")
	s.MessageLine = status.MakeStatusLine(0, 0, width)
	s.MessageLine.Clear()

	//Bottom Player status line
	log.Println("Initializing Player Status (Bottom)...")
	s.PlayerStatus = status.MakeStatusLine(0, height-1, width)
	s.PlayerStatus.Clear()

	//Actual terminal dimensions
	s.TerminalWidth = width
	s.TerminalHeight = height

	//Level dimensions
	levelWidth := (width - 1)
	levelHeight := (height - 1)

	roomCount := 20
	itemCount := 10

	//Initialize current level
	s.CurrentLevel = level.MakeLevel(itemCount, roomCount, levelWidth, levelHeight)

	//Make our player monster
	s.Player = monster.MakeMonster("wizard", '@', entity.EntityTypePlayer)
	s.Player.Position = s.CurrentLevel.GetRandomPassableTile()
	s.CurrentLevel.Entities.Add(s.Player)

	s.CurrentLevel.ListRoomsInLog()
}

func NewGameState() *GameState {
	s := &GameState{true, make(chan termbox.Event), nil, nil, nil, nil, 0, 0}
	return s
}

func (s *GameState) StartEventRoutine() {
	go func(s *GameState) {
		for {
			s.Events <- termbox.PollEvent()
		}
	}(s)
}

func (s *GameState) ClearMessageLines() {
	s.MessageLine.Clear()
	s.PlayerStatus.Clear()
}

func (s *GameState) UpdateState() {
	s.CurrentLevel.UpdateMap()
	s.ClearMessageLines()
	termbox.Flush()
}

func (g GameState) ClearTerminal() {
	for x := 0; x < g.TerminalWidth; x++ {
		for y := 0; y < g.TerminalHeight; y++ {
			termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
		}
	}
}
