package ruleset

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/effect"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/rule"
	"github.com/swwu/battlemap-server/scripting"
)

type Ruleset interface {
	Effects() map[string]classes.Effect
	Rules() map[string]classes.Rule

	ReadData(root string) error
}

type ruleset struct {
	effects map[string]classes.Effect
	rules   map[string]classes.Rule
}

type intermediateEffect struct {
	id          string
	displayName string
	displayType string
	ruleIds     []string
}

func NewRuleset() Ruleset {
	ret := &ruleset{
		effects: map[string]classes.Effect{},
		rules:   map[string]classes.Rule{},
	}

	return ret
}

func NewRulesetFromData(path string) Ruleset {
	ret := NewRuleset()
	ret.ReadData(path)
	return ret
}

func (rs *ruleset) Effects() map[string]classes.Effect {
	return rs.effects
}

func (rs *ruleset) Rules() map[string]classes.Rule {
	return rs.rules
}

// read all js files from data/effects to make effects
func (rs *ruleset) ReadData(root string) error {

	intermediateEffects := map[string]intermediateEffect{}
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if path[len(path)-3:] == ".js" {
			logging.Info.Println("Loading:", path)
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			engine := scripting.GetEngine()

			compiledScript := engine.Compile(bytes, nil)
			context := rs.constructGlobalContext(intermediateEffects)

			cbChan := make(chan int)
			go context.Scope(func(cs v8.ContextScope) {
				cs.Run(compiledScript)
				cbChan <- 1
			})
			<-cbChan
		}
		return nil
	})

	for effId, intEffect := range intermediateEffects {
		rules := []classes.Rule{}
		for _, ruleId := range intEffect.ruleIds {
			realRule, ruleExists := rs.rules[ruleId]
			if ruleExists {
				rules = append(rules, realRule)
			}
		}
		rs.effects[effId] = effect.NewScriptEffect(
			intEffect.id,
			intEffect.displayName,
			intEffect.displayType,
			rules,
		)
	}

	return err
}

/*
Construct the global scope's object template
*/
func (rs *ruleset) constructGlobalContext(
	intermediateEffects map[string]intermediateEffect) *v8.Context {
	engine := scripting.GetEngine()

	global := engine.NewObjectTemplate()

	// define namespace is used to define effects etc
	defineTemplate := engine.NewObjectTemplate()
	defineTemplate.Bind("effect", func(obj *v8.Object) {
		newEff := intermediateEffect{
			id:          scripting.StringFromV8Object(obj, "id", "defaultId"),
			displayName: scripting.StringFromV8Object(obj, "displayName", "unnamed"),
			displayType: scripting.StringFromV8Object(obj, "displayType", "none"),
			ruleIds:     scripting.StringArrFromV8Object(obj, "rules", []string{}),
		}
		// TODO: check for id collision
		intermediateEffects[newEff.id] = newEff
	})

	defineTemplate.Bind("rule", func(obj *v8.Object) {
		newRule := rule.NewRule(
			scripting.StringFromV8Object(obj, "id", "defaultId"),
			rule.MakeV8RuleEvalFn(scripting.FnFromV8Object(obj, "eval", nil)),
		)
		rs.rules[newRule.Id()] = newRule
	})
	global.SetAccessor("define",
		// get
		func(name string, info v8.AccessorCallbackInfo) {
			info.ReturnValue().Set(engine.NewInstanceOf(defineTemplate))
		},
		// set - shouldn't ever be called because readonly
		func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
			logging.Warning.Println("Attempted to overwrite global.define")
		},
		nil,
		v8.PA_ReadOnly,
	)

	// dice namespace is used to evaluate dice expressions
	diceTemplate := engine.NewObjectTemplate()
	diceTemplate.Bind("roll", func(obj *v8.Object) {
	})
	global.SetAccessor("dice",
		// get
		func(name string, info v8.AccessorCallbackInfo) {
			info.ReturnValue().Set(engine.NewInstanceOf(diceTemplate))
		},
		// set - shouldn't ever be called because readonly
		func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
			logging.Warning.Println("Attempted to overwrite global.dice")
		},
		nil,
		v8.PA_ReadOnly,
	)

	return engine.NewContext(global)
}
