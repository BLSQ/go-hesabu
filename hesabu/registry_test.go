package hesabu

import (
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
	}

	for _, table := range tables {
		functionToCall := Functions()[table.functionToCall]

		result, err := functionToCall(table.args...)
		if err != nil {
			t.Errorf("errored")
		}
		if result != table.expected {
			t.Errorf("%s(%v) was incorrect, got: %v, want: %v.", table.functionToCall, table.args, result, table.expected)
		}
	}
}

func TestIfFunctionWithIncorrectBool(t *testing.T) {
	inputData := []interface{}{1, 2, 3}
	_, err := Functions()["IF"](inputData...)
	if err == nil {
		t.Logf("1 is not a bool, this should have failed")
		t.Fail()
	}
}
