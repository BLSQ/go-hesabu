package hesabu

import (
	"github.com/Knetic/govaluate"
	"testing"
)

type ParserTest struct {
	Name     string
	Input    string
	Expected string
}

func TestCleaner(t *testing.T) {
	parserTests := []ParserTest{
		{
			Name:     "Sanity check",
			Input:    "a + b",
			Expected: "a + b",
		},
		{
			Name:     "Replace AND",
			Input:    "a AND b",
			Expected: "a && b",
		},
		{
			Name:     "Replace OR",
			Input:    "a OR b",
			Expected: "a || b",
		},
		{
			Name:     "Leaves <= alone",
			Input:    "a <= b",
			Expected: "a <= b",
		},
		{
			Name:     "Leaves == alone",
			Input:    "a == b",
			Expected: "a == b",
		},
		{
			Name:     "Replace single = with ==",
			Input:    "a=b && b     =     c && d = e",
			Expected: "a==b && b    ==    c && d==e",
		},
	}
	runEvaluationTests(parserTests, t)
}

func runEvaluationTests(parserTests []ParserTest, t *testing.T) {
	functions := map[string]govaluate.ExpressionFunction{}
	for _, parserTest := range parserTests {
		equations := map[string]string{"testing": parserTest.Input}
		parsedEquations := Parse(equations, functions)
		was := parsedEquations.Equations["testing"].String()
		if parserTest.Expected != was {
			t.Logf("Test '%s' '%s' vs '%s'", parserTest.Name, parserTest.Expected, was)
			t.Fail()
			continue
		}
	}
}
