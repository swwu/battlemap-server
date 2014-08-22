package scripting

import (
	"github.com/swwu/v8.go"
)

var engine *v8.Engine = nil
var context *v8.Context = nil

func GetEngine() *v8.Engine {
	if engine == nil {
		engine = v8.NewEngine()
	}

	return engine
}
