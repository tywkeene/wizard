package main

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/menu"
	"github.com/tywkeene/wizard/state"
	"log"
)

func main() {
	defer termbox.Close()
	gameState := state.NewGameState()
	gameState.Init()
	startMenu := menu.NewMenu(0, 0,
		gameState.TerminalWidth, gameState.TerminalHeight, menu.StartMenuInputHandle)
	startMenu.AddOption("Wizard", menu.EmptyOptionHandle)
	startMenu.AddOption("Start Game", menu.StartMenuStartGame)
	startMenu.AddOption("Exit", menu.StartMenuExitGame)
	startMenu.Execute(gameState)
	gameState.MainLoop()
	log.Println("Main loop exited\n")
}
