package hesabu

import (
	"strings"
	"testing"
)

func TestGeneric(t *testing.T) {
	tables := []struct {
		functionToCall string
		args           []any
		expected       any
	}{
		{"max", []any{4.0, 1.0}, 4.0},
		{"max", []any{1.0, 5.0}, 5.0},
		{"max", []any{-1.0, 5.0}, 5.0},
		{"max", []any{-1.0, -2.0}, -1.0},

		{"min", []any{4.0, 1.0}, 1.0},
		{"min", []any{2.0, 5.0}, 2.0},
		{"min", []any{-1.0, 5.0}, -1.0},
		{"min", []any{-1.0, -2.0}, -2.0},

		{"score_table", []any{1.0, 0.0, 2.0, 50.0, 2.0, 10.0, 95.0}, 50.0},
		{"score_table", []any{3.0, 0.0, 2.0, 50.0, 2.0, 10.0, 95.0}, 95.0},

		{"safe_div", []any{1.0, 0.0}, 0.0},
		{"safe_div", []any{8.0, 2.0}, 4.0},

		{"if", []any{true, 9000, 3}, 9000},
		{"if", []any{false, 2, 9000}, 9000},

		{"avg", []any{1.0, 2.0, 3.0}, 2.0},

		{"sum", []any{1.0, 2.0, 3.0}, 6.0},

		{"stdevp", []any{1.0, 2.0, 13.0, 3.0, 7.0, 9.0}, 4.258977446393546},

		{"round", []any{33.3333333}, 33.0},
		{"round", []any{33.3333333, 2.0}, 33.33},

		{"floor", []any{33.3333333}, 33.0},
		{"floor", []any{-33.3333333}, -34.0},
		{"floor", []any{33.3333333, 10.0}, 30.0},
		{"floor", []any{-33.3333333, 10.0}, -40.0},

		{"ceiling", []any{33.3333333}, 34.0},
		{"ceiling", []any{-33.3333333}, -33.0},
		{"ceiling", []any{33.3333333, 10.0}, 40.0},
		{"ceiling", []any{-33.3333333, 10.0}, -30.0},

		{"trunc", []any{1.2345678}, 1.0},
		{"trunc", []any{-1.2345678}, -1.0},
		{"trunc", []any{1.2345678, 2.0}, 1.23},
		{"trunc", []any{1.2345678, 3.0}, 1.234},
		{"trunc", []any{1.2345678, 4.0}, 1.2345},
		{"trunc", []any{1.2345678, 5.0}, 1.23456},
		{"round", []any{1.2345678, 5.0}, 1.23457},

		{"abs", []any{1.0}, 1.0},
		{"abs", []any{-1.0}, 1.0},

		{"sqrt", []any{4.0}, 2.0},

		{"access", []any{1.0, 2.0, 3.0, 1.0}, 2.0},
		{"access", []any{1.0, 2.0, 3.0, 2.0}, 3.0},

		{"strlen", []any{"1234567"}, 7.0},

		{"cal_days_in_month", []any{2020.0, 2.0}, 29.0},
		{"cal_days_in_month", []any{2020, 2}, 29.0},
		{"cal_days_in_month", []any{2020, 12}, 31.0},
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
	inputData := []any{true}
	_, err := Functions()["sqrt"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestSqrtFunctionWithIncorrectNegative(t *testing.T) {
	inputData := []any{-1.0}
	_, err := Functions()["sqrt"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestIfFunctionWithIncorrectBool(t *testing.T) {
	inputData := []any{1, 2, 3}
	_, err := Functions()["IF"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestSumFunctionWithIncorrectValues(t *testing.T) {
	inputData := []any{"a", "b", "c"}
	_, err := Functions()["SUM"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestRandBetweenFunction(t *testing.T) {
	inputData := []any{1.0, 10.0}
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
	inputData := []any{1.0, 2.0, 8.0}
	_, err := Functions()["access"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestCalDaysInMonthInvalidYear(t *testing.T) {
	inputData := []any{1.0, 2.0}
	_, err := Functions()["cal_days_in_month"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestCalDaysInMonthInvalidMonthLower(t *testing.T) {
	inputData := []any{2020.0, 0}
	_, err := Functions()["cal_days_in_month"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestCalDaysInMonthInvalidMonthGreater(t *testing.T) {
	inputData := []any{2020.0, 13}
	_, err := Functions()["cal_days_in_month"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}

func TestCalDaysInMonthInvalidMontType(t *testing.T) {
	inputData := []any{2020.0, "a"}
	_, err := Functions()["cal_days_in_month"](inputData...)
	if err, ok := err.(*customFunctionError); !ok {
		t.Logf("else, %v", err)
		t.Fail()
	}
}
