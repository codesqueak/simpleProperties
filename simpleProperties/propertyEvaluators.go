package simpleProperties

import (
	"container/list"
	"strings"
)

func DefaultEvaluator() func(*Properties) {
	return func(p *Properties) {
		println("-- expression evaluator --")
	}
}

// BasicEvaluator Try to evaluate all properties
// Step 1 - Update all properties using know values but ignroing defaults
// Step 2 - Once no more can be evaluated, select one default, and repeat step 1
// Step 3 - Once no more evaluations can be made, fail if unevaluated properties still exist, else all good
func BasicEvaluator() func(*Properties) {
	return func(p *Properties) {
		for true {
			println("-- basic evaluator - step 1 evaluate --")
			var changed bool
			unresolved := p.GetEvalKeys()
			for _, lhsName := range unresolved {
				itemsList := p.evalExprMap[lhsName]
				element := itemsList.Front()
				item := element.Value.(*exprParts)
				// do we have an existing property value for this ?
				value := p.GetProperty(item.name)
				changed = doEvaluation(p, value, itemsList, element, lhsName, item)
			}
			println("-- basic evaluator - step 1 complete --")
			if !changed {
				println("-- basic evaluator - step 2 use default --")
				unresolved = p.GetEvalKeys()
				for _, lhsName := range unresolved {
					// do we have a default property value for this ?
					itemsList := p.evalExprMap[lhsName]
					element := itemsList.Front()
					item := element.Value.(*exprParts)
					defaultValue := item.defaultValue
					changed = doEvaluation(p, defaultValue, itemsList, element, lhsName, item)
				}
				println("-- basic evaluator - step 2 complete --")
			}
			if !changed {
				break
			}
		}
	}
}

func doEvaluation(p *Properties, value string, itemsList *list.List, element *list.Element, lhsName string, item *exprParts) bool {
	if value != "" {
		// update rhs expression
		itemsList.Remove(element)
		rhs := p.evalKeyValueMap[lhsName]
		toBeReplaced := item.full
		println("Substituting ", value, "for", toBeReplaced, "in", rhs)
		evaluatedRhs := strings.Replace(rhs, toBeReplaced, value, -1)
		if containsExpression(evaluatedRhs) {
			// not fully evaluated so just update partially resolved expression
			p.evalKeyValueMap[lhsName] = evaluatedRhs
		} else {
			// finished so remove from expr valuation data and move to resolved properties
			delete(p.evalKeyValueMap, lhsName)
			delete(p.evalExprMap, lhsName)
			p.keyValueMap[lhsName] = evaluatedRhs
		}
		return true
	}
	return false
}
