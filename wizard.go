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
	startMenu := menu.NewMenu(1, 1, s.TerminalWidth-2, s.TerminalHeight-2, menu.StartMenuInputHandle)
	startMenu.AddOption("Start", menu.StartGameHandle)
	startMenu.AddOption("Exit", menu.ExitGameHandle)
	switch startMenu.Execute(s) {
	case 1:
		game.MainLoop(s)
		break
	case 2:
		break
	}
	log.Println("Main loop exited\n")
}
