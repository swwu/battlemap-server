package entity

import (
  "github.com/swwu/battlemap-server/effects"
  "sort"
)

type Footprint struct {

}

type Collider interface {
  Footprint() Footprint
}

// an entity is defined by its variables and its effects
type Entity interface {
  Variables() map[string]int

  Reset()
  Calculate()
  Recalculate()

  AddEffect(eff effects.Effect)
}

type entity struct {
  variables map[string]int

  effects []effects.Effect
}

// sorting boilerplate
type ByPriority []effects.Effect
func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
  return a[i].Priority() < a[j].Priority()
}


func NewEntity() (ent Entity) {
  return &entity {
    variables: make(map[string]int),
    effects: make([]effects.Effect, 0),
  }
}


func (ent *entity) Variables() map[string]int {
  return ent.variables
}

func (ent *entity) Reset() {
  ent.variables = make(map[string]int)
}

func (ent *entity) Calculate() {
  sort.Sort(ByPriority(ent.effects))
  for _,eff := range ent.effects {
    eff.DoEffect(&ent.variables)
  }
}

func (ent *entity) Recalculate() {
  ent.Reset()
  ent.Calculate()
}


func (ent *entity) AddEffect(eff effects.Effect) {
  ent.effects = append(ent.effects, eff)
}

