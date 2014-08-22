package main

import (
	"fmt"
	//"io/ioutil"
	"os"

	"github.com/swwu/battlemap-server/entity"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/ruleset"
)

func main() {

	logging.Init( /*ioutil.Discard*/ os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	rules := ruleset.NewRuleset()
	rules.ReadData("data/test_data")

	a := entity.NewEntity()

	a.AddEffect(rules.Effects()["baseStats"])
	a.AddEffect(rules.Effects()["statMods"])
	a.Recalculate()

	fmt.Println(a)

}
