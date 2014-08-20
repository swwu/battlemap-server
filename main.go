package main

import (
  "fmt"
  //"io/ioutil"
  "os"

  "github.com/swwu/battlemap-server/entity"
  "github.com/swwu/battlemap-server/effect"
  "github.com/swwu/battlemap-server/logging"
)

func main() {

  logging.Init(/*ioutil.Discard*/os.Stdout, os.Stdout, os.Stdout, os.Stderr)

  effects, _ := effect.ReadEffects()
  fmt.Println(effects)

  a :=  entity.NewEntity()

  a.AddEffect(effects["baseDef"])
  a.AddEffect(effects["statMods"])
  a.Recalculate()

  fmt.Println(a)

}

