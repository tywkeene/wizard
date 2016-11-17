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

type Game struct {
	Running      bool
	Events       chan termbox.Event
	CurrentLevel *level.Level
	MessageLine  *status.StatusLine
	PlayerStatus *status.StatusLine
	Player       *monster.Monster
}

var GameState *Game

func newGame() *Game {
	return &Game{true, make(chan termbox.Event), nil, nil, nil, nil}
}

func (g *Game) StartEventRoutine() {
	go func(g *Game) {
		for {
			g.Events <- termbox.PollEvent()
		}
	}(g)
}

func (g *Game) UpdateState() {
	g.CurrentLevel.DrawMap()
	termbox.Flush()
	g.MessageLine.Clear()
	g.PlayerStatus.Clear()
}

func (g *Game) mainLoop() {
	g.MessageLine.Println("Welcome to wizard!")
	for g.Running == true {
		g.UpdateState()
		select {
		case ev := <-g.Events:
			switch {
			case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
				g.Running = false
			case ev.Ch == 'k': //up
				g.Player.Move(g.Player.X, g.Player.Y-1)
				break
			case ev.Ch == 'j': //down
				g.Player.Move(g.Player.X, g.Player.Y+1)
				break
			case ev.Ch == 'h': //left
				g.Player.Move(g.Player.X-1, g.Player.Y)
				break
			case ev.Ch == 'l': //right
				g.Player.Move(g.Player.X+1, g.Player.Y)
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
	GameState = newGame()
	GameState.StartEventRoutine()
	log.Println("Intitializing RNG...")
	dice.Init()

	width, height := termbox.Size()
	log.Printf("Got terminal size: [%dx%d]...\n", width, height)

	log.Println("Initializing Message Line (Top)...")
	GameState.MessageLine = status.MakeStatusLine(0, 0, width)
	GameState.MessageLine.Clear()

	log.Println("Initializing Player Status (Bottom)...")
	GameState.PlayerStatus = status.MakeStatusLine(0, height-1, width)
	GameState.PlayerStatus.Clear()

	levelWidth := (width - 1)
	levelHeight := (height - 1)
	roomCount := 10
	GameState.CurrentLevel = level.MakeLevel(roomCount, levelWidth, levelHeight)

	stats := &monster.Status{Lvl: 1, Exp: 0, Hp: 25, Str: 12}
	GameState.Player = monster.MakeMonster(0, 0, "wizard", '@', stats)
}

func main() {
	defer termbox.Close()
	GameState.mainLoop()
	log.Println("Main loop exited\n")
}
