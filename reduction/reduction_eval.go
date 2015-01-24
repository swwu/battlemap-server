package reduction

import (
	"fmt"

	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/classes"
	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/scripting"
)

func MakeV8EvalFn(evalFn *v8.Function) evalFn {
	if evalFn != nil {
		return func(deps map[string]float64, mods map[string]classes.ReducerVariable) {
			fmt.Println(deps, mods)
			engine := scripting.GetEngine()

			context := engine.NewContext(nil)

			var retVal float64
			cbChan := make(chan int)
			go context.Scope(func(cs v8.ContextScope) {

				// dependency object is {<dependencyId>: <value>}
				depObj := engine.NewObject()
				for depKey, depValue := range deps {
					val := engine.NewNumber(depValue)
					depObj.ToObject().SetProperty(depKey, val, v8.PA_ReadOnly)
				}

				// modify object is {<modifyId>:func(accumulate_value)}
				modObjTemplate := engine.NewObjectTemplate()
				for modifyId, modVar := range mods {
					modObjTemplate.SetAccessor(modifyId,
						// get
						func(name string, info v8.AccessorCallbackInfo) {
							info.ReturnValue().Set(engine.NewInstanceOf(modVar.V8Accessor()))
						},
						// set
						func(name string, value *v8.Value, info v8.AccessorCallbackInfo) {
							logging.Warning.Println("Attempted to overwrite entity.vars")
						},
						nil,
						v8.PA_ReadOnly,
					)
				}
				modObj := engine.NewInstanceOf(modObjTemplate)

				retVal = scripting.NumberFromV8Value(evalFn.Call(depObj, modObj), 0)

				cbChan <- 1
			})
			<-cbChan
		}
	} else {
		return func(deps map[string]float64, mods map[string]classes.ReducerVariable) {}
	}
}
