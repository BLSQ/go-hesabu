package hesabu

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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

var ShouldLog = false

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
			equationDependencies[key] = expression.Vars()
		}
	}
	return ParsedEquations{Equations: equations, Dependencies: equationDependencies, RawEquations: rawEquations, Errors: errorsCollector}
}

// Solve the equation in correct order and return map of values
func (parsedEquations ParsedEquations) Solve() (map[string]any, error) {
	topsort := toposort.ReversedSort(parsedEquations.Dependencies)

	solutions := make(map[string]any, len(parsedEquations.RawEquations))
	if len(topsort) == 0 {
		evalError := EvalError{Message: "cycle between equations", Source: "general", Expression: "general"}
		return nil, &CustomError{EvalError: evalError}
	}

	for _, key := range topsort {
		expression, ok := parsedEquations.Equations[key]
		if !ok {
			// Key is missing an expression, we allow it to go on because
			// Evaluate will produce a better error message when it needs
			// this key.
			continue
		}
		result, err := expression.Evaluate(solutions)

		if err != nil {
			return parsedEquations.newSingleError(key, err.Error())
		}

		if v, ok := result.(float64); ok {
			if math.IsInf(v, 0) {
				return parsedEquations.newSingleError(key, "Divide by zero")
			} else if math.IsNaN(v) {
				return parsedEquations.newSingleError(key, "NaN")
			} else {
				solutions[key] = v
			}
		} else {
			solutions[key] = result
		}
	}

	return solutions, nil
}

func (parsedEquations ParsedEquations) newSingleError(key string, message string) (map[string]any, error) {
	equation := parsedEquations.RawEquations[key]
	evalError := EvalError{Message: message, Source: key, Expression: equation}
	return nil, &CustomError{EvalError: evalError}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// There can still be legacy formulas that use the old AND, OR and '='
// syntax, the clean method will replace this with:
//       AND => &&
//       OR => ||
//       = => == (only for equality comparison)
//
func clean(expression string) (cleanExpression string) {
	if !needsCleaning(expression) {
		return expression
	}
	cleanExpression = expression
	cleanExpression = strings.Replace(cleanExpression, " AND ", " && ", -1)
	cleanExpression = strings.Replace(cleanExpression, " OR ", " || ", -1)
	cleanExpression = strings.Replace(cleanExpression, " and ", " && ", -1)
	cleanExpression = strings.Replace(cleanExpression, " or ", " || ", -1)

	if strings.Contains(cleanExpression, "=") {
		cleanExpression = replaceSingleEquals(cleanExpression)
	}

	return cleanExpression
}

func needsCleaning(expression string) bool {
	// Numbers never need cleaning
	if isNumeric(expression) {
		return false
	}

	// Invalid formulas
	if strings.HasPrefix(expression, "=") {
		return false
	}
	if strings.HasSuffix(expression, "=") {
		return false
	}
	return true
}

// Replaces a single '=' with its double cousing '==', it takes care
// to not replace single equals that are actually a different symbol
// ('==', '<=', '>=')
//
// Why don't you just use a regex? Turns out, regex in go is not that
// fast. A regex replace for the single equals was on average almost a
// third slower than not using a regex.
func replaceSingleEquals(in string) string {
	// Characters that combined with an equal form a special symbol
	reserved := map[byte]int{'=': 1, '<': 1, '>': 1, '!': 1}
	var t = []rune{}

	// Loop over the runes from the in string
	for pos, char := range in {
		t = append(t, char)
		if char == '=' {
			// Make sure that neither the next, or the previous character is
			// one of the reserveds
			_, previousIsReserved := reserved[in[pos-1]]
			_, nextIsReserved := reserved[in[pos+1]]
			if !(previousIsReserved || nextIsReserved) {
				t = append(t, '=')
			}
		}
	}
	// Convert runes back to string
	return string(t)
}
