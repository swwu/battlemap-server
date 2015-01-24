package classes

import (
	"github.com/swwu/v8.go"
)

type V8AccessorProvider interface {
	V8Accessor() *v8.ObjectTemplate
}

/*
 * Effects mutate the state of entities and are persistent
 */
type Effect interface {
	Id() string
	DisplayName() string
	DisplayType() string

	RuleIds() []string
	Rules() []Rule
}
