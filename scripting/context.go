package scripting

import (
  "github.com/swwu/v8.go"
)

var engine *v8.Engine = nil

var contexts map[string]*v8.Context = map[string]*v8.Context{}

func GetEngine() *v8.Engine {
  if engine == nil {
    engine = v8.NewEngine()
  }

  return engine
}


func GetNullContextScope() v8.ContextScope {
  context := GetEngine().NewContext(nil)

  csChan := make(chan v8.ContextScope)
  go context.Scope(func(cs v8.ContextScope) {
    println("ASDASD")
    csChan <- cs
  })
  println("ASDASDASD")
  ret := <-csChan
  return ret
}

