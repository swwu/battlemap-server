package classes

type Rule interface {
	Id() string

	// evaluating a rule places its reductions and variables into the entity's
	// context
	Eval(ent Entity)
}
