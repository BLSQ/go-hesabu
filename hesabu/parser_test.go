package hesabu

import (
	"encoding/json"
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

func TestVariableAsArray(t *testing.T) {
	functions := Functions() //map[string]govaluate.ExpressionFunction{}
	equations := map[string]string{
		"a":      "array(1,2,-3,4,5)",
		"sum":    "sum(a)",
		"max":    "max(a)",
		"min":    "min(a)",
		"avg":    "avg(a)",
		"access": "access(a,1)",
	}
	parsedEquations := Parse(equations, functions)
	if len(parsedEquations.Errors) > 0 {
		t.Logf("Did not expect any errors while parsing: %v", parsedEquations.Errors)
		t.Fail()
	}

	solution, err := parsedEquations.Solve()
	if err != nil {
		t.Logf("Did not expect an error: %s", err)
		t.Fail()
	}
	if solution["sum"] != (1 + 2 + -3.0 + 4.0 + 5.0) {
		t.Logf("Solution does not match our sum: %f", solution["sum"])
		t.Fail()
	}
	if solution["max"] != 5.0 {
		t.Logf("Solution does not match our max: %f", solution["max"])
		t.Fail()
	}
	if solution["min"] != -3.0 {
		t.Logf("Solution does not match our min: %f", solution["min"])
		t.Fail()
	}
	if solution["avg"] != 1.8 {
		t.Logf("Solution does not match our avg: %f", solution["avg"])
		t.Fail()
	}
	if solution["access"] != 2.0 {
		t.Logf("Solution does not match our access: %f", solution["access"])
		t.Fail()
	}
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
			Name:     "Can handle negative numbers",
			Input:    "sum(eval_array('a', (1,-2,5), 'b', (3,4,-5), 'a - b'))",
			Solution: (1.0 - 3.0 + -2.0 - 4.0 + 5.0 - -5.0),
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
			"c":       "4",
			"d":       "5",
			"e":       "6",
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

func TestUnboundVariables(t *testing.T) {
	functions := Functions()
	equations := map[string]string{
		"result": "a + b + c",
	}
	parsedEquations := Parse(equations, functions)
	_, err := parsedEquations.Solve()
	if _, ok := err.(*CustomError); !ok {
		t.Logf("Expected an eval error because a, b and c were never defined")
		t.Fail()
	}
}

func TestDetectsAndReportsNan(t *testing.T) {
	functions := Functions()
	equations := map[string]string{
		"a":      "0",
		"b":      "0",
		"result": "a/b",
	}
	parsedEquations := Parse(equations, functions)
	_, err := parsedEquations.Solve()
	if _, ok := err.(*CustomError); !ok {
		t.Logf("Expected an NaN error because 0/0 is Not A Number")
		t.Fail()
	}
}

func TestDetectsAndReportsInf(t *testing.T) {
	functions := Functions()
	equations := map[string]string{
		"a":      "5",
		"b":      "0",
		"result": "a/b",
	}
	parsedEquations := Parse(equations, functions)
	_, err := parsedEquations.Solve()
	if _, ok := err.(*CustomError); !ok {
		t.Logf("Expected an Inf error because 5/0 is quite a lot")
		t.Fail()
	}
}

func BenchmarkParse(b *testing.B) {
	functions := Functions()
	raw, err := ioutil.ReadFile("../test/very_large_set_of_equations.json")
	if err != nil {
		panic("file not read")
	}
	var equations map[string]string
	err = json.Unmarshal(raw, &equations)
	if err != nil {
		panic("Could not read JSON")
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(equations, functions)
	}
}

func BenchmarkSolve(b *testing.B) {
	functions := Functions()
	raw, err := ioutil.ReadFile("../test/very_large_set_of_equations.json")
	if err != nil {
		panic("file not read")
	}
	var equations map[string]string
	err = json.Unmarshal(raw, &equations)
	if err != nil {
		panic("Could not read JSON")
	}
	parsedEquations := Parse(equations, functions)
	if len(parsedEquations.Errors) > 0 {
		panic("Error while parsing")
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		parsedEquations.Solve()
	}
}
