package rule

import (
	"github.com/swwu/battlemap-server/classes"
)

type ruleEvalFn func(ent classes.Entity)

type rule struct {
	id     string
	evalFn ruleEvalFn
}

func NewRule(id string, evalFn ruleEvalFn) classes.Rule {
	return &rule{
		id:     id,
		evalFn: evalFn,
	}
}

func (r *rule) Id() string {
	return r.id
}

func (r *rule) Eval(ent classes.Entity) {
	r.evalFn(ent)
}
