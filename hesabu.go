package main

import (
	"encoding/json"
	"fmt"

	"./hesabu"

	"io/ioutil"
	"log"
	"os"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

func main() {

	rawEquations := getEquations(os.Args[1])

	functions := hesabu.Functions()

	parsedEquations := hesabu.Parse(rawEquations, functions)

	solutions := parsedEquations.Solve()
	logSolution(solutions)
}

func logSolution(solutions map[string]interface{}) {
	b, _ := json.MarshalIndent(solutions, "", "  ")
	// Convert bytes to string.
	s := string(b)
	fmt.Println(s)
}

func getEquations(file string) map[string]string {
	raw, err := ioutil.ReadFile(file)
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
