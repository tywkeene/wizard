package main

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
	Running      bool
	Events       chan termbox.Event
	CurrentLevel *level.Level
	MessageLine  *status.StatusLine
	PlayerStatus *status.StatusLine
	Player       *monster.Monster
}

var (
	State       *GameState
	LevelWidth  int
	LevelHeight int
)

func NewGameState() *GameState {
	return &GameState{true, make(chan termbox.Event), nil, nil, nil, nil}
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

func (g *GameState) mainLoop() {
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
			case ev.Key == termbox.KeyF12: //Regen map
				g.CurrentLevel = level.MakeLevel(20, LevelWidth, LevelHeight)
				g.CurrentLevel.AddEntity(State.Player)
				break
			}
		}
	}
}

func init() {
	logging.OpenLog("./wizard.log")
	log.Println("Initializng termbox...")
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	log.Println("Initializng game state...")
	State = NewGameState()
	State.StartEventRoutine()
	log.Println("Intitializing RNG...")
	dice.Init()

	width, height := termbox.Size()
	log.Printf("Got terminal size: [%dx%d]...\n", width, height)

	log.Println("Initializing Message Line (Top)...")
	State.MessageLine = status.MakeStatusLine(0, 0, width)
	State.MessageLine.Clear()

	log.Println("Initializing Player Status (Bottom)...")
	State.PlayerStatus = status.MakeStatusLine(0, height-1, width)
	State.PlayerStatus.Clear()

	LevelWidth = (width - 1)
	LevelHeight = (height - 1)
	roomCount := 20
	State.CurrentLevel = level.MakeLevel(roomCount, LevelWidth, LevelHeight)

	State.Player = monster.MakeMonster(1, 1, "wizard", '@')
	State.CurrentLevel.AddEntity(State.Player)
}

func main() {
	defer termbox.Close()
	State.mainLoop()
	log.Println("Main loop exited\n")
}
