package simpleProperties

import (
	"container/list"
	"strings"
)

// BasicEvaluator Try to evaluate all properties
// Step 1 - Update all properties using know values but ignoring defaults
// Step 2 - Once no more can be evaluated, select one default, and repeat step 1
//
//	2a - When selecting a default, check for properties that don't appear in the lhs of any expression
//	2b - If no default is available for (2a), use the first one that is available
//
// Step 3 - Once no more evaluations can be made, fail if unevaluated properties still exist, else all good
func BasicEvaluator() func(*Properties) {
	return func(p *Properties) {
		for true {
			var changed bool
			unresolved := p.GetEvalKeys()
			for _, lhsName := range unresolved {
				itemsList := p.evalExprMap[lhsName]
				var next *list.Element
				for element := itemsList.Front(); element != nil; element = next {
					next = element.Next()
					item := element.Value.(*exprParts)
					// do we have an existing property value for this ?
					value := p.GetProperty(item.name)
					changed = changed || doEvaluation(p, value, itemsList, element, lhsName, item, false)
				}
			}
			// start looking at defaults
			if !changed {
				unresolved = p.GetEvalKeys()
				for evalCheck := 0; evalCheck < 2; evalCheck++ { // 0 == check for lhs evaluator, 1 = don't check, use whatever is available
					if changed {
						break
					}
					// look at all unresolved items
					for _, lhsName := range unresolved {
						if changed {
							break
						}
						// do we have a default property value for this ?
						itemsList := p.evalExprMap[lhsName]
						// go through all unresolved rhs items until we run out of elements
						var next *list.Element
						for element := itemsList.Front(); element != nil; element = next {
							next = element.Next()
							item := element.Value.(*exprParts)
							name := item.name
							if evalCheck == 0 && hasPotentialEvaluator(p, name) { // may yet get evaluated. Ignore until other defaults expended
								continue
							}
							defaultValue := item.defaultValue
							changed = doEvaluation(p, defaultValue, itemsList, element, lhsName, item, true)
							if changed {
								break
							}
						}
					}
				}
			}
			if !changed {
				break
			}
		}
	}
}

// does the named value have a potential evaluator, e.g. abc = ${something} ?
// if so, this should be used in default value assignment only after all named values without potential assignments
// have been used first
func hasPotentialEvaluator(p *Properties, name string) bool {
	evalKeys := p.GetEvalKeys()
	for _, key := range evalKeys {
		if key == name {
			return true
		}
	}
	return false
}

// if we have resolved something -->
// a) remove it from the list of things to resolve
// b) replace all the placeholders in the rhs with the resolved value
// c) if the lhs is fully resolved, put it into the remove it from the to be evaluated map and put it into the kv map
func doEvaluation(p *Properties, resolvedValue string, itemsList *list.List, element *list.Element, lhsName string, item *exprParts, defaultReplacement bool) bool {
	if resolvedValue != "" {
		// update rhs expression
		itemsList.Remove(element)
		rhs := p.evalKeyValueMap[lhsName]
		toBeReplaced := item.full
		var replaceQuantity int
		if defaultReplacement {
			replaceQuantity = 1 // so we don't replace all with same default value
		} else {
			replaceQuantity = -1 // its not a default value, i.e. resolved lhs so safe to replace all
		}
		evaluatedRhs := strings.Replace(rhs, toBeReplaced, resolvedValue, replaceQuantity)
		if containsExpression(evaluatedRhs) {
			// not fully evaluated so just update partially resolved expression
			p.evalKeyValueMap[lhsName] = evaluatedRhs
		} else {
			// finished so remove from expr valuation data and move to resolved properties
			delete(p.evalKeyValueMap, lhsName)
			delete(p.evalExprMap, lhsName)
			p.keyValueMap[lhsName] = evaluatedRhs
			return true
		}
	}
	return false
}
