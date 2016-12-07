package game

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/item"
	"github.com/tywkeene/wizard/menu"
	"github.com/tywkeene/wizard/state"
)

func MainLoop(s *state.GameState) {
	var updateState bool = false
	player := s.Player

	for s.Running == true {
		s.UpdateState()
		updateState = false

		playerPos := player.GetPosition()
		tileEntities := s.CurrentLevel.GetEntitiesAtPosition(playerPos)
		if len(tileEntities) == 1 {
			s.MessageLine.Println("There is a " + tileEntities[0].GetName() + " here")
		}
		for updateState == false {
			termbox.Flush()
			select {
			case ev := <-s.Events:
				switch {
				case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
					return
				case ev.Ch == ',': //Pickup item at player's position
					if len(tileEntities) == 0 {
						s.MessageLine.Println("There's nothing here to pick up")
					} else {
						tileItem := tileEntities[0].(*item.Item)
						player.PickupItem(tileItem)
						s.CurrentLevel.Entities.Remove(tileItem)
						s.MessageLine.Println("You pick up the " + tileItem.GetName())
					}
					break
				case ev.Ch == 'i': //Inventory
					invent := player.Items
					if len(invent.List) == 0 {
						s.MessageLine.Println("You have no items")
						break
					}
					inventMenu := menu.NewMenu(1, 1, s.TerminalWidth-2, s.TerminalHeight-2, menu.InventMenuInputHandle)
					for _, i := range invent.List {
						inventMenu.AddOption(i.GetName(), nil)
					}
					inventMenu.Execute(s)
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
