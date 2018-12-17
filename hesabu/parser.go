package hesabu

import (
	"fmt"
	"log"
	"math"
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

// returned as err
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
		exp = clean(exp)
		// https://github.com/Knetic/govaluate/blob/master/EvaluableExpression.go
		expression, err := govaluate.NewEvaluableExpressionWithFunctions(exp, functions)
		if err != nil {
			errorsCollector = append(errorsCollector, EvalError{Source: key, Message: err.Error(), Expression: exp})
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
	if len(topsort) == 0 {
		evalError := EvalError{Message: "cycle between equations", Source: "general", Expression: "general"}
		return make(map[string]interface{}), &CustomError{EvalError: evalError}
	}

	for _, key := range topsort {

		result, err := parsedEquations.Equations[key].Evaluate(solutions)
		if err != nil {
			return parsedEquations.newSingleError(key, err.Error())
		}

		v, ok := result.(float64)
		if !ok {
			log.Printf("%v is not float64, %v (%s)", result, key, parsedEquations.RawEquations[key])
		} else {
			log.Printf("%s = %v (%s)", key, v, parsedEquations.RawEquations[key])
			if math.IsInf(v, 0) {
				return parsedEquations.newSingleError(key, "Divide by zero")
			} else {
				solutions[key] = v
			}
		}

		vBool, okBool := result.(bool)
		if okBool {
			solutions[key] = vBool
		}
		vString, okString := result.(string)
		if okString {
			solutions[key] = vString
		}
	}
	return solutions, nil
}

func (parsedEquations ParsedEquations) newSingleError(key string, message string) (map[string]interface{}, error) {
	equation := parsedEquations.RawEquations[key]
	evalError := EvalError{Message: message, Source: key, Expression: equation}
	return make(map[string]interface{}), &CustomError{EvalError: evalError}
}

// There can still be legacy formulas that use the old AND, OR and '='
// syntax, the clean method will replace this with:
//       AND => &&
//       OR => ||
//       = => == (only for equality comparison)
//
func clean(expression string) (cleanExpression string) {
	and_regex := regexp.MustCompile(`(?i)\bAND\b`)
	or_regex := regexp.MustCompile(`(?i)\bOR\b`)
	single_equals_regex := regexp.MustCompile(`(\b|\s)(=)(\b|\s)`)

	cleanExpression = expression
	cleanExpression = and_regex.ReplaceAllString(cleanExpression, "&&")
	cleanExpression = or_regex.ReplaceAllString(cleanExpression, "||")
	cleanExpression = single_equals_regex.ReplaceAllString(cleanExpression, "==")

	return cleanExpression
}
