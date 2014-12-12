package rule

import "github.com/swwu/battlemap-server/classes"

type evalFn func(map[string]float64, map[string]classes.ReducerVariable)

type rule struct {
	id            string
	dependencyIds []string
	modifyIds     []string
	evalFn        evalFn
}

func NewRule(id string, dependencyIds []string, modifyIds []string,
	evalFn evalFn) classes.Rule {
	return &rule{
		id:            id,
		dependencyIds: dependencyIds,
		modifyIds:     modifyIds,
		evalFn:        evalFn,
	}
}

func (r *rule) Id() string {
	return r.id
}

func (r *rule) DependencyIds() []string {
	return r.dependencyIds
}
func (r *rule) ModifyIds() []string {
	return r.modifyIds
}

func (r *rule) Dependencies(ent classes.Entity) []classes.Variable {
	ret := []classes.Variable{}
	for _, id := range r.DependencyIds() {
		ret = append(ret, ent.VariableContext().Variable(id))
	}
	return ret
}
func (r *rule) Modifies(ent classes.Entity) []classes.Variable {
	ret := []classes.Variable{}
	for _, id := range r.ModifyIds() {
		ret = append(ret, ent.VariableContext().Variable(id))
	}
	return ret
}

func (r *rule) Eval(ent classes.Entity) {
	r.evalFn(map[string]float64{}, map[string]classes.ReducerVariable{})
}
