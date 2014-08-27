package ruleset

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/action"
	"github.com/swwu/battlemap-server/effect"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/scripting"
)

type Ruleset interface {
	Effects() map[string]effect.Effect

	ReadData(root string) error
}

type ruleset struct {
	effects map[string]effect.Effect

	v8context *v8.Context
}

func NewRuleset() Ruleset {
	ret := &ruleset{
		effects: map[string]effect.Effect{},
	}
	ret.constructGlobalContext()

	return ret
}

func NewRulesetFromData(path string) Ruleset {
	ret := NewRuleset()
	ret.ReadData(path)
	return ret
}

func (rs *ruleset) Effects() map[string]effect.Effect {
	return rs.effects
}

// read all js files from data/effects to make effects
func (rs *ruleset) ReadData(root string) error {
	effects := make([]effect.Effect, 0)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if path[len(path)-3:] == ".js" {
			logging.Info.Println("Loading:", path)
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			engine := scripting.GetEngine()

			compiledScript := engine.Compile(bytes, nil)
			context := rs.v8context

			cbChan := make(chan int)
			go context.Scope(func(cs v8.ContextScope) {
				cs.Run(compiledScript)
				cbChan <- 1
			})
			<-cbChan
		}
		return nil
	})

	for _, effect := range effects {
		rs.effects[effect.Id()] = effect
	}

	return err
}

/*
Construct the global scope's object template
*/
func (rs *ruleset) constructGlobalContext() {
	engine := scripting.GetEngine()

	global := engine.NewObjectTemplate()

	// define namespace is used to define effects etc
	defineTemplate := engine.NewObjectTemplate()
	defineTemplate.Bind("effect", func(obj *v8.Object) {
		newEff := effect.NewScriptEffect(
			scripting.StringFromV8Object(obj, "id", "defaultId"),
			scripting.StringFromV8Object(obj, "displayName", "unnamed"),
			scripting.StringFromV8Object(obj, "displayType", "none"),
			scripting.FnFromV8Object(obj, "onEffect", nil),
		)
		// TODO: check for id collision
		rs.effects[newEff.Id()] = newEff
	})
	defineTemplate.Bind("action", func(obj *v8.Object) {
		newAction := action.NewScriptAction(
		//scripting.StringFromV8Object(obj, "id", "defaultId"),
		//scripting.StringFromV8Object(obj, "displayName", "unnamed"),
		//scripting.StringFromV8Object(obj, "displayType", "none"),
		//scripting.FnFromV8Object(obj, "onEffect", nil),
		)
		// TODO: check for id collision
		//rs.effects[newEff.Id()] = newEff
		logging.Trace.Println(newAction)
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
	diceTemplate.Bind("effect", func(obj *v8.Object) {
		newEff := effect.NewScriptEffect(
			scripting.StringFromV8Object(obj, "id", "defaultId"),
			scripting.StringFromV8Object(obj, "displayName", "unnamed"),
			scripting.StringFromV8Object(obj, "displayType", "none"),
			scripting.FnFromV8Object(obj, "onEffect", nil),
		)
		// TODO: check for id collision
		rs.effects[newEff.Id()] = newEff
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

	rs.v8context = engine.NewContext(global)
}
