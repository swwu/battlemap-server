package variable

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/scripting"
)

type reducerFn func(float64, float64) float64

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
	// TODO: maybe move this to a separate function explicitly called in the
	// entity eval loop
	rv.value = rv.defaultValue

	for _, op := range rv.reducerOps {
		rv.value = op.Reduce(rv.value)
	}
	return rv.value
}

func (rv *reducerVariable) V8Accessor() *v8.ObjectTemplate {
	engine := scripting.GetEngine()

	objTemplate := engine.NewObjectTemplate()
	objTemplate.Bind("add", func(val *v8.Value) {
		rv.AddReducerOp(NewBasicReducerOp(func(l float64, r float64) float64 {
			return l + r
		}, 1, scripting.NumberFromV8Value(val, 0)))
	})
	objTemplate.Bind("max", func(val *v8.Value) {
		rv.AddReducerOp(NewBasicReducerOp(func(l float64, r float64) float64 {
			if l > r {
				return l
			} else {
				return r
			}
		}, 2, scripting.NumberFromV8Value(val, 0)))
	})

	return objTemplate
}

func (rv *reducerVariable) AddReducerOp(op classes.ReducerOp) {
	rv.reducerOps = append(rv.reducerOps, op)
}

// intended to be a commutative associative endomorphism on a binary operation (e.g. addition).
type basicReducerOp struct {
	reducerFn  reducerFn
	precedence int
	value      float64 // the value of the curried argument in the binary operation
}

func NewBasicReducerOp(reducerFn reducerFn,
	precedence int, value float64) classes.ReducerOp {
	return &basicReducerOp{
		reducerFn:  reducerFn,
		precedence: precedence,
		value:      value,
	}
}

func (ro *basicReducerOp) Condition() bool {
	return true
}

func (ro *basicReducerOp) Precedence() int {
	return ro.precedence
}

func (ro *basicReducerOp) Reduce(orig float64) float64 {
	if ro.reducerFn != nil {
		return ro.reducerFn(orig, ro.value)
	} else {
		return 0
	}

}
