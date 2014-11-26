package variable

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/scripting"
)

type Variable interface {
	Id() string
	Context() VariableContext
	DependencyIds() []string
	ModifyIds() []string
	Value() float64
	OnEval()
}

type scriptVariable struct {
	id            string
	context       *variableContext
	dependencyIds []string
	modifyIds     []string
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

func (sv *scriptVariable) ModifyIds() []string {
	return sv.modifyIds
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

		// dependency object is {<dependencyId>: <value>}
		depObj := engine.NewObject()
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

		// modify object is {<modifyId>:func(accumulate_value)}
		modObjTemplate := engine.NewObjectTemplate()
		for _, modifyId := range sv.modifyIds {
			modVar := sv.context.accumVariables[modifyId]

			modObjTemplate.Bind(modifyId, func(val float64) {
				modVar.Accum(val)
			})
		}
		modObj := engine.NewInstanceOf(modObjTemplate)

		retVal = scripting.NumberFromV8Value(sv.onEvalFn.Call(depObj, modObj), 0)

		cbChan <- 1
	})
	<-cbChan
	sv.value = retVal
}

type DataVariable interface {
	Variable
	SetValue(val float64)
}

type dataVariable struct {
	id      string
	context *variableContext
	value   float64
}

func (dv *dataVariable) Id() string {
	return dv.id
}

func (dv *dataVariable) Context() VariableContext {
	return dv.context
}

func (dv *dataVariable) DependencyIds() []string {
	return []string{}
}

func (dv *dataVariable) ModifyIds() []string {
	return []string{}
}

func (dv *dataVariable) Value() float64 {
	return dv.value
}

func (dv *dataVariable) OnEval() {
}

func (dv *dataVariable) SetValue(val float64) {
	dv.value = val
}

type AccumVariable interface {
	Variable
	SetDependencyIds(deps []string)
	Accum(more float64)
}

type accumVariable struct {
	id            string
	context       VariableContext
	value         float64
	init          float64
	dependencyIds []string
	accumFn       func(float64, float64) float64
}

func (av *accumVariable) Id() string {
	return av.id
}

func (av *accumVariable) Context() VariableContext {
	return av.context
}

func (av *accumVariable) DependencyIds() []string {
	return av.dependencyIds
}

func (av *accumVariable) SetDependencyIds(ids []string) {
	av.dependencyIds = ids
}

func (av *accumVariable) ModifyIds() []string {
	return []string{} // accumulators don't modify things
}

func (av *accumVariable) Value() float64 {
	return av.value
}

func (av *accumVariable) OnEval() {
	av.value = av.accumFn(av.value, av.init)
}

func (av *accumVariable) Accum(more float64) {
	av.value = av.accumFn(av.value, more)
}

func addAccumFn(a float64, b float64) float64 {
	return a + b
}

func maxAccumFn(a float64, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}
