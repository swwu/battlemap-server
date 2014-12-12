package variable

import "github.com/swwu/battlemap-server/classes"

type reducerFn func(float64) float64
type conditionFn func() bool

type reducerVariable struct {
	id           string
	defaultValue float64
	value        float64
	reducerOps   []classes.ReducerOp
}

func (rv *reducerVariable) Id() string {
	return rv.id
}

func (rv *reducerVariable) Value() float64 {
	return rv.value
}

func (rv *reducerVariable) AddReducerOp(op classes.ReducerOp) {
	rv.reducerOps = append(rv.reducerOps, op)
}

// reducer variables reduce across operations
type basicReducerOp struct {
	effect      classes.Effect
	conditionFn conditionFn
	reducerFn   reducerFn
	precedence  int
}

func NewBasicReducerOp(effect classes.Effect,
	reducerFn reducerFn, conditionFn conditionFn, precedence int) classes.ReducerOp {
	return &basicReducerOp{
		effect:      effect,
		conditionFn: conditionFn,
		reducerFn:   reducerFn,
		precedence:  precedence,
	}
}

func (ro *basicReducerOp) Effect() classes.Effect {
	return ro.effect
}

func (ro *basicReducerOp) Condition() bool {
	return true
}

func (ro *basicReducerOp) Precedence() int {
	return ro.precedence
}
