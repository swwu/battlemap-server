package main

import (
  "github.com/swwu/battlemap-server/entity"
  "github.com/swwu/battlemap-server/effects"
  "github.com/swwu/battlemap-server/dice"
  "fmt"
)

func main() {
  a :=  entity.NewEntity()

  a.AddEffect(effects.NewBaseCharEffect(&map[string]int{
      "STR": 10,
      "DEX": 10,
      "CON": 10,
      "WIS": 10,
      "INT": 10,
      "CHA": 10,
  }))
  a.Recalculate()

  //die := dice.NewDiceExpression(2,10)
  //fmt.Println(die.DisplayString())
  expr, err := dice.ParseDiceExpression("1d4+2d6+3 + \t   6")
  fmt.Println(err)
  fmt.Println(expr.DisplayString())
  fmt.Println(a)
}

