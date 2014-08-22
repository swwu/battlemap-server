package entity

import (
  "github.com/swwu/v8.go"

  "github.com/swwu/battlemap-server/scripting"
  "github.com/swwu/battlemap-server/effect"
  "github.com/swwu/battlemap-server/logging"
)

type Footprint struct {

}

type Collider interface {
  Footprint() Footprint
}

// an entity is defined by its variables and its effect
type Entity interface {
  Variables() map[string]float64

  Reset()
  Calculate()
  Recalculate()

  AddEffect(eff effect.Effect)

  // returns a *v8.Value instead of *v8.Object (since object can't be easily
  // converted back to value)
  V8Accessor() *v8.ObjectTemplate
}

type entity struct {
  variables map[string]float64

  effects []effect.Effect
}


func NewEntity() (ent Entity) {
  return &entity {
    variables: map[string]float64{},
    effects: []effect.Effect{},
  }
}


func (ent *entity) Variables() map[string]float64 {
  return ent.variables
}

func (ent *entity) Reset() {
  ent.variables = map[string]float64{}
}

func (ent *entity) Calculate() {
  for _,eff := range ent.effects {
    eff.OnEffect(ent)
  }
}

func (ent *entity) Recalculate() {
  ent.Reset()
  ent.Calculate()
}


func (ent *entity) AddEffect(eff effect.Effect) {
  ent.effects = append(ent.effects, eff)
}


func (ent *entity) V8Accessor() *v8.ObjectTemplate {
  engine := scripting.GetEngine()

  objTemplate := engine.NewObjectTemplate()

  objTemplate.SetNamedPropertyHandler(
    // get
    func(name string, info v8.PropertyCallbackInfo) {
      info.ReturnValue().Set(engine.NewNumber(ent.variables[name]))
    },
    // set
    func(name string, value *v8.Value, info v8.PropertyCallbackInfo) {
      if value.IsNumber() {
        ent.variables[name] = scripting.NumberFromV8Value(value, ent.variables[name])
        info.ReturnValue().Set(value)
      } else {
        logging.Warning.Println(
          "Attempted to insert non-numerical value into entity variables")
      }
    },
    // query
    func(name string, info v8.PropertyCallbackInfo) {
    },
    // delete
    func(name string, info v8.PropertyCallbackInfo) {
    },
    // enumerate
    func(info v8.PropertyCallbackInfo) {
    },
    nil,
    )

  return objTemplate
}


