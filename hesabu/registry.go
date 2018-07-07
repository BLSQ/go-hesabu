package hesabu

import (
	"github.com/Knetic/govaluate"
)

func Functions() map[string]govaluate.ExpressionFunction {
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
		"if":       ifFunction,
		"IF":       ifFunction,
		"If":       ifFunction,
		"max":      maxFunction,
		"MAX":      maxFunction,
		"Max":      maxFunction,
		"min":      minFunction,
		"MIN":      minFunction,
		"Min":      minFunction,
		"safe_div": safeDivFuntion,
		"sum":      sumFunction,
		"SUM":      sumFunction,
		"Sum":      sumFunction,
		"avg":      averageFunction,
		"AVG":      averageFunction,
		"ABS":      averageFunction,
	}
	return functions
}
