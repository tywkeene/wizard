package entity

import (
	"github.com/tywkeene/wizard/position"
	"log"
)

type Entity interface {
	GetName() string
	GetPosition() *position.Position
	GetSymbol() rune
	IsPassable() bool
	GetType() int
	GetID() int
	Move(int, int)
	Draw()
}

type EntityList struct {
	List map[int]Entity
}

const (
	EntityTypePlayer = iota
	EntityTypeMonster
	EntityTypeItem
)

func NewEntityList() *EntityList {
	return &EntityList{make(map[int]Entity)}
}

func (el *EntityList) GetAllOfType(entityType int) []Entity {
	list := make([]Entity, 0)
	for _, e := range el.List {
		if e.GetType() == entityType {
			list = append(list, e)
		}
	}
	return list
}

func (el *EntityList) Add(e Entity) {
	p := e.GetPosition()
	id := len(el.List)
	el.List[id] = e
	log.Printf("Added entity [ID:%d] '%s' (%c) @[X:%d/Y:%d]",
		e.GetID(), e.GetName(), e.GetSymbol(), p.X, p.Y)
}

func (el *EntityList) Remove(e Entity) {
	delete(el.List, e.GetID())
	p := e.GetPosition()
	log.Printf("Removed entity [ID:%d] '%s' (%c) @[X:%d/Y:%d]",
		e.GetID(), e.GetName(), e.GetSymbol(), p.X, p.Y)
}

func (el *EntityList) Count() int {
	return len(el.List)
}

func (el *EntityList) Get(ID int) Entity {
	return el.List[ID]
}

func TypeToString(typeValue int) string {
	typeStrings := map[int]string{
		EntityTypePlayer:  "player",
		EntityTypeMonster: "monster",
		EntityTypeItem:    "item",
	}
	return typeStrings[typeValue]
}
