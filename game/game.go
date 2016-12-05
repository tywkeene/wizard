package game

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/item"
	"github.com/tywkeene/wizard/state"
)

func MainLoop(s *state.GameState) {
	s.MessageLine.Println("Welcome to wizard!")
	var updateState bool
	player := s.Player

	for s.Running == true {
		s.UpdateState()
		updateState = false

		playerPos := player.GetPosition()
		tileEntities := s.CurrentLevel.GetEntitiesAtPosition(playerPos)
		if len(tileEntities) == 1 {
			s.MessageLine.Println("There is a " + tileEntities[0].GetName())
		}
		for updateState == false {
			select {
			case ev := <-s.Events:
				switch {
				case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
					return
				case ev.Ch == ',': //Pickup item at player's position
					if len(tileEntities) == 0 {
						s.MessageLine.Println("There's nothing here to pick up")
						updateState = false
						break
					} else {
						tileItem := tileEntities[0].(*item.Item)
						player.PickupItem(tileItem)
						s.CurrentLevel.Entities.Remove(tileItem)
						s.MessageLine.Println("You pick up the " + tileItem.GetName())
						updateState = true
					}
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
