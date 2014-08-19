package effects


type Effect interface {
  DoEffect(variables *map[string]int)
  DisplayType() string
  Priority() int // applied in priority order (low to high)
  /* in general +- should be on priority 0, /* should be on priority 1 */
}


// base character effect
type baseCharEffect struct {
  attributes *map[string]int
}

func NewBaseCharEffect(attributes *map[string]int) *baseCharEffect {
  return &baseCharEffect{
    attributes: attributes,
  }
}

func (eff *baseCharEffect) DoEffect(variables *map[string]int) {
  for k, v := range *eff.attributes {
    (*variables)[k] += v
  }
}

func (eff *baseCharEffect) Priority() int {
  return 0
}

func (eff *baseCharEffect) DisplayType() string {
  return "base"
}





