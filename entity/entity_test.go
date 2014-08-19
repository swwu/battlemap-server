package entity

import (
  "testing"
  "github.com/swwu/battlemap-server/effects"
)

// effects are applied at all
func TestEffectAddition(t *testing.T) {
  ent := NewEntity()
  ent.AddEffect(effects.NewBaseCharEffect(&map[string]int{
      "STR": 10,
      "DEX": 10,
      "CON": 10,
      "WIS": 10,
      "INT": 10,
      "CHA": 10,
  }))
  ent.Recalculate()

  if (len(ent.Variables()) < 6) {
    t.Fail()
  }
}

// effects are applied in priority order
func TestEffectSorting(t *testing.T) {
  t.Fail()
}
