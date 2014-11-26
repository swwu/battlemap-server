package entity

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/effect"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/scripting"
	"github.com/swwu/battlemap-server/variable"
)

type Footprint struct {
}

type Collider interface {
	Footprint() Footprint
}

// an entity is defined by its variables and its effect
type Entity interface {
	VariableContext() variable.VariableContext

	BaseValues() map[string]float64

	SetVars(vars map[string]float64)

	Reset()
	Calculate()
	Recalculate()

	AddEffect(eff effect.Effect)

	// returns a *v8.Value instead of *v8.Object (since object can't be easily
	// converted back to value)
	V8Accessor() *v8.ObjectTemplate

	JsonDump() ([]byte, error)
	JsonPut(jsonString []byte) error
}

type entityJson struct {
	Id         string             `json:"id"`
	Vars       map[string]float64 `json:"vars"`
	BaseValues map[string]float64 `json:"baseValues"`
	Effects    map[string]bool    `json:"effects"`
}

type entity struct {
	variableContext variable.VariableContext

	baseValues map[string]float64

	effects []effect.Effect
}

func NewEntity() Entity {
	return &entity{
		variableContext: variable.NewContext(),
		baseValues:      map[string]float64{},
		effects:         []effect.Effect{},
	}
}

func (ent *entity) VariableContext() variable.VariableContext {
	return ent.variableContext
}

func (ent *entity) BaseValues() map[string]float64 {
	return ent.baseValues
}

func (ent *entity) SetVars(vars map[string]float64) {
	ent.Reset()

	for id, value := range vars {
		if ent.VariableContext().DataVariableExists(id) {
			ent.baseValues[id] = value
		} else {
			logging.Warning.Println("Attempting to set non-data variable", id, ", skipping")
		}
	}

	ent.Calculate()
}

func (ent *entity) Reset() {
	ent.variableContext = variable.NewContext()

	// evaluate all effects to instantiate variables
	for _, eff := range ent.effects {
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
	dependencyOrder, err := ent.variableContext.DependencyOrdering()
	if err != nil {
		logging.Error.Println(err)
	}
	for _, variable := range dependencyOrder {
		variable.OnEval()
	}
}

func (ent *entity) Recalculate() {
	ent.Reset()
	ent.Calculate()
}

func (ent *entity) AddEffect(eff effect.Effect) {
	ent.effects = append(ent.effects, eff)
}

func (ent *entity) variableFromV8Object(obj *v8.Object) (variable.Variable, error) {
	return ent.variableContext.SetScriptVariable(
		scripting.StringFromV8Object(obj, "id", ""),
		scripting.StringArrFromV8Object(obj, "depends", []string{}),
		scripting.StringArrFromV8Object(obj, "modifies", []string{}),
		scripting.FnFromV8Object(obj, "onEval", nil),
	)
}

func (ent *entity) accumVariableFromV8Object(obj *v8.Object) (variable.Variable, error) {
	return ent.variableContext.SetAccumVariable(
		scripting.StringFromV8Object(obj, "id", ""),
		scripting.StringFromV8Object(obj, "op", "+"), // default operation is add
		scripting.NumberFromV8Object(obj, "init", 0), // default value is 0
	)
}

func (ent *entity) dataVariableFromV8Object(obj *v8.Object) (variable.Variable, error) {
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
	// a proxy is a variable whose purpose is to modify other variables
	varTemplate.Bind("newProxy", func(obj *v8.Object) {
		ent.variableFromV8Object(obj)
	})
	varTemplate.Bind("newAccum", func(obj *v8.Object) {
		ent.accumVariableFromV8Object(obj)
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
	ent.SetVars(jsonStruct.BaseValues)

	return nil
}
