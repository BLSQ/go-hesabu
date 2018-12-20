package hesabu

import (
	"encoding/json"
	"github.com/Knetic/govaluate"
	"io/ioutil"
	"testing"
)

type ParserTest struct {
	Name                 string
	Input                string
	Expected             string
	ExpectedErrorMessage string
	Solution             interface{}
}

func TestCleaner(t *testing.T) {
	parserTests := []ParserTest{
		{
			Name:     "Sanity check",
			Input:    "a + b",
			Expected: "a + b",
			Solution: 3.0,
		},
		{
			Name:     "Replace AND",
			Input:    "abool AND bbool",
			Expected: "abool && bbool",
			Solution: false,
		},
		{
			Name:     "Replace and",
			Input:    "abool and bbool",
			Expected: "abool && bbool",
			Solution: false,
		},
		{
			Name:     "Replace OR",
			Input:    "abool OR bbool",
			Expected: "abool || bbool",
			Solution: true,
		},
		{
			Name:     "Replace or",
			Input:    "abool or bbool",
			Expected: "abool || bbool",
			Solution: true,
		},
		{
			Name:     "Leaves <= alone",
			Input:    "a <= b",
			Expected: "a <= b",
			Solution: true,
		},
		{
			Name:     "Leaves == alone",
			Input:    "a == b",
			Expected: "a == b",
			Solution: false,
		},
		{
			Name:     "Replace single = with ==",
			Input:    "a=b && b     =     c && d = e",
			Expected: "a==b && b     ==     c && d == e",
			Solution: false,
		},
		{
			Name:     "Leaves alone variable containing AND",
			Input:    "operand=1",
			Expected: "operand==1",
		},
		{
			Name:     "Leaves alone variable containing or",
			Input:    "operator=1",
			Expected: "operator==1",
		},
		{
			Name:                 "Malformed formulas return an error",
			Input:                "=operator=1",
			ExpectedErrorMessage: "Invalid token: '='",
		},
		{
			Name:                 "Malformed formulas return an error",
			Input:                "operator=1=",
			ExpectedErrorMessage: "Invalid token: '='",
		},
	}
	runEvaluationTests(parserTests, t)
}

func runEvaluationTests(parserTests []ParserTest, t *testing.T) {
	functions := map[string]govaluate.ExpressionFunction{}
	for _, parserTest := range parserTests {
		equations := map[string]string{
			"abool":   "true",
			"bbool":   "false",
			"a":       "1",
			"b":       "2",
			"testing": parserTest.Input,
		}
		parsedEquations := Parse(equations, functions)
		if parserTest.Expected != "" {
			solution, err := parsedEquations.Solve()
			if err != nil {
				t.Logf("err not nil : %s", err)
			} else {
				if parserTest.Solution != nil && parserTest.Solution != solution["testing"] {
					t.Logf("Test '%s' '%F' vs '%s'", parserTest.Name, parserTest.Solution, solution["testing"])
					t.Fail()
					continue
				}
			}
			if len(parsedEquations.Errors) < 0 {
				t.Logf("%s - had an error but should not have had", parserTest.Name)
				t.Fail()
				continue
			}
			was := parsedEquations.Equations["testing"].String()
			if parserTest.Expected != was {
				t.Logf("Test '%s' '%s' vs '%s'", parserTest.Name, parserTest.Expected, was)
				t.Fail()
				continue
			}
		}

		if parserTest.ExpectedErrorMessage != "" {
			errors := parsedEquations.Errors
			if len(errors) < 1 {
				t.Logf("%s - Expected an error with %s but none was returned", parserTest.Name, parserTest.ExpectedErrorMessage)
				t.Fail()
				continue
			}
			was := errors[0].Message
			if parserTest.ExpectedErrorMessage != was {
				t.Logf("Test '%s' '%s' vs '%s'", parserTest.Name, parserTest.Expected, was)
				t.Fail()
				continue
			}
		}

	}
}

// Without Clean						: BenchmarkParse-4   	      10	 165449475 ns/op
// With replaceSingleEqual  : BenchmarkParse-4   	      10	 188343765 ns/op
// With regex								: BenchmarkParse-4   	       5	 318208619 ns/op
func BenchmarkParse(b *testing.B) {
	functions := map[string]govaluate.ExpressionFunction{}
	raw, err := ioutil.ReadFile("../test/large_set_of_equations.json")
	if err != nil {
		panic("file not read")
	}
	var equations map[string]string
	err = json.Unmarshal(raw, &equations)
	if err != nil {
		panic("Could not read JSON")
	}
	b.ResetTimer()
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Parse(equations, functions)
	}
}
