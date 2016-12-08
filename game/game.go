package game

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/state"
)

func MainLoop(s *state.GameState) {
	var updateState bool = false
	player := s.Player

	for s.Running == true {
		s.UpdateState()
		updateState = false

		playerPos := player.GetPosition()
		for updateState == false {
			termbox.Flush()
			select {
			case ev := <-s.Events:
				switch {
				case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
					return

				case ev.Ch == 'b': //south west
					player.Move(playerPos.X-1, playerPos.Y+1)
					updateState = true
					break
				case ev.Ch == 'n': //south east
					player.Move(playerPos.X+1, playerPos.Y+1)
					updateState = true
					break

				case ev.Ch == 'y': //north west
					player.Move(playerPos.X-1, playerPos.Y-1)
					updateState = true
					break
				case ev.Ch == 'u': //north east
					player.Move(playerPos.X+1, playerPos.Y-1)
					updateState = true
					break

				case ev.Ch == 'k': //up
					player.Move(playerPos.X, playerPos.Y-1)
					updateState = true
					break
				case ev.Ch == 'j': //down
					player.Move(playerPos.X, playerPos.Y+1)
					updateState = true
					break
				case ev.Ch == 'h': //left
					player.Move(playerPos.X-1, playerPos.Y)
					updateState = true
					break
				case ev.Ch == 'l': //right
					player.Move(playerPos.X+1, playerPos.Y)
					updateState = true
					break
				}
			}
		}
	}
}
