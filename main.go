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
	fmt.Println("HELLO")

	rules := ruleset.NewRuleset()
	rules.ReadData("data/test_data")

	a := entity.NewEntity()

	a.AddEffect(rules.Effects()["baseStats"])
	a.AddEffect(rules.Effects()["baseEntityRules"])
	a.AddEffect(rules.Effects()["fighterClass"])
	a.Recalculate()

	vc := a.VariableContext()

	vals := []string{"str", "str_mod", "fighter_lvl", "bab", "will_save",
		"will_save_insight_bonus", "will_save_untyped_bonus", "fort_save",
		"ref_save", "hp"}
	for _, val := range vals {
		fmt.Println(val, "-", vc.Variable(val).Value())
	}
}
