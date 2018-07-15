package hesabu

import (
	"math"
	"math/rand"

	"github.com/Knetic/govaluate"
)

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
	rules := args[1:len(args)]
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

// Functions return function registry
func Functions() map[string]govaluate.ExpressionFunction {

	accessFunction := func(args ...interface{}) (interface{}, error) {
		index := int(args[len(args)-1].(float64))
		return args[index], nil
	}

	absFunction := func(args ...interface{}) (interface{}, error) {
		return math.Abs(args[0].(float64)), nil
	}
	roundFunction := func(args ...interface{}) (interface{}, error) {
		places := 0
		if len(args) == 2 {
			places = int(args[1].(float64))
		}
		f := args[0].(float64)
		shift := math.Pow(10, float64(places))
		return (math.Round(f*shift) / shift), nil
	}

	ifFunction := func(args ...interface{}) (interface{}, error) {
		var result interface{}
		if args[0].(bool) {
			result = args[1]
		} else {
			result = args[2]
		}
		return result, nil
	}
	maxFunction := func(args ...interface{}) (interface{}, error) {
		max := args[0].(float64)
		for _, arg := range args {
			if arg.(float64) > max {
				max = arg.(float64)
			}
		}
		return max, nil
	}
	minFunction := func(args ...interface{}) (interface{}, error) {
		min := args[0].(float64)
		for _, arg := range args {
			if arg.(float64) < min {
				min = arg.(float64)
			}
		}
		return min, nil
	}
	sumFunction := func(args ...interface{}) (interface{}, error) {
		total := float64(0.0)
		for _, arg := range args {
			total += arg.(float64)
		}
		return total, nil
	}

	averageFunction := func(args ...interface{}) (interface{}, error) {
		total := float64(0)
		for _, x := range args {
			total += x.(float64)
		}
		return (total / float64(len(args))), nil
	}

	safeDivFuntion := func(args ...interface{}) (interface{}, error) {
		if args[1].(float64) == 0 {
			return float64(0), nil
		}
		return (args[0].(float64) / args[1].(float64)), nil
	}
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
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
	return functions
}
