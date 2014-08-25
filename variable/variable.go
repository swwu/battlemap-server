package variable

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/scripting"
)

type Variable interface {
	Id() string
	Context() VariableContext
	DependencyIds() []string
	Value() float64
	OnEval()
}

type scriptVariable struct {
	id            string
	context       VariableContext
	dependencyIds []string
	onEvalFn      *v8.Function
	value         float64
}

func (sv *scriptVariable) Id() string {
	return sv.id
}

func (sv *scriptVariable) Context() VariableContext {
	return sv.context
}

func (sv *scriptVariable) DependencyIds() []string {
	return sv.dependencyIds
}

func (sv *scriptVariable) Value() float64 {
	return sv.value
}

func (sv *scriptVariable) OnEval() {
	if sv.onEvalFn == nil {
		// function is nil? nothing to do
		return
	}

	engine := scripting.GetEngine()

	context := engine.NewContext(nil)

	var retVal float64
	cbChan := make(chan int)
	go context.Scope(func(cs v8.ContextScope) {
		depObj := engine.NewObject()
		// pass the function the object {<dependencyId>: <value>}
		for _, depName := range sv.dependencyIds {
			depVar := sv.Context().Variable(depName)
			var val *v8.Value
			if depVar == nil {
				val = engine.NewNumber(0)
			} else {
				val = engine.NewNumber(depVar.Value())
			}
			depObj.ToObject().SetProperty(depName, val, v8.PA_ReadOnly)
		}
		retVal = scripting.NumberFromV8Value(sv.onEvalFn.Call(depObj), 0)

		cbChan <- 1
	})
	<-cbChan
	sv.value = retVal
}
