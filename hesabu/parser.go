package hesabu

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Knetic/govaluate"
	toposort "github.com/otaviokr/topological-sort"
)

// ParsedEquations raw equation, EvaluableExpression and dependencies
type ParsedEquations struct {
	RawEquations map[string]string
	Equations    map[string]govaluate.EvaluableExpression
	Dependencies map[string][]string
	Errors       []EvalError
}

type CustomError struct {
	EvalError EvalError
}

func (e *CustomError) Error() string {
	return e.EvalError.Error()
}

// Eval or parsing errors
type EvalError struct {
	Source     string `json:"source"`
	Expression string `json:"expression"`
	Message    string `json:"message"`
}

func (e *EvalError) Error() string {
	return fmt.Sprintf("%s %s %s", e.Message, e.Expression, e.Source)
}

// Parse string equation in a EvaluableExpressions and their dependencies
func Parse(rawEquations map[string]string, functions map[string]govaluate.ExpressionFunction) ParsedEquations {
	var errorsCollector []EvalError
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
			errorsCollector = append(errorsCollector, EvalError{Source: key, Message: err.Error(), Expression: fixedExpression})
		} else {
			equations[key] = *expression
			log.Printf("vars  %v", expression.Vars())
			equationDependencies[key] = expression.Vars()
		}
	}
	return ParsedEquations{Equations: equations, Dependencies: equationDependencies, RawEquations: rawEquations, Errors: errorsCollector}
}

// Solve the equation in correct order and return map of values
func (parsedEquations ParsedEquations) Solve() (map[string]interface{}, error) {
	topsort := toposort.ReversedSort(parsedEquations.Dependencies)

	solutions := make(map[string]interface{}, len(parsedEquations.RawEquations))
	log.Println("topsort %v", topsort)
	for _, key := range topsort {
		result, err := parsedEquations.Equations[key].Evaluate(solutions)
		if err != nil {
			equation := parsedEquations.RawEquations[key]
			evalError := EvalError{Message: "evaluate error " + err.Error(), Source: key, Expression: equation}
			log.Printf("%s => %v (%s)", key, err.Error(), parsedEquations.RawEquations[key])
			return make(map[string]interface{}), &CustomError{EvalError: evalError}
		}

		v, ok := result.(float64)
		if !ok {
			log.Printf("%v is not float64, %v (%s)", result, key, parsedEquations.RawEquations[key])
		} else {
			log.Printf("%s = %v (%s)", key, v, parsedEquations.RawEquations[key])
			solutions[key] = v
		}

	}
	return solutions, nil
}
