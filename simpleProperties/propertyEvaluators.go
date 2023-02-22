package simpleProperties

func DefaultEvaluator(_ string) func(*Properties) {
	return func(p *Properties) {
		println("-- default evaluator --")
	}
}
