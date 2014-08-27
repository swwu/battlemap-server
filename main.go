package main

import (
	"fmt"
	//"io/ioutil"
	"os"

	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/ruleset"
	"github.com/swwu/battlemap-server/server"
)

func main() {

	logging.Init( /*ioutil.Discard*/ os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	fmt.Println("HELLO")

	gamespaces := map[string]server.Gamespace{}

	rulesets := map[string]ruleset.Ruleset{}
	rulesets["test"] = ruleset.NewRulesetFromData("data/test_data")

	gamespaces["testspace"] = server.NewGamespace(rulesets["test"])

	/*
		rules := rulesets["test"]

		ent := entity.NewEntity()

		ent.AddEffect(rules.Effects()["baseStats"])
		ent.AddEffect(rules.Effects()["baseEntityRules"])
		ent.AddEffect(rules.Effects()["fighterClass"])
		ent.Recalculate()

		vc := ent.VariableContext()

		vals := []string{"str", "str_mod", "fighter_lvl", "bab", "will_save",
			"will_save_insight_bonus", "will_save_untyped_bonus", "fort_save",
			"ref_save", "hp"}
		for _, val := range vals {
			fmt.Println(val, "-", vc.Variable(val).Value())
		}
	*/

	server.Serve(gamespaces, rulesets)
}
