package hesabu

import (
	"encoding/json"
	"github.com/Knetic/govaluate"
	"io/ioutil"
	"testing"
)

type ParserTest struct {
	Name               string
	Input              string
	Expected           string
	ParserErrorMessage string
	Solution           interface{}
	SolutionError      string
}

func TestArrayFunction(t *testing.T) {
	testCases := []ParserTest{
		{
			Name:     "array in sum",
			Input:    "sum(array(1,2,3))",
			Solution: 6.0,
		},
		{
			Name:     "ARRAY in avg",
			Input:    "avg(ARRAY(1,2,3,4,5))",
			Solution: 3.0,
		},
	}
	runEvaluationTests(testCases, t)
}

func TestEvalArrayFunction(t *testing.T) {
	testCases := []ParserTest{
		{
			Name:     "Simple case",
			Input:    "sum(eval_array('a', (1,2,5), 'b', (3,4,5), 'a + b'))",
			Solution: 20.0,
		},
		{
			Name:     "Simple case with use of array function",
			Input:    "sum(eval_array('a', array(1,2,5), 'b', array(3,4,5), 'a + b'))",
			Solution: 20.0,
		},
		{
			Name:     "More complex",
			Input:    "sum(eval_array('quantity_is_null', (0,1,0,1,0,1), 'stock_is_null', (1,1,0,0,0,1), 'if(quantity_is_null + stock_is_null == 2, 1, 0)'))",
			Solution: 2.0,
		},
		{
			Name:          "Malformed formulas return an error",
			Input:         "sum(eval_array('a', array(1), 'b', array(1,2), 'a + b'))",
			SolutionError: "customErrorFunction",
			Solution:      9000,
		},
	}
	runEvaluationTests(testCases, t)
}

func TestCleanerLeavesAlone(t *testing.T) {
	leave_me_alones := []string{"a<=basic",
		"a>=basic",
		"a==basic",
		"a!=basic",
		"a <=basic",
		"a >=basic",
		"a ==basic",
		"a!= basic",
		"a <= basic",
		"a >= basic",
		"a == basic",
		"a != basic",
		"a<= basic",
		"a>= basic",
		"a== basic",
		"a !=basic"}

	var parserTests []ParserTest
	for _, leave_me_alone := range leave_me_alones {
		parserTests = append(parserTests, ParserTest{
			Name:     leave_me_alone,
			Input:    leave_me_alone,
			Expected: leave_me_alone,
		})
	}
	runEvaluationTests(parserTests, t)
}

func TestCleanerSingleEquals(t *testing.T) {
	replace_single_equals := map[string]string{
		"a=b":  "a==b",
		"a =b": "a ==b",
		"a= b": "a== b",
	}
	var parserTests []ParserTest
	for input, output := range replace_single_equals {
		parserTests = append(parserTests, ParserTest{
			Name:     input,
			Input:    input,
			Expected: output,
		})
	}
	runEvaluationTests(parserTests, t)
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
			Name:     "Replace single = with ==",
			Input:    "a=b && b     =     c && d = e",
			Expected: "a==b && b     ==     c && d == e",
			Solution: false,
		},
		{
			Name:     "AND and equals",
			Input:    "(a == 1) AND (b = 2) or (a = 1) || (b == 2)",
			Expected: "(a == 1) && (b == 2) || (a == 1) || (b == 2)",
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
			Name:               "Malformed formulas return an error",
			Input:              "=operator=1",
			ParserErrorMessage: "Invalid token: '='",
		},
		{
			Name:               "Malformed formulas return an error",
			Input:              "operator=1=",
			ParserErrorMessage: "Invalid token: '='",
		},
	}
	runEvaluationTests(parserTests, t)
}

func runEvaluationTests(parserTests []ParserTest, t *testing.T) {
	functions := Functions() //map[string]govaluate.ExpressionFunction{}
	for _, parserTest := range parserTests {
		equations := map[string]string{
			"abool":   "true",
			"bbool":   "false",
			"basic":   "4",
			"a":       "1",
			"b":       "2",
			"testing": parserTest.Input,
		}
		parsedEquations := Parse(equations, functions)

		if parserTest.Expected != "" {
			if len(parsedEquations.Errors) > 0 {
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

		if parserTest.Solution != nil {
			solution, err := parsedEquations.Solve()
			if err != nil {
				if parserTest.SolutionError != "" {
					return
				}
				t.Logf("Has an error while test has not set a Solution Error : %s", err)
				t.Fail()
				continue
			}

			if parserTest.Solution != solution["testing"] {
				t.Logf("Test '%s' '%s' vs '%s'", parserTest.Name, parserTest.Solution, solution["testing"])
				t.Fail()
				continue
			}
		}

		if parserTest.ParserErrorMessage != "" {
			errors := parsedEquations.Errors
			if len(errors) < 1 {
				t.Logf("%s - Expected an error with %s but none was returned", parserTest.Name, parserTest.ParserErrorMessage)
				t.Fail()
				continue
			}
			was := errors[0].Message
			if parserTest.ParserErrorMessage != was {
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
