package state

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
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
	MessageLine    *status.StatusLine
	PlayerStatus   *status.StatusLine
	Player         *monster.Monster
	TerminalWidth  int
	TerminalHeight int
}

func (s *GameState) Init() {
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
	s.CurrentLevel = level.MakeLevel(roomCount, levelWidth, levelHeight)

	s.Player = monster.MakeMonster("wizard", '@')
	s.CurrentLevel.AddEntity(s.Player)
}

func NewGameState() *GameState {
	s := &GameState{true, make(chan termbox.Event), nil, nil, nil, nil, 0, 0}
	return s
}

func (g *GameState) StartEventRoutine() {
	go func(g *GameState) {
		for {
			g.Events <- termbox.PollEvent()
		}
	}(g)
}

func (g *GameState) UpdateState() {
	g.CurrentLevel.UpdateMap()
	termbox.Flush()
	g.MessageLine.Clear()
	g.PlayerStatus.Clear()
}

func (g *GameState) MainLoop() {
	g.MessageLine.Println("Welcome to wizard!")
	for g.Running == true {
		g.UpdateState()
		playerPos := g.Player.GetPosition()
		select {
		case ev := <-g.Events:
			switch {
			case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
				g.Running = false
			case ev.Ch == 'k': //up
				g.Player.Move(playerPos.X, playerPos.Y-1)
				break
			case ev.Ch == 'j': //down
				g.Player.Move(playerPos.X, playerPos.Y+1)
				break
			case ev.Ch == 'h': //left
				g.Player.Move(playerPos.X-1, playerPos.Y)
				break
			case ev.Ch == 'l': //right
				g.Player.Move(playerPos.X+1, playerPos.Y)
				break
			case ev.Key == termbox.KeyF12:
				break
			}
		}
	}
}
