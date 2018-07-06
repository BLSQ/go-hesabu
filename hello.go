package main

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"

	toposort "github.com/otaviokr/topological-sort"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

func main() {
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
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
		"if":  ifFunction,
		"IF":  ifFunction,
		"If":  ifFunction,
		"max": maxFunction,
		"MAX": maxFunction,
		"Max": maxFunction,
		"min": minFunction,
		"MIN": minFunction,
		"Min": minFunction,
		"safe_div": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
		"sum": sumFunction,
		"SUM": sumFunction,
		"Sum": sumFunction,
		"avg": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
		"AVG": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
		"ABS": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
	}

	rawEquations := getEquations()

	// https://github.com/philopon/go-toposort/blob/master/toposort.go
	//graph := toposort.NewGraph() //(10 * len(rawEquations))

	equations := make(map[string]govaluate.EvaluableExpression, len(rawEquations))
	equationDependencies := make(map[string][]string, len(rawEquations))
	for key, exp := range rawEquations {
		//graph.AddNode(key)

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
		//for _, dependency := range expression.Vars() {
		//graph.AddNode(dependency)
		//graph.AddEdge(key, dependency)
		//}
	}

	topsort := toposort.ReversedSort(equationDependencies)

	//if !ok {
	//	panic("cycle detected")
	//}
	log.Println(topsort)

	//reverse(topsort)
	solutions := make(map[string]interface{}, len(rawEquations))
	log.Println("topsort %v", topsort)
	for _, key := range topsort {
		result, err := equations[key].Evaluate(solutions)
		if err != nil {
			log.Fatalf("error %v", err)
		}

		v, ok := result.(float64)
		if !ok {
			log.Printf("%v is not float64, %v (%s)", result, key, rawEquations[key])
		} else {
			log.Printf("%s = %v (%s)", key, v, rawEquations[key])
			solutions[key] = v
		}
	}

	b, _ := json.MarshalIndent(solutions, "", "  ")
	// Convert bytes to string.
	s := string(b)
	fmt.Println(s)
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func getEquations() map[string]string {
	raw, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var results map[string]string
	err = json.Unmarshal(raw, &results)
	if err != nil {
		log.Printf("equations not loaded %v ", err)
	}
	log.Printf("equations loaded: %d ", len(results))
	//log.Println("map:", results)
	return results
}
