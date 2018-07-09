package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"./hesabu"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func main() {

	rawEquations := getEquations(os.Args[1])
	parsedEquations := hesabu.Parse(rawEquations, hesabu.Functions())
	log.Printf("during parsing %v ", parsedEquations.Errors)
	if len(parsedEquations.Errors) > 0 {
		logErrors(parsedEquations.Errors)
	} else {
		solutions, err := parsedEquations.Solve()
		if err != nil {
			var evalErrors []hesabu.EvalError
			var hesabuerr hesabu.EvalError

			ok, err2 := err.(*hesabu.CustomError)
			if !err2 {
				panic("ddd")
			}
			if ok != nil {
				hesabuerr = ok.EvalError
			}

			evalErrors = append(evalErrors, hesabuerr)
			logErrors(evalErrors)
		} else {
			logSolution(solutions)
		}
	}

}

func logErrors(errors []hesabu.EvalError) {
	log.Printf("during parsing %v ", errors)
	var content = make(map[string]interface{}, 1)
	content["errors"] = errors
	b, _ := json.MarshalIndent(content, "", "  ")
	s := string(b)
	fmt.Println(s)
}

func logSolution(solutions map[string]interface{}) {
	b, _ := json.MarshalIndent(solutions, "", "  ")
	s := string(b)
	fmt.Println(s)
}

func getEquations(file string) map[string]string {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
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
