package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BLSQ/go-hesabu/hesabu"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var debugFlag = flag.Bool("d", false, "Extra debug logging")
var versionFlag = flag.Bool("v", false, "Prints version")

func init() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	if os.Getenv("HESABU_DEBUG") == "true" || *debugFlag {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	raw, error := getInput(flag.Args())
	if error != nil {
		fmt.Printf(`
You need to either supply a filename or pipe to hesabu

      bin/hesabucli path/to/yourfilename.json
      echo '{"a": "1 + 2 + b", "b": "7"}' | bin/hesabucli
`)
		os.Exit(1)
	}

	rawEquations, error := getEquations(raw)
	if error != nil {
		errs := []hesabu.EvalError{
			{
				Message:    "Invalid JSON",
				Source:     "general",
				Expression: "general",
			},
		}
		logErrors(errs)
		os.Exit(1)
	}

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

func getInput(flag_arguments []string) ([]byte, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	var str []byte
	if fi.Mode()&os.ModeNamedPipe == 0 {
		if len(flag_arguments) < 1 {
			return nil, errors.New("No filename supplied")
		}
		raw, err := ioutil.ReadFile(flag_arguments[0])
		if err != nil {
			return nil, err
		}
		str = raw
	} else {
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		str = raw
	}
	return str, nil
}

func getEquations(raw []byte) (map[string]string, error) {
	log.Printf("equations to parse %s", string(raw))
	var results map[string]string
	err := json.Unmarshal(raw, &results)
	if err != nil {
		log.Printf("equations not loaded %v ", err)
		return nil, err
	}
	log.Printf("equations loaded: %d ", len(results))
	return results, nil
}
