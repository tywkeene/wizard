package item

import (
	"github.com/nsf/termbox-go"
	"github.com/tywkeene/wizard/dice"
	"github.com/tywkeene/wizard/entity"
	"github.com/tywkeene/wizard/position"
)

type ItemActionHandle func(e *entity.Entity)

type ItemInfo struct {
	Description string
	SpawnChance int
}

type Item struct {
	Info     *ItemInfo
	Apply    ItemActionHandle
	Name     string
	ID       int
	Position *position.Position
	Symbol   rune
	Passable bool
	Type     int
}

var ItemPlaceHolder = &Item{
	Position: &position.Position{-1, -1, -1, -1, 1, 1},
	Symbol:   'X',
	Name:     "Placeholder Item",
	Passable: true,
	Type:     entity.EntityTypeItem,
	Info: &ItemInfo{
		Description: "A placeholder item used by the creator to test the universe",
		SpawnChance: 1,
	},
	Apply: func(e *entity.Entity) {},
}

var ( //Scrolls
	ItemIdentifyScroll = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '?',
		Name:     "Scroll of identify",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "A useful piece of parchment describing a magical item",
			SpawnChance: 10,
		},
	}
)

var ( //All coins
	ItemSteelCoin = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '$',
		Name:     "Steel Coin",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "An old currency, some kingdoms do not accept it.",
			SpawnChance: 20,
		},
	}

	ItemCopperCoin = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '$',
		Name:     "Copper Coin",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "Still almost useless. Popular amongst peasants",
			SpawnChance: 20,
		},
	}

	ItemSilverCoin = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '$',
		Name:     "Silver Coin ",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "Base currency for most kingdoms.",
			SpawnChance: 10,
		},
	}

	ItemGoldCoin = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '$',
		Name:     "Gold Coin ",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "Most valuable currency. Only the richest possess it.",
			SpawnChance: 5,
		},
	}
)

var ( //Consumables
	ItemMarysHerb = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '~',
		Name:     "Mary's Herb",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "Reduces pain. Causes less loss of health points for a time",
			SpawnChance: 10,
		},
	}
	ItemMushroomEnlightenment = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '~',
		Name:     "Mushroom of Enlightenment",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "Causes great hallucinations. May broaden user's mind, granting additional magic points for a time",
			SpawnChance: 10,
		},
	}
)

var ( //All wands
	ItemTeleWand = &Item{
		Position: &position.Position{-1, -1, -1, -1, 1, 1},
		Symbol:   '/',
		Name:     "Wand of Teleportation",
		Passable: true,
		Type:     entity.EntityTypeItem,
		Apply:    func(e *entity.Entity) {},
		Info: &ItemInfo{
			Description: "A curious wand that bends space and time",
			SpawnChance: 5,
		},
	}
)

const (
	ItemTypeScroll = iota
	ItemTypeCoin
	ItemTypeConsumable
	ItemTypeWand
)

var (
	AllScrolls     = []*Item{ItemIdentifyScroll}
	AllCoins       = []*Item{ItemSteelCoin, ItemCopperCoin, ItemSilverCoin, ItemGoldCoin}
	AllConsumables = []*Item{ItemMarysHerb, ItemMushroomEnlightenment}
	AllWands       = []*Item{ItemTeleWand}
	AllItems       = [][]*Item{AllScrolls, AllCoins, AllConsumables, AllWands}
)

func GetRandomItem() *Item {
	randomItemType := dice.MakeDie(ItemTypeScroll, ItemTypeWand).Roll()
	randomItem := dice.MakeDie(0, len(AllItems[randomItemType])).Roll()
	item := AllItems[randomItemType][randomItem]
	return item
}

func (i *Item) Use(e *entity.Entity) {
	i.Apply(e)
}

func (i *Item) GetName() string {
	return i.Name
}

func (i *Item) GetID() int {
	return i.ID
}

func (i *Item) GetPosition() *position.Position {
	return i.Position
}

func (i *Item) GetSymbol() rune {
	return i.Symbol
}

func (i *Item) IsPassable() bool {
	return i.Passable
}

func (i *Item) GetType() int {
	return i.Type
}

func (m *Item) Move(x int, y int) {
	m.Position.PrevX = m.Position.X
	m.Position.PrevY = m.Position.Y
	m.Position.X = x
	m.Position.Y = y
}

func (m *Item) Draw() {
	p := m.GetPosition()
	termbox.SetCell(p.X, p.Y, m.GetSymbol(), termbox.ColorWhite, termbox.ColorBlack)
}
