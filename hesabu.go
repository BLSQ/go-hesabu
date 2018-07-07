package main

import (
	"./hesabu"
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetOutput(ioutil.Discard)
}
func handler(w http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	//log.Println(string(body))
	var reqRawEquations map[string]interface{}
	err = json.Unmarshal(body, &reqRawEquations)
	if err != nil {
		log.Printf("equations not loaded %v ", err)
	}
	rawEquations := make(map[string]string, len(reqRawEquations))
	for k, v := range reqRawEquations {
		rawEquations[k] = fmt.Sprintf("%v", v)
	}
	parsedEquations := hesabu.Parse(rawEquations, hesabu.Functions())

	solutions := parsedEquations.Solve()
	b, _ := json.MarshalIndent(solutions, "", "  ")
	s := string(b)
	fmt.Fprintf(w, s)
}

func main() {

	if os.Args[1] == "s" {
		withoutGz := http.HandlerFunc(handler)
		withGz := gziphandler.GzipHandler(withoutGz)
		http.Handle("/", withGz)
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		rawEquations := getEquations(os.Args[1])
		parsedEquations := hesabu.Parse(rawEquations, hesabu.Functions())
		solutions := parsedEquations.Solve()
		logSolution(solutions)
	}
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
