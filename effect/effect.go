package effect

import (
  "path/filepath"
  "os"
  "io/ioutil"

  "github.com/swwu/v8.go"

  "github.com/swwu/battlemap-server/scripting"
  "github.com/swwu/battlemap-server/logging"
)

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
  id string
  displayName string
  displayType string
  onEffectFn *v8.Function
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


/*
 * Functions for loading scriptEffects from ,js files
*/

func GenerateScriptEffect(defaultId string, script []byte, effs []Effect,
  cb func(effs []Effect) error) {
  engine := scripting.GetEngine()
  global := engine.NewObjectTemplate()

  global.Bind("addEffect", func(obj *v8.Object) {

    ret := &scriptEffect{
      id: defaultId,
      displayName: "unnamed property",
      displayType: "none",
      onEffectFn: nil,
    }

    ret.id          = scripting.StringFromV8Object(obj, "id", "defaultId")
    ret.displayName = scripting.StringFromV8Object(obj, "displayName", "unnamed")
    ret.displayType = scripting.StringFromV8Object(obj, "displayType", "none")
    ret.onEffectFn  = scripting.FnFromV8Object(obj, "onEffect", nil)

    effs = append(effs, ret)

  })

  compiledScript := engine.Compile(script, nil)
  context := engine.NewContext(global)

  context.Scope(func(cs v8.ContextScope) {
    cs.Run(compiledScript)
    cb(effs)
  })
}


// read all js files from data/effects to make effects
func ReadEffects() (map[string]Effect, error) {
  effects := make([]Effect, 0)

  root := "data/effects/"
  err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
    if path[len(path)-3:] == ".js" {
      logging.Info.Println("Loading:",path)
      bytes, err := ioutil.ReadFile(path)
      if err != nil {
        return err
      }

      cbChan := make(chan int)
      go GenerateScriptEffect(path[len(root):len(path)-3], bytes, effects,
      func(effs []Effect) error {
        effects = effs
        cbChan <- 1
        return nil
      })
      <-cbChan
    }
    return nil
  })

  ret := map[string]Effect{}
  for _, effect := range effects {
    ret[effect.Id()] = effect
  }

  return ret, err
}

