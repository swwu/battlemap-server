package main

import (
	"fmt"
	//"io/ioutil"
	"os"

	"github.com/swwu/battlemap-server/entity"
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

	rules := rulesets["test"]

	ent := entity.NewEntity()

	ent.AddEffect(rules.Effects()["baseEntityRules"])
	ent.AddEffect(rules.Effects()["fighterClass"])

	ent.SetBaseValues(map[string]float64{
		"str_base":    14,
		"dex_base":    14,
		"con_base":    14,
		"int_base":    14,
		"wis_base":    14,
		"cha_base":    14,
		"fighter_lvl": 10,
	})

	ent.Recalculate()

	vc := ent.VariableContext()

	vals := []string{"str", "str_mod", "dex", "dex_mod", "fighter_lvl", "bab",
		"melee_ab", "will_save", "will_save_insight_bonus",
		"will_save_untyped_bonus", "fort_save", "ref_save", "ac_base",
		"ac_abmod_bonus", "ac", "ac_touch", "ac_flatfooted", "cmb", "cmd",
		"cmb_trip", "cmb_grapple", "hp"}
	for _, val := range vals {
		fmt.Println(val, "-", vc.Variable(val).Value())
	}

	gamespaces["testspace"].SetEntity("test", ent)

	server.Serve(gamespaces, rulesets)
}
