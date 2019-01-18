package hesabu

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/Knetic/govaluate"
)

var functions = map[string]govaluate.ExpressionFunction{
	"strlen":      strlen,
	"if":          ifFunction,
	"IF":          ifFunction,
	"If":          ifFunction,
	"max":         maxFunction,
	"MAX":         maxFunction,
	"Max":         maxFunction,
	"min":         minFunction,
	"MIN":         minFunction,
	"Min":         minFunction,
	"safe_div":    safeDivFuntion,
	"SAFE_DIV":    safeDivFuntion,
	"Safe_div":    safeDivFuntion,
	"sum":         sumFunction,
	"SUM":         sumFunction,
	"Sum":         sumFunction,
	"avg":         averageFunction,
	"AVG":         averageFunction,
	"ABS":         absFunction,
	"abs":         absFunction,
	"access":      accessFunction,
	"ACCESS":      accessFunction,
	"score_table": scoreTableFunction,
	"score_Table": scoreTableFunction,
	"SCORE_TABLE": scoreTableFunction,
	"round":       roundFunction,
	"ROUND":       roundFunction,
	"randbetween": randbetweenFunction,
	"RANDBETWEEN": randbetweenFunction,
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

func ifFunction(args ...interface{}) (interface{}, error) {
	var result interface{}
	bool, ok := args[0].(bool)
	if !ok {
		return nil, fmt.Errorf("Expected '%v' to be a boolean expression.", args[0])
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
	for _, arg := range args {
		total += arg.(float64)
	}
	return total, nil
}

func wrap(arg interface{}) []interface{} {
	arr, ok := arg.([]interface{})
	if !ok {
		arr = make([]interface{}, 1)
		arr[0] = arg
	}
	return arr
}

var evalExps = make(map[string]*govaluate.EvaluableExpression)

func evalArrayFunction(args ...interface{}) (interface{}, error) {
	key1 := args[0].(string)
	array1 := wrap(args[1])
	key2 := args[2].(string)
	array2 := wrap(args[3])
	if len(array1) != len(array2) {
		return nil, fmt.Errorf("Expected '%v' and '%v' to have same size of values (%v => %d and %v => %d)", key1, key2, array1, len(array1), array2, len(array2))
	}
	meta_formula := args[4].(string)
	var expression *govaluate.EvaluableExpression
	var err error
	if v, aok := evalExps[meta_formula]; aok {
		expression = v
	} else {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(meta_formula, functions)
	}

	log.Printf("Expression error %v", err)
	if err != nil {
		return nil, err
	}
	var results []interface{}
	for i, item1 := range array1 {
		item2 := array2[i]
		parameters := make(map[string]interface{}, 2)
		parameters[key1] = item1
		parameters[key2] = item2
		result, error_eval := expression.Evaluate(parameters)
		log.Printf("result %v, error_eval %v", result, error_eval)
		if error_eval != nil {
			return nil, error_eval
		}
		results = append(results, result)
	}

	// result := meta_formula
	log.Printf("key 1 %v with %v\nkey2 %v with %v\nsuch meta: %v\nexpression%v", key1, array1, key2, array2, meta_formula, expression)
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

// Functions return function registry
func Functions() map[string]govaluate.ExpressionFunction {

	functions := map[string]govaluate.ExpressionFunction{
		"strlen":      strlen,
		"if":          ifFunction,
		"IF":          ifFunction,
		"If":          ifFunction,
		"max":         maxFunction,
		"MAX":         maxFunction,
		"Max":         maxFunction,
		"min":         minFunction,
		"MIN":         minFunction,
		"Min":         minFunction,
		"safe_div":    safeDivFuntion,
		"SAFE_DIV":    safeDivFuntion,
		"Safe_div":    safeDivFuntion,
		"sum":         sumFunction,
		"SUM":         sumFunction,
		"Sum":         sumFunction,
		"avg":         averageFunction,
		"AVG":         averageFunction,
		"ABS":         absFunction,
		"abs":         absFunction,
		"access":      accessFunction,
		"ACCESS":      accessFunction,
		"score_table": scoreTableFunction,
		"score_Table": scoreTableFunction,
		"SCORE_TABLE": scoreTableFunction,
		"round":       roundFunction,
		"ROUND":       roundFunction,
		"randbetween": randbetweenFunction,
		"RANDBETWEEN": randbetweenFunction,
		"eval_array":  evalArrayFunction,
	}
	return functions
}
