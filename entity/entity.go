package entity

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/reduction"
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

	effects    map[string]classes.Effect
	rules      map[string]classes.Rule
	reductions map[string]classes.Reduction
}

func NewEntity() classes.Entity {
	return &entity{
		variableContext: variable.NewContext(),
		baseValues:      map[string]float64{},
		effects:         map[string]classes.Effect{},
		rules:           map[string]classes.Rule{},
		reductions:      map[string]classes.Reduction{},
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
		ent.baseValues[id] = value
	}

	ent.Calculate()
}

func (ent *entity) ReductionDependencyOrdering() ([]classes.Reduction, error) {
	queue := []classes.Reduction{}
	sortedList := []classes.Reduction{}

	// given a reductionid, number of dependencies left to evaluate for it before we
	// can safely enqueue it
	depsLeft := map[string]int{}

	// given a varname, what reductions are modified by that variable
	varMods := map[string][]classes.Reduction{}
	// given a varname, what reductions does that variable depend on
	varDeps := map[string][]classes.Reduction{}

	// first pass, initialize leaves into queue and varMods/varDeps
	for _, reduction := range ent.reductions {
		depVars := reduction.Dependencies(ent)
		modVars := reduction.Modifies(ent)

		for _, depVar := range depVars {
			varMods[depVar.Id()] = append(varMods[depVar.Id()], reduction)
		}
		for _, modVar := range modVars {
			varDeps[modVar.Id()] = append(varDeps[modVar.Id()], reduction)
		}
	}

	// util functions to easily get an array of dep/mod reductions for each
	// reduction
	// these only work once varMods/varDeps are populated
	getMods := func(curReduction classes.Reduction) []classes.Reduction {
		modVars := curReduction.Modifies(ent)
		allMods := map[classes.Reduction]bool{}

		for _, modVar := range modVars {
			for _, modReduction := range varMods[modVar.Id()] {
				allMods[modReduction] = true
			}
		}

		ret := make([]classes.Reduction, 0, len(allMods))
		for k := range allMods {
			ret = append(ret, k)
		}
		return ret
	}
	getDeps := func(curReduction classes.Reduction) []classes.Reduction {
		depVars := curReduction.Dependencies(ent)
		allDeps := map[classes.Reduction]bool{}

		for _, depVar := range depVars {
			for _, depReduction := range varDeps[depVar.Id()] {
				allDeps[depReduction] = true
			}
		}

		ret := make([]classes.Reduction, 0, len(allDeps))
		for k := range allDeps {
			ret = append(ret, k)
		}
		return ret
	}

	// queue up all the leaf reductions (those without dependenies)
	for _, reduction := range ent.reductions {
		if len(getDeps(reduction)) == 0 {
			queue = append(queue, reduction)
		}
	}

	// second pass, use varMods/varDeps to populate
	for _, reduction := range ent.reductions {
		deps := getDeps(reduction)
		depsLeft[reduction.Id()] = len(deps)
	}

	// go through the queue, enqueueing things
	for len(queue) > 0 {
		curReduction := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		sortedList = append(sortedList, curReduction)

		for _, modReduction := range getMods(curReduction) {
			// decrement the number untraversed incoming edges
			depsLeft[modReduction.Id()] -= 1
			// if we're out of untraversed incoming edges then enqueue
			if depsLeft[modReduction.Id()] == 0 {
				// I guess technically it's a stack because the ordering don't mattuh
				queue = append(queue, modReduction)
			}
		}
	}

	// if any of our reductions still have untraversed incoming edges then we have a
	// cycle (we weren't able to traverse all the reductions)
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

	// reset all rules
	ent.rules = map[string]classes.Rule{}

	// evaluate all effects to get rules
	for _, eff := range ent.effects {
		for _, rule := range eff.Rules() {
			ent.rules[rule.Id()] = rule
		}
	}

	// evaluate all rules to get reductions and variables
	for _, rule := range ent.rules {
		rule.Eval(ent)
	}
}

func (ent *entity) Calculate() {
	// apply base values
	for valueVar, baseValue := range ent.baseValues {
		if dataVar := ent.variableContext.Variable(valueVar); dataVar == nil {
			ent.variableContext.SetDataVariable(valueVar, baseValue)
		} else {
			logging.Warning.Println("Attempting to set non-data variable", valueVar, " with basevalues, skipping")
		}
	}

	// evaluate all variable nodes in dependency order
	orderedReductions, err := ent.ReductionDependencyOrdering()
	if err != nil {
		logging.Error.Println(err)
	}
	logging.Warning.Println(orderedReductions)
	for _, reduction := range orderedReductions {
		reduction.Eval(ent)
	}
}

func (ent *entity) Recalculate() {
	ent.Reset()
	ent.Calculate()
}

func (ent *entity) AddEffect(eff classes.Effect) {
	if eff != nil {
		ent.effects[eff.Id()] = eff
	} else {
		logging.Warning.Println("Tried to add nil effect to entity")
	}
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

func (ent *entity) reductionFromV8Object(obj *v8.Object) (classes.Reduction, error) {
	id := scripting.StringFromV8Object(obj, "id", "")
	newRed := reduction.NewReduction(
		id,
		scripting.StringArrFromV8Object(obj, "depends", []string{}),
		scripting.StringArrFromV8Object(obj, "modifies", []string{}),
		reduction.MakeV8EvalFn(scripting.FnFromV8Object(obj, "eval", nil)),
	)
	ent.reductions[id] = newRed
	return newRed, nil
}

func (ent *entity) V8VariableAccessor() *v8.ObjectTemplate {
	return nil
}

func (ent *entity) V8Accessor() *v8.ObjectTemplate {
	engine := scripting.GetEngine()

	reductionTemplate := engine.NewObjectTemplate()
	reductionTemplate.Bind("new", func(obj *v8.Object) {
		ent.reductionFromV8Object(obj)
	})

	varTemplate := engine.NewObjectTemplate()
	varTemplate.Bind("new", func(obj *v8.Object) {
		ent.variableFromV8Object(obj)
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
	objTemplate.SetAccessor("reductions",
		// get
		func(name string, info v8.AccessorCallbackInfo) {
			info.ReturnValue().Set(engine.NewInstanceOf(reductionTemplate))
		},
		// set
		func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
			logging.Warning.Println("Attempted to overwrite entity.reductions")
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
