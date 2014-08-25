package entity

import (
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

	Reset()
	Calculate()
	Recalculate()

	AddEffect(eff effect.Effect)

	// returns a *v8.Value instead of *v8.Object (since object can't be easily
	// converted back to value)
	V8Accessor() *v8.ObjectTemplate
}

type entity struct {
	variableContext variable.VariableContext

	effects []effect.Effect
}

func NewEntity() (ent Entity) {
	return &entity{
		variableContext: variable.NewContext(),
		effects:         []effect.Effect{},
	}
}

func (ent *entity) VariableContext() variable.VariableContext {
	return ent.variableContext
}

func (ent *entity) Reset() {
	ent.variableContext = variable.NewContext()
}

func (ent *entity) Calculate() {
	for _, eff := range ent.effects {
		eff.OnEffect(ent)
	}

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
	id := scripting.StringFromV8Object(obj, "id", "")

	dependencies := scripting.StringArrFromV8Object(obj, "dependencies", []string{})

	return ent.variableContext.SetVariable(
		id,
		dependencies,
		scripting.FnFromV8Object(obj, "onEval", nil),
	)
}

func (ent *entity) V8Accessor() *v8.ObjectTemplate {
	engine := scripting.GetEngine()

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

	return objTemplate
}
