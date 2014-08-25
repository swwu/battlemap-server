package variable

import (
	"fmt"

	"github.com/swwu/v8.go"
)

type VariableContext interface {
	Variable(id string) Variable
	SetVariable(id string, dependencies []string, onEvalFn *v8.Function) (Variable, error)

	DependencyOrdering() ([]Variable, error)
	//Eval()
}

type variableContext struct {
	variables map[string]Variable
}

func NewContext() VariableContext {
	return &variableContext{
		variables: map[string]Variable{},
	}
}

func (vc *variableContext) SetVariable(id string,
	dependencyIds []string, onEvalFn *v8.Function) (Variable, error) {
	if _, exists := vc.variables[id]; exists {
		return nil, fmt.Errorf("Variable with id", id, "already exists")
	}

	newVar := &scriptVariable{
		id:            id,
		context:       vc,
		dependencyIds: dependencyIds,
		onEvalFn:      onEvalFn,
	}
	vc.variables[id] = newVar
	return newVar, nil
}

func (vc *variableContext) Variable(id string) Variable {
	return vc.variables[id]
}

// performs a topological sort on all the variables in the context
func (vc *variableContext) DependencyOrdering() ([]Variable, error) {
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