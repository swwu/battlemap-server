package entity

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/scripting"
	"github.com/swwu/battlemap-server/variable"
)

type entityJson struct {
	Id         string             `json:"id"`
	Vars       map[string]float64 `json:"vars"`
	BaseValues map[string]float64 `json:"baseValues"`
	Effects    map[string]bool    `json:"effects"`
}

type entity struct {
	variableContext classes.VariableContext

	baseValues map[string]float64

	effects []classes.Effect
	rules   []classes.Rule
}

func NewEntity() classes.Entity {
	return &entity{
		variableContext: variable.NewContext(),
		baseValues:      map[string]float64{},
		effects:         []classes.Effect{},
	}
}

func (ent *entity) VariableContext() classes.VariableContext {
	return ent.variableContext
}

func (ent *entity) BaseValues() map[string]float64 {
	return ent.baseValues
}

func (ent *entity) SetBaseValues(vars map[string]float64) {
	ent.Reset()

	for id, value := range vars {
		if dataVar := ent.VariableContext().DataVariable(id); dataVar != nil {
			ent.baseValues[id] = value
		} else {
			logging.Warning.Println("Attempting to set non-data variable", id, ", skipping")
		}
	}

	ent.Calculate()
}

func (ent *entity) RuleDependencyOrdering() ([]classes.Rule, error) {
	queue := []classes.Rule{}
	sortedList := []classes.Rule{}

	// given a ruleid, number of dependencies left to evaluate for it before we
	// can safely enqueue it
	depsLeft := map[string]int{}

	// given a varname, what rules are modified by that variable
	varMods := map[string][]classes.Rule{}
	// given a varname, what rules does that variable depend on
	varDeps := map[string][]classes.Rule{}

	// first pass, initialize leaves into queue and varMods/varDeps
	for _, rule := range ent.rules {
		depVars := rule.Dependencies(ent)
		modVars := rule.Modifies(ent)
		if len(depVars) == 0 {
			queue = append(queue, rule)
		}

		for _, depVar := range depVars {
			varMods[depVar.Id()] = append(varMods[depVar.Id()], rule)
		}
		for _, modVar := range modVars {
			varDeps[modVar.Id()] = append(varDeps[modVar.Id()], rule)
		}
	}

	// util functions to easily get an array of deps/mods for each rule
	// these only work once varMods/varDeps are populated
	getMods := func(curRule classes.Rule) []classes.Rule {
		modVars := curRule.Modifies(ent)
		allMods := map[classes.Rule]bool{}

		for _, modVar := range modVars {
			for _, modRule := range varMods[modVar.Id()] {
				allMods[modRule] = true
			}
		}

		ret := make([]classes.Rule, 0, len(allMods))
		for k := range allMods {
			ret = append(ret, k)
		}
		return ret
	}
	getDeps := func(curRule classes.Rule) []classes.Rule {
		depVars := curRule.Dependencies(ent)
		allDeps := map[classes.Rule]bool{}

		for _, depVar := range depVars {
			for _, depRule := range varDeps[depVar.Id()] {
				allDeps[depRule] = true
			}
		}

		ret := make([]classes.Rule, 0, len(allDeps))
		for k := range allDeps {
			ret = append(ret, k)
		}
		return ret
	}

	// second pass, use varMods/varDeps to populate
	for _, rule := range ent.rules {
		deps := getDeps(rule)
		depsLeft[rule.Id()] = len(deps)
	}

	// go through the queue, enqueueing things
	for len(queue) > 0 {
		curRule := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		sortedList = append(sortedList, curRule)

		for _, modRule := range getMods(curRule) {
			// decrement the number untraversed incoming edges
			depsLeft[modRule.Id()] -= 1
			// if we're out of untraversed incoming edges then enqueue
			if depsLeft[modRule.Id()] == 0 {
				// I guess technically it's a stack because the ordering don't mattuh
				queue = append(queue, modRule)
			}
		}
	}

	// if any of our rules still have untraversed incoming edges then we have a
	// cycle (we weren't able to traverse all the rules)
	for _, depsLeft := range depsLeft {
		if depsLeft > 0 {
			return nil, fmt.Errorf("Cannot toposort - dependency graph has cycles")
		}
	}

	// otherwise yey
	return sortedList, nil
}

func (ent *entity) Reset() {
	ent.variableContext = variable.NewContext()

	// evaluate all effects to instantiate variables
	for _, eff := range ent.effects {
		//TODO: next up is effects add rules
		eff.OnEffect(ent)
	}
}

func (ent *entity) Calculate() {
	// apply base values
	for valueVar, baseValue := range ent.baseValues {
		if dataVar := ent.variableContext.DataVariable(valueVar); dataVar != nil {
			dataVar.SetValue(baseValue)
		} else {
			logging.Warning.Println("Entity base value was not a data variable")
		}
	}

	// evaluate all variable nodes in dependency order
	orderedRules, err := ent.RuleDependencyOrdering()
	if err != nil {
		logging.Error.Println(err)
	}
	for _, rule := range orderedRules {
		rule.Eval(ent)
	}
}

func (ent *entity) Recalculate() {
	ent.Reset()
	ent.Calculate()
}

func (ent *entity) AddEffect(eff classes.Effect) {
	ent.effects = append(ent.effects, eff)
}

func (ent *entity) variableFromV8Object(obj *v8.Object) (classes.Variable, error) {
	return ent.variableContext.SetReducerVariable(
		scripting.StringFromV8Object(obj, "id", ""),
		scripting.NumberFromV8Object(obj, "init", 0), // default value is 0
	)
}

func (ent *entity) dataVariableFromV8Object(obj *v8.Object) (classes.Variable, error) {
	return ent.variableContext.SetDataVariable(
		scripting.StringFromV8Object(obj, "id", ""),
		scripting.NumberFromV8Object(obj, "init", 0), // default value is 0
	)
}

func (ent *entity) V8VariableAccessor() *v8.ObjectTemplate {
	return nil
}

func (ent *entity) V8Accessor() *v8.ObjectTemplate {
	engine := scripting.GetEngine()

	varTemplate := engine.NewObjectTemplate()
	varTemplate.Bind("new", func(obj *v8.Object) {
		ent.variableFromV8Object(obj)
	})

	// a data variable is essentially a literal whose value can mutate and is
	// not supplied in an effect
	varTemplate.Bind("newData", func(obj *v8.Object) {
		ent.dataVariableFromV8Object(obj)
	})

	labelTemplate := engine.NewObjectTemplate()
	labelTemplate.Bind("new", func(fn *v8.Function) {
	})

	objTemplate := engine.NewObjectTemplate()
	objTemplate.SetAccessor("vars",
		// get
		func(name string, info v8.AccessorCallbackInfo) {
			info.ReturnValue().Set(engine.NewInstanceOf(varTemplate))
		},
		// set
		func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
			logging.Warning.Println("Attempted to overwrite entity.vars")
		},
		nil,
		v8.PA_ReadOnly,
	)
	objTemplate.SetAccessor("labels",
		// get
		func(name string, info v8.AccessorCallbackInfo) {
			info.ReturnValue().Set(engine.NewInstanceOf(labelTemplate))
		},
		// set
		func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
			logging.Warning.Println("Attempted to overwrite entity.labels")
		},
		nil,
		v8.PA_ReadOnly,
	)

	return objTemplate
}

func (ent *entity) JsonDump() ([]byte, error) {
	jsonStruct := &entityJson{
		Id:         "test",
		Vars:       map[string]float64{},
		BaseValues: ent.baseValues,
		Effects:    map[string]bool{},
	}
	vars := ent.VariableContext().Variables()
	for k, v := range vars {
		if !math.IsNaN(v.Value()) {
			jsonStruct.Vars[k] = v.Value()
		}
	}
	jsonString, err := json.Marshal(jsonStruct)
	if err != nil {
		logging.Warning.Println("Error marshaling entity")
		return nil, fmt.Errorf("Error marshaling entity")
	}
	return jsonString, nil
}

func (ent *entity) JsonPut(jsonString []byte) error {
	jsonStruct := &entityJson{}
	err := json.Unmarshal(jsonString, jsonStruct)
	if err != nil {
		logging.Warning.Println("Error unmarshaling entity")
		fmt.Println(err)
		fmt.Println(string(jsonString))
		return fmt.Errorf("Error unmarshaling entity")
	}
	ent.SetBaseValues(jsonStruct.BaseValues)

	return nil
}
