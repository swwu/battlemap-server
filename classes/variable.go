package classes

import (
	"github.com/swwu/v8.go"
)

type VariableContext interface {
	Variable(id string) Variable
	DataVariable(id string) DataVariable
	ReducerVariable(id string) ReducerVariable

	Variables() map[string]Variable

	SetDataVariable(id string, defaultValue float64) (DataVariable, error)
	SetReducerVariable(id string, defaultValue float64) (ReducerVariable, error)
}

type Variable interface {
	Id() string
	Value() float64
}

type DataVariable interface {
	Variable
	SetValue(val float64)
}

type ReducerVariable interface {
	Variable
	AddReducerOp(op ReducerOp)
	V8Accessor() *v8.ObjectTemplate
}

type ReducerOp interface {
	Condition() bool
	Precedence() int
	Reduce(orig float64) float64
}
