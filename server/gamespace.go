package server

import (
	"github.com/swwu/battlemap-server/entity"
	"github.com/swwu/battlemap-server/ruleset"
)

type Gamespace interface {
	Entity(id string) entity.Entity
	Ruleset() ruleset.Ruleset
}

type gamespace struct {
	entities map[string]entity.Entity
	ruleset  ruleset.Ruleset
}

func NewGamespace(ruleset ruleset.Ruleset) Gamespace {
	return &gamespace{
		entities: map[string]entity.Entity{},
		ruleset:  ruleset,
	}
}

func (gs *gamespace) Entity(id string) entity.Entity {
	return gs.entities[id]
}

func (gs *gamespace) Ruleset() ruleset.Ruleset {
	return gs.ruleset
}