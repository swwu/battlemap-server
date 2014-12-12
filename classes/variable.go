package classes

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
}

type ReducerOp interface {
	Effect() Effect
	Condition() bool
	Precedence() int
	//Reduce(orig float64)
}
