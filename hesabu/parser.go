package hesabu

import (
	"fmt"
	"github.com/Knetic/govaluate"
	toposort "github.com/otaviokr/topological-sort"
	"log"
	"os"
	"regexp"
)

// ParsedEquations raw equation, EvaluableExpression and dependencies
type ParsedEquations struct {
	RawEquations map[string]string
	Equations    map[string]govaluate.EvaluableExpression
	Dependencies map[string][]string
}

// Parse string equation in a EvaluableExpressions and their dependencies
func Parse(rawEquations map[string]string, functions map[string]govaluate.ExpressionFunction) ParsedEquations {

	equations := make(map[string]govaluate.EvaluableExpression, len(rawEquations))
	equationDependencies := make(map[string][]string, len(rawEquations))
	for key, exp := range rawEquations {
		mutiple := regexp.MustCompile(`(=)+`)
		single := regexp.MustCompile(`=`)
		interFixedExpression := mutiple.ReplaceAllString(exp, "=")
		fixedExpression := single.ReplaceAllString(interFixedExpression, "==")
		if fixedExpression != exp {
			log.Printf("fixed \n\t%s vs \n\t%s", exp, fixedExpression)
		}
		// https://github.com/Knetic/govaluate/blob/master/EvaluableExpression.go
		expression, err := govaluate.NewEvaluableExpressionWithFunctions(fixedExpression, functions)
		if err != nil {
			log.Fatalf("error %#v %#v %#v", err, key, exp)
			fmt.Printf("error")
			os.Exit(1)
		}
		equations[key] = *expression

		log.Printf("vars  %v", expression.Vars())
		equationDependencies[key] = expression.Vars()
	}
	return ParsedEquations{Equations: equations, Dependencies: equationDependencies, RawEquations: rawEquations}
}

// Solve the equation in correct order and return map of values
func (parsedEquations ParsedEquations) Solve() map[string]interface{} {
	topsort := toposort.ReversedSort(parsedEquations.Dependencies)

	log.Println(topsort)

	solutions := make(map[string]interface{}, len(parsedEquations.RawEquations))
	log.Println("topsort %v", topsort)
	for _, key := range topsort {
		result, err := parsedEquations.Equations[key].Evaluate(solutions)
		if err != nil {
			log.Fatalf("error %v", err)
		}

		v, ok := result.(float64)
		if !ok {
			log.Printf("%v is not float64, %v (%s)", result, key, parsedEquations.RawEquations[key])
		} else {
			log.Printf("%s = %v (%s)", key, v, parsedEquations.RawEquations[key])
			solutions[key] = v
		}
	}
	return solutions
}
