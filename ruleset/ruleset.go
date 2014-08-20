package ruleset

import (
  "github.com/swwu/battlemap-server/effect"
)

type RuleSet interface {
}

type ruleSet struct {
  effects map[string]effect.Effect

}
