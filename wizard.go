package main

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/menu"
	"github.com/tywkeene/wizard/state"
	"log"
)

func main() {
	defer termbox.Close()
	s := state.NewGameState()
	s.Init()
	startMenu := menu.NewMenu(0, 0,
		s.TerminalWidth, s.TerminalHeight, menu.StartMenuInputHandle)
	startMenu.AddOption("Wizard", menu.EmptyOptionHandle)
	startMenu.AddOption("Start Game", menu.StartMenuStartGame)
	startMenu.AddOption("Exit", menu.StartMenuExitGame)
	startMenu.Execute(s)
	s.ClearTerminal()
	s.MainLoop()
	log.Println("Main loop exited\n")
}
