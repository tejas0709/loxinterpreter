package main

import (
	"io"
	"os"
	"testing"
)

func TestStatementsAndState(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		// Print statements
		{"print 123;", "123\n", false},
		{"print \"hello\";", "hello\n", false},
		{"print nil;", "nil\n", false},

		// Variable declarations
		{"var a = 123; print a;", "123\n", false},
		{"var b; print b;", "nil\n", false},
		{"var c = \"test\"; print c;", "test\n", false},

		// Variable reassignments
		{"var x = 10; x = 20; print x;", "20\n", false},
		{"var y; y = 50; print y;", "50\n", false},

		// Expression statements
		{"123;", "", false},
		{"\"test\";", "", false},

		// Undefined variables (should error, comment out to check error test cases)
		{"print z;", "", true},
		{"var x = 10; z = x + 1;", "", true},

		// Nested expressions and scope
		{"var a = 10; var b = a + 20; print b;", "30\n", false},
		{"var outer = 10; { var inner = 20; print inner; } print outer;", "20\n10\n", false},
		{"{ var a = 5; print a; } print a;", "", true}, // "a" is not defined outside the block
	
		{"var x = 10; { var x = 20; print x; } print x;", "20\n10\n", false},      // Nested block scoping
		{"{ var a = 1; var b = 2; print a; print b; }", "1\n2\n", false},         //  Multiple variable declarations in a block
		{"var x = 10; { var x = x + 5; print x; }", "15\n", false},              //  Shadowing variables
		{"var a = 10; { a = a + 5; print a; }", "15\n", false},                 //   Complex nested expressions and reassignments
	
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var output string
			var didError bool

			// Capture output from the interpreter
			func() {
				defer func() {
					if r := recover(); r != nil {
						didError = true
					}
				}()

				// Redirect standard output to capture interpreter output
				originalStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				// Run the interpreter
				scanner := NewScanner(tt.input, nil)
				tokens := scanner.ScanTokens()

				parser := NewParser(tokens, nil)
				statements, err := parser.ParseStatements()
				if err != nil {
					didError = true
					return
				}

				interpreter := NewInterpreter()
				interpreter.InterpretStatements(statements)

				// Capture and restore standard output
				w.Close()
				outBytes, _ := io.ReadAll(r)
				output = string(outBytes)
				os.Stdout = originalStdout
			}()

			if didError != tt.shouldError {
				t.Errorf("Expected error: %v, but got: %v", tt.shouldError, didError)
				return
			}

			if !tt.shouldError && output != tt.expected {
				t.Errorf("Expected output: %q, but got: %q", tt.expected, output)
			}
		})
	}
}

