package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasError bool
	}{
		{
			name:  "single character tokens",
			input: "( ) { } , . - + ; * [ ]",
			expected: []string{
				fmt.Sprintf("%d ( %v", TokenLeftParen, nil),
				fmt.Sprintf("%d ) %v", TokenRightParen, nil),
				fmt.Sprintf("%d { %v", TokenLeftBrace, nil),
				fmt.Sprintf("%d } %v", TokenRightBrace, nil),
				fmt.Sprintf("%d , %v", TokenComma, nil),
				fmt.Sprintf("%d . %v", TokenDot, nil),
				fmt.Sprintf("%d - %v", TokenMinus, nil),
				fmt.Sprintf("%d + %v", TokenPlus, nil),
				fmt.Sprintf("%d ; %v", TokenSemicolon, nil),
				fmt.Sprintf("%d * %v", TokenStar, nil),
				fmt.Sprintf("%d [ %v", TokenLeftBracket, nil),
				fmt.Sprintf("%d ] %v", TokenRightBracket, nil),
			},
			hasError: false,
		},
		{
			name:  "numbers",
			input: "123 45.67",
			expected: []string{
				fmt.Sprintf("%d 123 %v", TokenNumber, float64(123)),
				fmt.Sprintf("%d 45.67 %v", TokenNumber, 45.67),
			},
			hasError: false,
		},
		{
			name:  "identifiers",
			input: "varName abc123 _test",
			expected: []string{
				fmt.Sprintf("%d varName %v", TokenIdentifier, nil),
				fmt.Sprintf("%d abc123 %v", TokenIdentifier, nil),
				fmt.Sprintf("%d _test %v", TokenIdentifier, nil),
			},
			hasError: false,
		},
		{
			name:  "string literals",
			input: "\"hello\" \"world\"",
			expected: []string{
				fmt.Sprintf("%d \"hello\" %v", TokenString, "hello"),
				fmt.Sprintf("%d \"world\" %v", TokenString, "world"),
			},
			hasError: false,
		},
		{
			name:  "operators",
			input: "! != == > >= < <=",
			expected: []string{
				fmt.Sprintf("%d ! %v", TokenBang, nil),
				fmt.Sprintf("%d != %v", TokenBangEqual, nil),
				fmt.Sprintf("%d == %v", TokenEqualEqual, nil),
				fmt.Sprintf("%d > %v", TokenGreater, nil),
				fmt.Sprintf("%d >= %v", TokenGreaterEqual, nil),
				fmt.Sprintf("%d < %v", TokenLess, nil),
				fmt.Sprintf("%d <= %v", TokenLessEqual, nil),
			},
			hasError: false,
		},
		{
			name:  "keywords",
			input: "and class else if nil or true false var while",
			expected: []string{
				fmt.Sprintf("%d and %v", TokenAnd, nil),
				fmt.Sprintf("%d class %v", TokenClass, nil),
				fmt.Sprintf("%d else %v", TokenElse, nil),
				fmt.Sprintf("%d if %v", TokenIf, nil),
				fmt.Sprintf("%d nil %v", TokenNil, nil),
				fmt.Sprintf("%d or %v", TokenOr, nil),
				fmt.Sprintf("%d true %v", TokenTrue, nil),
				fmt.Sprintf("%d false %v", TokenFalse, nil),
				fmt.Sprintf("%d var %v", TokenVar, nil),
				fmt.Sprintf("%d while %v", TokenWhile, nil),
			},
			hasError: false,
		},
		{
			name:  "single line comments",
			input: "// this is a comment\n42",
			expected: []string{
				fmt.Sprintf("%d 42 %v", TokenNumber, float64(42)),
			},
			hasError: false,
		},
		{
			name:  "block comments",
			input: "/* this is\na block\ncomment */42",
			expected: []string{
				fmt.Sprintf("%d 42 %v", TokenNumber, float64(42)),
			},
			hasError: false,
		},
		{
			name:  "nested expressions with comments",
			input: "(123 /* comment */ + /* another */ abc)",
			expected: []string{
				fmt.Sprintf("%d ( %v", TokenLeftParen, nil),
				fmt.Sprintf("%d 123 %v", TokenNumber, float64(123)),
				fmt.Sprintf("%d + %v", TokenPlus, nil),
				fmt.Sprintf("%d abc %v", TokenIdentifier, nil),
				fmt.Sprintf("%d ) %v", TokenRightParen, nil),
			},
			hasError: false,
		},
		{
			name:     "unterminated block comment",
			input:    "/* unterminated",
			expected: []string{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errBuf bytes.Buffer
			sc := NewScanner(tt.input, &errBuf)
			tokens := sc.ScanTokens()

			// Remove EOF token for comparison
			tokens = tokens[:len(tokens)-1]

			var result []string
			for _, token := range tokens {
				result = append(result, token.String())
			}

			if tt.hasError {
				if errBuf.Len() == 0 {
					t.Error("Expected error but got none")
				}
				return
			}

			if errBuf.Len() > 0 {
				t.Errorf("Unexpected error: %s", errBuf.String())
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Token count mismatch.\nGot: %d tokens %v\nExpected: %d tokens %v",
					len(result), result,
					len(tt.expected), tt.expected)
				return
			}

			for i := 0; i < len(result); i++ {
				if tt.expected[i] != result[i] {
					t.Errorf("Token mismatch at position %d:\nExpected: %q\nGot: %q",
						i, tt.expected[i], result[i])
				}
			}
		})
	}
}

// Additional test for EOF token
func TestEOFToken(t *testing.T) {
	sc := NewScanner("", nil)
	tokens := sc.ScanTokens()

	if len(tokens) != 1 {
		t.Errorf("Expected 1 token (EOF), got %d tokens", len(tokens))
		return
	}

	if tokens[0].TokenType != TokenEof {
		t.Errorf("Expected EOF token, got %v", tokens[0])
	}
}