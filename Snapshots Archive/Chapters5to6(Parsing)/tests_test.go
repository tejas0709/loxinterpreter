package main

import (
	"fmt"
	"testing"
	"bytes"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		// Simple literals
		{"123", "Literal(123)", false},
		{"\"hello\"", "Literal(hello)", false},

		// Unary expressions
		{"-123", "Unary(-, Literal(123))", false},
		{"!true", "Unary(!, Literal(true))", false},

		// Binary expressions
		{"1 + 2", "Binary(Literal(1), +, Literal(2))", false},
		{"3 * (4 - 5)", "Binary(Literal(3), *, Grouping(Binary(Literal(4), -, Literal(5))))", false},

		// Comparison operators
		{"4 > 3", "Binary(Literal(4), >, Literal(3))", false},
		{"5 <= 6", "Binary(Literal(5), <=, Literal(6))", false},

		// Equality operators
		{"7 == 7", "Binary(Literal(7), ==, Literal(7))", false},
		{"8 != 9", "Binary(Literal(8), !=, Literal(9))", false},

		// Nested expressions
		{"(1 + 2) * 3", "Binary(Grouping(Binary(Literal(1), +, Literal(2))), *, Literal(3))", false},

		// Invalid expressions
		{"(1 + )", "", true},  // Missing operand
		{"5 + * 2", "", true}, // Invalid operator usage
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Setup scanner and parser
			scanner := NewScanner(tt.input, nil)
			tokens := scanner.ScanTokens()

			parser := NewParser(tokens, &bytes.Buffer{})
			expr, err := parser.Parse()

			// Check for errors
			if (err != nil) != tt.hasError {
				t.Errorf("Unexpected error state. Got error: %v, Expected error: %v", err, tt.hasError)
				return
			}

			if !tt.hasError {
				// Convert expression to string for comparison
				result := stringifyExpr(expr)
				if result != tt.expected {
					t.Errorf("Parsed expression does not match. Got: %s, Expected: %s", result, tt.expected)
				}
			}
		})
	}
}

// Helper function to stringify an expression for comparison
func stringifyExpr(expr Expr) string {
	switch e := expr.(type) {
	case *Binary:
		return fmt.Sprintf("Binary(%s, %s, %s)", stringifyExpr(e.Left), e.Operator.Lexeme, stringifyExpr(e.Right))
	case *Grouping:
		return fmt.Sprintf("Grouping(%s)", stringifyExpr(e.Expression))
	case *Literal:
		return fmt.Sprintf("Literal(%v)", e.Value)
	case *Unary:
		return fmt.Sprintf("Unary(%s, %s)", e.Operator.Lexeme, stringifyExpr(e.Right))
	default:
		return "Unknown"
	}
}