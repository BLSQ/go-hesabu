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