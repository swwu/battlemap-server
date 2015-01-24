package classes

import (
	"github.com/swwu/v8.go"
)

// an entity is defined by its variables and its effect
type Entity interface {
	VariableContext() VariableContext

	BaseValues() map[string]float64

	SetBaseValues(vars map[string]float64)

	ReductionDependencyOrdering() ([]Reduction, error)

	Reset()
	Calculate()
	Recalculate()

	AddEffect(eff Effect)

	// returns a *v8.Value instead of *v8.Object (since object can't be easily
	// converted back to value)
	V8Accessor() *v8.ObjectTemplate

	JsonDump() ([]byte, error)
	JsonPut(jsonString []byte) error
}
