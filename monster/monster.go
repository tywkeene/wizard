package monster

import (
	"github.com/tywkeene/wizard/dice"
	"log"
	"time"
)

type Status struct {
	Lvl int
	Exp int
	Hp  int
	Str int
}

type Monster struct {
	X           int
	Y           int
	PrevX       int
	PrevY       int
	CurrentRoom int
	Name        string
	Symbol      rune
	Stats       *Status
	//Inventory []*items.Item
}

var gnome = &Monster{0, 0, 0, 0, -1, "gnome", 'G', &Status{Lvl: 1, Exp: 5, Hp: 5, Str: 4}}
var goblin = &Monster{0, 0, 0, 0, -1, "goblin", 'g', &Status{Lvl: 2, Exp: 8, Hp: 7, Str: 3}}
var orc = &Monster{0, 0, 0, 0, -1, "orc", 'o', &Status{Lvl: 3, Exp: 12, Hp: 10, Str: 5}}
var ogre = &Monster{0, 0, 0, 0, -1, "ogre", 'Y', &Status{Lvl: 4, Exp: 15, Hp: 20, Str: 12}}

var AllMonsters = []*Monster{gnome, goblin, orc, ogre}

func NewRandomMonster() Monster {
	return *AllMonsters[dice.MakeDie(0, len(AllMonsters)).Roll()]
}

func RandomStats(lvl int, maxHp int, maxStr int) *Status {
	hp := dice.MakeDie(10, maxHp).Roll()
	str := dice.MakeDie(10, maxStr).Roll()
	log.Printf("Generated stats: Level: %d [HP:%d][Str:%d", lvl, hp, str)
	return &Status{lvl, 0, hp, str}
}

func logMonsterList(list []*Monster) {
	log.Println("=====================================")
	for _, m := range list {
		log.Printf("Initialized monster at position [%d:%d]: (%s) (%s)", m.X, m.Y, m.Name, string(m.Symbol))
		log.Printf("\tStats: Level: %d [HP:%d Str:%d]\n", m.Stats.Lvl, m.Stats.Hp, m.Stats.Str)
	}
	log.Println("=====================================")
}

func GrabRandomMonsterList(count int) []*Monster {
	log.Printf("Grabbing %d monsters...", count)
	list := make([]*Monster, 0)
	for i := 0; i < count; i++ {
		randomMonster := NewRandomMonster()
		list = append(list, &randomMonster)
		time.Sleep(1)
	}
	logMonsterList(list)
	return list
}

func MakeMonster(x int, y int, name string, symbol rune, stats *Status) *Monster {
	log.Printf("Initialized monster at position [%d:%d]: (%s) (%s)", x, y, name, string(symbol))
	return &Monster{X: x, Y: y,
		PrevX: x, PrevY: y,
		CurrentRoom: -1,
		Name:        name,
		Symbol:      symbol,
		Stats:       stats}
}

func (m *Monster) Move(newX int, newY int) {
	m.PrevX = m.X
	m.PrevY = m.Y
	m.X = newX
	m.Y = newY
}

func (m *Monster) SetCurrentRoom(room int) {
	m.CurrentRoom = room
}
