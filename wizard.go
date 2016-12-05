package main

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/game"
	"github.com/tywkeene/wizard/menu"
	"github.com/tywkeene/wizard/state"
	"log"
)

func main() {
	defer termbox.Close()
	s := state.NewGameState()
	s.Initialize()
	s.ClearTerminal()
	startMenu := menu.NewMenu(1, 1, s.TerminalWidth-2, s.TerminalHeight-2, menu.StartMenuInputHandle)
	startMenu.AddOption("Start", menu.StartMenuStartGame)
	startMenu.AddOption("Exit", menu.StartMenuExitGame)
	startMenu.Execute(s)
	game.MainLoop(s)
	log.Println("Main loop exited\n")
}
