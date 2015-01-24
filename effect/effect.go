package effect

import "github.com/swwu/battlemap-server/classes"

// javascript-code effect
type scriptEffect struct {
	id          string
	displayName string
	displayType string
	rules       []classes.Rule
}

func NewScriptEffect(id string, displayName string, displayType string,
	rules []classes.Rule) classes.Effect {
	return &scriptEffect{
		id:          id,
		displayName: displayName,
		displayType: displayType,
		rules:       rules,
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

func (eff *scriptEffect) RuleIds() []string {
	ret := make([]string, 0, len(eff.rules))
	for _, rule := range eff.rules {
		ret = append(ret, rule.Id())
	}
	return []string{}
}

func (eff *scriptEffect) Rules() []classes.Rule {
	return eff.rules
}
