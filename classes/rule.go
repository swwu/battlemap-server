package classes

type Rule interface {
	Id() string
	DependencyIds() []string
	ModifyIds() []string
	Dependencies(ent Entity) []Variable
	Modifies(ent Entity) []Variable
	Eval(ent Entity)
}
