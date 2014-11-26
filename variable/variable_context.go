package variable

import (
	"fmt"

	"github.com/swwu/v8.go"
)

// A general scope for game-related variables. Implements dependency patterns
// to guarantee consistency of variable evaluations.
type VariableContext interface {
	Variable(id string) Variable
	DataVariable(id string) DataVariable
	AccumVariable(id string) AccumVariable

	SetDataVariable(id string, value float64) (DataVariable, error)
	SetScriptVariable(id string, dependencies []string,
		modifies []string, onEvalFn *v8.Function) (Variable, error)
	SetAccumVariable(id string, op string, init float64) (AccumVariable, error)

	DataVariableExists(id string) bool

	Variables() map[string]Variable

	DependencyOrdering() ([]Variable, error)

	//Eval()
}

type variableContext struct {
	// all variables ever go here, even if they also go in other places
	variables map[string]Variable

	// data variables need to be tracked separately for mutations of underlying
	// state
	dataVariables map[string]DataVariable

	// accums need to be tracked separately to track modification correctness
	accumVariables map[string]AccumVariable

	// literal values that are accumulated onto accumulators before any other
	// eval is done
}

func NewContext() VariableContext {
	return &variableContext{
		variables:      map[string]Variable{},
		dataVariables:  map[string]DataVariable{},
		accumVariables: map[string]AccumVariable{},
	}
}

func (vc *variableContext) SetDataVariable(id string, value float64) (
	DataVariable, error) {
	if _, exists := vc.variables[id]; exists {
		return nil, fmt.Errorf("Variable with id", id, "already exists")
	}

	newVar := &dataVariable{
		id:      id,
		context: vc,
		value:   value,
	}
	vc.variables[id] = newVar
	vc.dataVariables[id] = newVar
	return newVar, nil
}

// Creates a new script variable. Script variables are immutable but can be
// changed by modifying accumulators or literal variables on which they depend
func (vc *variableContext) SetScriptVariable(id string, dependencyIds []string,
	modifyIds []string, onEvalFn *v8.Function) (Variable, error) {
	if _, exists := vc.variables[id]; exists {
		return nil, fmt.Errorf("Variable with id", id, "already exists")
	}

	newVar := &scriptVariable{
		id:            id,
		context:       vc,
		dependencyIds: dependencyIds,
		modifyIds:     modifyIds,
		onEvalFn:      onEvalFn,
	}
	vc.variables[id] = newVar
	return newVar, nil
}

// Creates a new accumulator. Accumulators are mutable using a given
// commutative operation
func (vc *variableContext) SetAccumVariable(id string,
	op string, init float64) (AccumVariable, error) {
	if _, exists := vc.variables[id]; exists {
		return nil, fmt.Errorf("Variable with id", id, "already exists")
	}

	newVar := &accumVariable{
		id:      id,
		context: vc,
		init:    init,
	}
	if op == "+" {
		newVar.accumFn = addAccumFn
	} else if op == "max" {
		newVar.accumFn = maxAccumFn
	}
	vc.variables[id] = newVar
	vc.accumVariables[id] = newVar
	return newVar, nil
}

func (vc *variableContext) DataVariableExists(id string) bool {
	_, exists := vc.dataVariables[id]
	return exists
}

func (vc *variableContext) Variable(id string) Variable {
	return vc.variables[id]
}

func (vc *variableContext) DataVariable(id string) DataVariable {
	return vc.dataVariables[id]
}

func (vc *variableContext) AccumVariable(id string) AccumVariable {
	return vc.accumVariables[id]
}

// Performs a topological sort on all the variables in the context
func (vc *variableContext) DependencyOrdering() ([]Variable, error) {
	// first, reset and regenerate dependencies for accumulators based on
	// modifyIds
	for _, accumVar := range vc.accumVariables {
		accumVar.SetDependencyIds([]string{})
	}
	for _, variable := range vc.variables {
		for _, modifyId := range variable.ModifyIds() {
			if modifyVar, exists := vc.accumVariables[modifyId]; exists {
				modifyVar.SetDependencyIds(append(modifyVar.DependencyIds(), variable.Id()))
			} else {
				return nil, fmt.Errorf("No accumulator with id %v", modifyId)
			}
		}
	}

	// now run toposort algorithm
	// our edges are the opposite of dependency, so edge(a -> b) exists iff b
	// depends on a
	sortedList := []Variable{}
	curList := []Variable{} // literals here (as they have no incoming edges)

	edgeIndex := map[string][]Variable{}

	// if markedEdgeCount[node.id] == len(node.dependencies) then we add node to curList
	markedEdgeCount := map[string]int{}

	for _, node := range vc.variables {
		depIds := node.DependencyIds()
		// initialize the startList with all the literals (no dependencies)
		if len(depIds) == 0 {
			curList = append(curList, node)
		}

		// initialize the edgeIndex (node -> things that depend on node)
		for _, depId := range depIds {
			if _, exists := vc.variables[depId]; !exists {
				return nil, fmt.Errorf("Cannot toposort - nonexistent dependency: %v", depId)
			}
			edgeIndex[depId] = append(edgeIndex[depId], node)
		}
	}

	// go through and order everything we can
	for len(curList) > 0 {
		curNode := curList[len(curList)-1]
		curList = curList[:len(curList)-1]
		sortedList = append(sortedList, curNode)

		for _, newNode := range edgeIndex[curNode.Id()] {
			markedEdgeCount[newNode.Id()] += 1
			if markedEdgeCount[newNode.Id()] == len(newNode.DependencyIds()) {
				curList = append(curList, newNode)
			}
		}
	}

	for nodeId, markedEdges := range markedEdgeCount {
		if markedEdges < len(vc.variables[nodeId].DependencyIds()) {
			return nil, fmt.Errorf("Cannot toposort - dependency graph has cycles")
		}
	}

	return sortedList, nil
}

func (vc *variableContext) Variables() map[string]Variable {
	return vc.variables
}
