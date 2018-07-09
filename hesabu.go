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
	if os.Getenv("HESABU_DEBUG") == "true" {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {

	rawEquations := getEquations()
	parsedEquations := hesabu.Parse(rawEquations, hesabu.Functions())
	if len(parsedEquations.Errors) > 0 {
		logErrors(parsedEquations.Errors)
		os.Exit(1)
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
			os.Exit(1)
		} else {
			logSolution(solutions)
		}
	}

}

func logErrors(errors []hesabu.EvalError) {
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

func getEquations() map[string]string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var str []byte
	if fi.Mode()&os.ModeNamedPipe == 0 {
		raw, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic("file not read" + os.Args[1])
		}
		str = raw
	} else {
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic("pipe not read")
		}
		str = raw
	}
	log.Printf("equations to parse %s", string(str))
	var results map[string]string
	err = json.Unmarshal(str, &results)
	if err != nil {
		log.Printf("equations not loaded %v ", err)
	}
	log.Printf("equations loaded: %d ", len(results))
	return results
}
