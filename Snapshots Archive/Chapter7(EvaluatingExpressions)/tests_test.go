package main

import (
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		// Literal values
		{"123", "123", false},
		{"\"hello\"", "hello", false},
		{"nil", "nil", false},

		// Unary expressions
		{"-123", "-123", false},
		{"!true", "false", false},
		{"!nil", "true", false},
		{"!123", "false", false},

		// Binary expressions
		{"1 + 2", "3", false},
		{"5 - 3", "2", false},
		{"2 * 3", "6", false},
		{"8 / 4", "2", false},

		// Operator precedence
		{"1 + 2 * 3", "7", false},
		{"(1 + 2) * 3", "9", false},

		// Comparison operators
		{"5 > 3", "true", false},
		{"3 < 4", "true", false},
		{"5 >= 5", "true", false},
		{"3 <= 3", "true", false},

		// Equality
		{"4 == 4", "true", false},
		{"4 != 5", "true", false},
		{"nil == nil", "true", false},
		{"nil != 0", "true", false},

		// String concatenation
		{"\"a\" + \"b\"", "ab", false},
		{"\"hello\" + \" \" + \"world\"", "hello world", false},
		{"\"a\" + \"\"", "a", false},

		// Invalid cases
		{"1 + \"hello\"", "", true},
		{"true + 1", "", true},
		{"123 / \"string\"", "", true},
		{"-nil", "", true},

		// Division by zero
		{"1 / 0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			scanner := NewScanner(tt.input, nil)
			tokens := scanner.ScanTokens()

			parser := NewParser(tokens, nil)
			expr, err := parser.Parse()
			if err != nil {
				t.Errorf("Parsing failed: %v", err)
				return
			}

			interpreter := NewInterpreter()
			var output string
			var didError bool
			func() {
				defer func() {
					if r := recover(); r != nil {
						didError = true
					}
				}()
				output = stringify(interpreter.evaluate(expr))
			}()

			if tt.shouldError != didError {
				t.Errorf("Expected error state %v, got %v", tt.shouldError, didError)
			}

			if !tt.shouldError && output != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, output)
			}
		})
	}
}

