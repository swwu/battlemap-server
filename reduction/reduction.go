package reduction

import "github.com/swwu/battlemap-server/classes"

type evalFn func(map[string]float64, map[string]classes.ReducerVariable)

type reduction struct {
	id            string
	dependencyIds []string
	modifyIds     []string
	evalFn        evalFn
}

func NewReduction(id string, dependencyIds []string, modifyIds []string,
	evalFn evalFn) classes.Reduction {
	return &reduction{
		id:            id,
		dependencyIds: dependencyIds,
		modifyIds:     modifyIds,
		evalFn:        evalFn,
	}
}

func (r *reduction) Id() string {
	return r.id
}

func (r *reduction) DependencyIds() []string {
	return r.dependencyIds
}
func (r *reduction) ModifyIds() []string {
	return r.modifyIds
}

func (r *reduction) Dependencies(ent classes.Entity) []classes.Variable {
	ret := []classes.Variable{}
	for _, id := range r.DependencyIds() {
		ret = append(ret, ent.VariableContext().Variable(id))
	}
	return ret
}
func (r *reduction) Modifies(ent classes.Entity) []classes.ReducerVariable {
	ret := []classes.ReducerVariable{}
	for _, id := range r.ModifyIds() {
		ret = append(ret, ent.VariableContext().ReducerVariable(id))
	}
	return ret
}

func (r *reduction) DependencyValues(ent classes.Entity) map[string]float64 {
	depVars := r.Dependencies(ent)
	ret := map[string]float64{}
	for _, depVar := range depVars {
		ret[depVar.Id()] = depVar.Value()
	}
	return ret
}
func (r *reduction) ModifiedReducers(ent classes.Entity) map[string]classes.ReducerVariable {
	modVars := r.Modifies(ent)
	ret := map[string]classes.ReducerVariable{}
	for _, modVar := range modVars {
		ret[modVar.Id()] = modVar
	}
	return ret
}

func (r *reduction) Eval(ent classes.Entity) {
	r.evalFn(r.DependencyValues(ent), r.ModifiedReducers(ent))
}
