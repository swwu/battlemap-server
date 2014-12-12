package server

import (
	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/ruleset"
)

type Gamespace interface {
	Entity(id string) classes.Entity
	SetEntity(id string, entity classes.Entity)
	Ruleset() ruleset.Ruleset
}

type gamespace struct {
	entities map[string]classes.Entity
	ruleset  ruleset.Ruleset
}

func NewGamespace(ruleset ruleset.Ruleset) Gamespace {
	return &gamespace{
		entities: map[string]classes.Entity{},
		ruleset:  ruleset,
	}
}

func (gs *gamespace) Entity(id string) classes.Entity {
	return gs.entities[id]
}

func (gs *gamespace) SetEntity(id string, entity classes.Entity) {
	gs.entities[id] = entity
}

func (gs *gamespace) Ruleset() ruleset.Ruleset {
	return gs.ruleset
}
