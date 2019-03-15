package hesabu

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Knetic/govaluate"
	"github.com/gleicon/go-descriptive-statistics"
)

type customFunctionError struct {
	functionName string
	err          string
}

func (e *customFunctionError) Error() string {
	return fmt.Sprintf("Error for %s-function, %s", e.functionName, e.err)
}

// Cache for evalArray-evaluations
var evalExps = make(map[string]*govaluate.EvaluableExpression)

// Functions used by `evalArray`
var functions = map[string]govaluate.ExpressionFunction{
	"ABS":         absFunction,
	"abs":         absFunction,
	"sqrt":        sqrtFunction,
	"SQRT":        sqrtFunction,
	"ACCESS":      accessFunction,
	"access":      accessFunction,
	"ARRAY":       arrayFunction,
	"array":       arrayFunction,
	"AVG":         averageFunction,
	"avg":         averageFunction,
	"stdev":       stdevFunction,
	"STDEV":       stdevFunction,
	"IF":          ifFunction,
	"If":          ifFunction,
	"if":          ifFunction,
	"MAX":         maxFunction,
	"Max":         maxFunction,
	"max":         maxFunction,
	"MIN":         minFunction,
	"Min":         minFunction,
	"min":         minFunction,
	"RANDBETWEEN": randbetweenFunction,
	"randbetween": randbetweenFunction,
	"ROUND":       roundFunction,
	"round":       roundFunction,
	"SAFE_DIV":    safeDivFuntion,
	"Safe_div":    safeDivFuntion,
	"safe_div":    safeDivFuntion,
	"SCORE_TABLE": scoreTableFunction,
	"score_Table": scoreTableFunction,
	"score_table": scoreTableFunction,
	"strlen":      strlen,
	"STRLEN":      strlen,
	"SUM":         sumFunction,
	"Sum":         sumFunction,
	"sum":         sumFunction,
}

func randbetweenFunction(args ...interface{}) (interface{}, error) {
	min := args[0].(float64)
	max := args[1].(float64)
	result := min + rand.Float64()*(max-min)
	return result, nil
}

/*
 SCORE_TABLE = lambda do |*args|
	target = args.shift
	matching_rules = args.each_slice(3).find do |lower, greater, result|
	  greater.nil? || result.nil? ? true : lower <= target && target < greater
	end
	matching_rules.last
  end
*/
func scoreTableFunction(args ...interface{}) (interface{}, error) {
	target := args[0].(float64)
	rules := args[1:]
	chunkSize := 3
	for i := 0; i < len(rules); i += chunkSize {
		end := i + chunkSize

		if end > len(rules) {
			end = len(rules)
		}

		page := rules[i:end]
		if len(page) == 3 {
			lower := page[0].(float64)
			greater := page[1].(float64)
			value := page[2].(float64)
			if lower <= target && target < greater {
				return value, nil
			}
		} else {
			return page[0].(float64), nil
		}
	}

	return args[0], nil
}

func accessFunction(args ...interface{}) (interface{}, error) {
	index := int(args[len(args)-1].(float64))
	return args[index], nil
}

func roundFunction(args ...interface{}) (interface{}, error) {
	places := 0
	if len(args) == 2 {
		places = int(args[1].(float64))
	}
	f := args[0].(float64)
	shift := math.Pow(10, float64(places))
	return (math.Round(f*shift) / shift), nil
}

func absFunction(args ...interface{}) (interface{}, error) {
	return math.Abs(args[0].(float64)), nil
}

func sqrtFunction(args ...interface{}) (interface{}, error) {
	float, ok := args[0].(float64)
	if !ok {
		return nil, &customFunctionError{
			functionName: "SQRT",
			err:          fmt.Sprintf("Expected '%v' to be a float64 expression.", args[0]),
		}
	}
	if float < 0 {
		return nil, &customFunctionError{
			functionName: "SQRT",
			err:          fmt.Sprintf("Expected '%v' to be a float 0 or positive.", args[0]),
		}
	}
	return math.Sqrt(float), nil
}

func ifFunction(args ...interface{}) (interface{}, error) {
	var result interface{}
	bool, ok := args[0].(bool)
	if !ok {
		return nil, &customFunctionError{
			functionName: "IF",
			err:          fmt.Sprintf("Expected '%v' to be a boolean expression.", args[0]),
		}
	}

	if bool {
		result = args[1]
	} else {
		result = args[2]
	}
	return result, nil
}

func safeDivFuntion(args ...interface{}) (interface{}, error) {
	if args[1].(float64) == 0 {
		return float64(0), nil
	}
	return (args[0].(float64) / args[1].(float64)), nil
}

func maxFunction(args ...interface{}) (interface{}, error) {
	max := args[0].(float64)
	for _, arg := range args {
		if arg.(float64) > max {
			max = arg.(float64)
		}
	}
	return max, nil
}

func minFunction(args ...interface{}) (interface{}, error) {
	min := args[0].(float64)
	for _, arg := range args {
		if arg.(float64) < min {
			min = arg.(float64)
		}
	}
	return min, nil
}

func sumFunction(args ...interface{}) (interface{}, error) {
	total := float64(0.0)
	for _, a := range args {
		if v, ok := a.(float64); ok {
			total += v
		} else {
			return nil, &customFunctionError{
				functionName: "sumFunction",
				err:          fmt.Sprintf("Unspoorted type to sum: %T", v),
			}
		}
	}
	return total, nil
}

func stdevFunction(args ...interface{}) (interface{}, error) {
	values := make(descriptive_statistics.Enum, len(args))
	for i := range args {
		values[i] = args[i].(float64)
	}
	return values.StandardDeviation(), nil
}

// A noop function in this context, mainly added for api parity with
// dentaku, so arrays can be explicilty marked as arrays.
//
// ARRAY(1,2,3) => (1,2,3)
func arrayFunction(args ...interface{}) (interface{}, error) {
	return args, nil
}

// `eval_array('a', (1,2,3), 'b', (2,3,4), 'b - a')`
//
// 'a'				=> key 1
// '(1,2,3)'	=> array 1
// 'b'				=> key 2
// '(2,3,4)		=> array 2 (needs to be same length as array 1)
// 'b-a'			=> metaformula
//
// Will loop over the arrays and apply the formula to each index, so
// in this example would result in:
//
//       (2-1, 3-2,4-3)
//       (1,1,1)
func evalArrayFunction(args ...interface{}) (interface{}, error) {
	key1 := args[0].(string)
	array1 := ensureSlice(args[1])
	key2 := args[2].(string)
	array2 := ensureSlice(args[3])
	meta_formula := args[4].(string)

	if len(array1) != len(array2) {
		errorMessage := fmt.Sprintf(
			"Expected '%v' and '%v' to have same size of values (%d and %d)",
			key1,
			key2,
			len(array1),
			len(array2))
		return nil, &customFunctionError{"evalArray", errorMessage}
	}

	var expression *govaluate.EvaluableExpression
	var err error
	if v, ok := evalExps[meta_formula]; ok {
		expression = v
	} else {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(meta_formula, functions)
	}

	if err != nil {
		return nil, &customFunctionError{
			functionName: "EVAL_ARRAY()",
			err:          fmt.Sprintf("Meta formula: %v", err),
		}
	}
	var results []interface{}
	for i, item1 := range array1 {
		item2 := array2[i]
		parameters := make(map[string]interface{}, 2)
		parameters[key1] = item1
		parameters[key2] = item2
		result, error_eval := expression.Evaluate(parameters)
		if error_eval != nil {
			return nil, &customFunctionError{
				functionName: "EVAL_ARRAY()",
				err:          fmt.Sprintf("%v. We only know '%s' and '%s'", error_eval, key1, key2),
			}
		}
		results = append(results, result)
	}

	return results, nil
}

func averageFunction(args ...interface{}) (interface{}, error) {
	total := float64(0)
	for _, x := range args {
		total += x.(float64)
	}
	return (total / float64(len(args))), nil
}

func strlen(args ...interface{}) (interface{}, error) {
	length := len(args[0].(string))
	return (float64)(length), nil
}

// Ensures that the interface passed is a slice, it's like Array.wrap
// but in golang.
func ensureSlice(arg interface{}) []interface{} {
	arr, ok := arg.([]interface{})
	if !ok {
		arr = make([]interface{}, 1)
		arr[0] = arg
	}
	return arr
}

// Functions return function registry
func Functions() map[string]govaluate.ExpressionFunction {
	all_functions := functions
	all_functions["eval_array"] = evalArrayFunction
	all_functions["EVAL_ARRAY"] = evalArrayFunction
	return all_functions
}
