package main

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"

	toposort "github.com/philopon/go-toposort"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
		"if": func(args ...interface{}) (interface{}, error) {
			var result interface{}
			if args[0].(bool) {
				result = args[1]
			} else {
				result = args[3]
			}
			return result, nil
		},
		"max": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
		"safe_div": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
	}

	rawEquations := getEquations()

	// https://github.com/philopon/go-toposort/blob/master/toposort.go
	graph := toposort.NewGraph(len(rawEquations))

	equations := make(map[string]govaluate.EvaluableExpression, len(rawEquations))

	for key, exp := range rawEquations {
		graph.AddNode(key)

		// https://github.com/Knetic/govaluate/blob/master/EvaluableExpression.go
		expression, err := govaluate.NewEvaluableExpressionWithFunctions(exp, functions)
		if err != nil {
			log.Fatalf("error %#v", err)
			os.Exit(1)
		}
		equations[key] = *expression

		log.Printf("vars  %v", expression.Vars())
		for _, dependency := range expression.Vars() {
			graph.AddNode(dependency)
			graph.AddEdge(key, dependency)
		}
	}

	topsort, ok := graph.Toposort()

	if !ok {
		panic("cycle detected")
	}
	fmt.Println(topsort)

	reverse(topsort)
	parameters := make(map[string]interface{}, len(rawEquations))
	fmt.Println(topsort)
	for _, key := range topsort {
		result, err := equations[key].Evaluate(parameters)
		if err != nil {
			log.Fatalf("error %v", err)
		}

		v, ok := result.(float64)
		if !ok {
			log.Printf("%v is not float64", result)
		} else {
			log.Printf("%s = %v (%s)", key, v, rawEquations[key])
			parameters[key] = v
		}

	}

}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func getEquations() map[string]string {
	raw, err := ioutil.ReadFile("./small.json")
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
