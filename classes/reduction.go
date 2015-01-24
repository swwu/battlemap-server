package classes

type Reduction interface {
	Id() string
	DependencyIds() []string
	ModifyIds() []string
	Dependencies(ent Entity) []Variable
	Modifies(ent Entity) []ReducerVariable
	DependencyValues(ent Entity) map[string]float64
	ModifiedReducers(ent Entity) map[string]ReducerVariable
	Eval(ent Entity)
}
