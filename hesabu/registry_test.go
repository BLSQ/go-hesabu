package hesabu

import (
	"strings"
	"testing"
)

func TestGeneric(t *testing.T) {
	tables := []struct {
		functionToCall string
		args           []interface{}
		expected       interface{}
	}{
		{"max", []interface{}{4.0, 1.0}, 4.0},
		{"max", []interface{}{1.0, 5.0}, 5.0},
		{"max", []interface{}{-1.0, 5.0}, 5.0},
		{"max", []interface{}{-1.0, -2.0}, -1.0},

		{"min", []interface{}{4.0, 1.0}, 1.0},
		{"min", []interface{}{2.0, 5.0}, 2.0},
		{"min", []interface{}{-1.0, 5.0}, -1.0},
		{"min", []interface{}{-1.0, -2.0}, -2.0},

		{"score_table", []interface{}{1.0, 0.0, 2.0, 50.0, 2.0, 10.0, 95.0}, 50.0},
		{"score_table", []interface{}{3.0, 0.0, 2.0, 50.0, 2.0, 10.0, 95.0}, 95.0},

		{"safe_div", []interface{}{1.0, 0.0}, 0.0},
		{"safe_div", []interface{}{8.0, 2.0}, 4.0},

		{"if", []interface{}{true, 9000, 3}, 9000},
		{"if", []interface{}{false, 2, 9000}, 9000},

		{"avg", []interface{}{1.0, 2.0, 3.0}, 2.0},

		{"sum", []interface{}{1.0, 2.0, 3.0}, 6.0},

		{"stdevp", []interface{}{1.0, 2.0, 13.0, 3.0, 7.0, 9.0}, 4.258977446393546},

		{"round", []interface{}{33.3333333}, 33.0},
		{"round", []interface{}{33.3333333, 2.0}, 33.33},

		{"floor", []interface{}{33.3333333}, 33.0},
		{"floor", []interface{}{-33.3333333}, -34.0},
		{"floor", []interface{}{33.3333333, 10.0}, 30.0},
		{"floor", []interface{}{-33.3333333, 10.0}, -40.0},

		{"ceiling", []interface{}{33.3333333}, 34.0},
		{"ceiling", []interface{}{-33.3333333}, -33.0},
		{"ceiling", []interface{}{33.3333333, 10.0}, 40.0},
		{"ceiling", []interface{}{-33.3333333, 10.0}, -30.0},

		{"trunc", []interface{}{1.2345678}, 1.0},
		{"trunc", []interface{}{-1.2345678}, -1.0},
		{"trunc", []interface{}{1.2345678, 2.0}, 1.23},
		{"trunc", []interface{}{1.2345678, 3.0}, 1.234},
		{"trunc", []interface{}{1.2345678, 4.0}, 1.2345},
		{"trunc", []interface{}{1.2345678, 5.0}, 1.23456},
		{"round", []interface{}{1.2345678, 5.0}, 1.23457},

		{"abs", []interface{}{1.0}, 1.0},
		{"abs", []interface{}{-1.0}, 1.0},

		{"sqrt", []interface{}{4.0}, 2.0},

		{"access", []interface{}{1.0, 2.0, 3.0, 1.0}, 2.0},
		{"access", []interface{}{1.0, 2.0, 3.0, 2.0}, 3.0},

		{"strlen", []interface{}{"1234567"}, 7.0},
	}

	for _, table := range tables {
		variants := []string{table.functionToCall, strings.ToUpper(table.functionToCall)}
		for _, variant := range variants {
			functionToCall, ok := Functions()[variant]
			if !ok {
				t.Errorf("Function %v was not found in functions table", variant)
				t.Fail()
				continue
			}
			result, err := functionToCall(table.args...)
			if err != nil {
				t.Errorf("errored")
			}
			if result != table.expected {
				t.Errorf("%s(%v) was incorrect, got: %v, want: %v.", variant, table.args, result, table.expected)
			}
		}
	}
}

func TestSqrtFunctionWithIncorrectBool(t *testing.T) {
	inputData := []interface{}{true}
	_, err := Functions()["sqrt"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestSqrtFunctionWithIncorrectNegative(t *testing.T) {
	inputData := []interface{}{-1.0}
	_, err := Functions()["sqrt"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestIfFunctionWithIncorrectBool(t *testing.T) {
	inputData := []interface{}{1, 2, 3}
	_, err := Functions()["IF"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestSumFunctionWithIncorrectValues(t *testing.T) {
	inputData := []interface{}{"a", "b", "c"}
	_, err := Functions()["SUM"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestRandBetweenFunction(t *testing.T) {
	inputData := []interface{}{1.0, 10.0}
	value, err := Functions()["randbetween"](inputData...)
	if err != nil {
		t.Logf("randbetween shouldn't return error")
		t.Fail()
	}
	fvalue := value.(float64)
	if fvalue < 1.0 || fvalue > 10.0 {
		t.Logf("randbetween should generate within range specified")
		t.Fail()
	}
}

func TestBothUpperCaseAndLowerCaseVariantsAreFound(t *testing.T) {
	for name := range Functions() {
		if strings.ToLower(name) == name {
			upper := strings.ToUpper(name)
			if _, ok := Functions()[upper]; !ok {
				t.Logf("%v found but no %v", name, upper)
				t.Fail()
			}
		}
	}
}

func TestAccessOutOfRangeError(t *testing.T) {
	inputData := []interface{}{1.0, 2.0, 8.0}
	_, err := Functions()["access"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}
