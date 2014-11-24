package effect

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/scripting"
)

/*
 * Effect groups are groupings of effects that are all required to produce a
 * certain gameplay mechanic (e.g. a class + its class abilities)
 */
type EffectGroup interface {
	Id() string
	Effects() []Effect
}

/*
 * Effects mutate the state of entities and are persistent
 */
type Effect interface {
	Id() string
	DisplayName() string
	DisplayType() string
	/* in general +- should be on priority 0, /* should be on priority 1 */

	OnEffect(ent V8AccessorProvider)
}

type V8AccessorProvider interface {
	V8Accessor() *v8.ObjectTemplate
}

// javascript-code effect
type scriptEffect struct {
	id          string
	displayName string
	displayType string
	onEffectFn  *v8.Function
}

func NewScriptEffect(id string, displayName string, displayType string,
	onEffectFn *v8.Function) Effect {
	return &scriptEffect{
		id:          id,
		displayName: displayName,
		displayType: displayType,
		onEffectFn:  onEffectFn,
	}
}

func (eff *scriptEffect) Id() string {
	return eff.id
}

func (eff *scriptEffect) DisplayName() string {
	return eff.displayName
}

func (eff *scriptEffect) DisplayType() string {
	return eff.displayType
}

func (eff *scriptEffect) OnEffect(ent V8AccessorProvider) {
	if eff.onEffectFn == nil {
		// function is nil? nothing to do
		return
	}

	engine := scripting.GetEngine()
	objTemplate := ent.V8Accessor()

	context := engine.NewContext(nil)

	cbChan := make(chan int)
	go context.Scope(func(cs v8.ContextScope) {
		eff.onEffectFn.Call(engine.NewInstanceOf(objTemplate))
		cbChan <- 1
	})

	<-cbChan

}
