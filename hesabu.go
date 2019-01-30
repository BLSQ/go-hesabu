package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BLSQ/go-hesabu/hesabu"
)

var (
	version = "0.0.4"
	commit  = "none"
	date    = "20190121"
)

var debugFlag = flag.Bool("d", false, "Extra debug logging")
var versionFlag = flag.Bool("v", false, "Prints version")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var ShouldLog = false

func init() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	if os.Getenv("HESABU_DEBUG") == "true" || *debugFlag {
		ShouldLog = true
		hesabu.ShouldLog = ShouldLog
		log.SetOutput(os.Stderr)
	}

	if *cpuprofile != "" {
		startProfilingCPU(*cpuprofile)
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
	}

	solutions, err := parsedEquations.Solve()
	if err != nil {
		if customError, ok := err.(*hesabu.CustomError); ok {
			evalError := customError.EvalError
			logErrors([]hesabu.EvalError{evalError})
			os.Exit(1)
		} else {
			panic("Only expected a custom error")
		}
	}

	logSolution(solutions)

	stopProfilingCPU()
	if *memprofile != "" {
		startProfilingMemory(*memprofile)
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
	b, err := json.MarshalIndent(solutions, "", "  ")
	if err != nil {
		logErrors([]hesabu.EvalError{
			{
				Source:     "General",
				Expression: "General",
				Message:    fmt.Sprintf("Could not generate JSON\n%s", err),
			}})
	}
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
	var results map[string]string
	err := json.Unmarshal(raw, &results)
	if err != nil {
		if ShouldLog {
			log.Printf("equations to parse %s", string(raw))
			log.Printf("equations not loaded %v ", err)
		}
		return nil, err
	}
	if ShouldLog {
		log.Printf("equations loaded: %d ", len(results))
	}
	return results, nil
}

func startProfilingCPU(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
}

func stopProfilingCPU() {
	pprof.StopCPUProfile()
}

func startProfilingMemory(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	defer f.Close()
}
