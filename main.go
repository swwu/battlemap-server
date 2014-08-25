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

	fmt.Println(vc)

	vals := []string{"fighter_lvl", "bab", "will_save", "fort_save", "ref_save"}
	fmt.Println(vc.Variable("str"))
	fmt.Println(vc.Variable("str_mod"))
	for _, val := range vals {
		fmt.Println(val, "-", vc.Variable(val).Value())
	}
}
