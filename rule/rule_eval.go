package rule

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/scripting"
)

func MakeV8RuleEvalFn(evalFn *v8.Function) ruleEvalFn {
	if evalFn != nil {
		return func(ent classes.Entity) {
			engine := scripting.GetEngine()

			context := engine.NewContext(nil)

			cbChan := make(chan int)
			go context.Scope(func(cs v8.ContextScope) {
				evalFn.Call(engine.NewInstanceOf(ent.V8Accessor()))
				cbChan <- 1
			})
			<-cbChan
		}
	} else {
		return func(ent classes.Entity) {}
	}
}
