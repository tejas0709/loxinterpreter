package main

import (
	"io"
	"os"
	"testing"
)


func TestControlFlow(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		// If-else statements
		{"if (true) print 1; else print 2;", "1\n", false},
		{"if (false) print 1; else print 2;", "2\n", false},
		{"if (false) print 1;", "", false}, // No `else` branch

		// Nested if-else
		{"if (true) if (false) print 1; else print 2;", "2\n", false},

		// While loops
		{"var i = 0; while (i < 3) { print i; i = i + 1; }", "0\n1\n2\n", false},
		{"var x = 5; while (x > 0) { print x; x = x - 1; }", "5\n4\n3\n2\n1\n", false},
		{"while (false) print 1;", "", false}, // No loop iteration

		// For loops
		{"for (var i = 0; i < 3; i = i + 1) print i;", "0\n1\n2\n", false},
		{"for (var x = 10; x > 5; x = x - 1) print x;", "10\n9\n8\n7\n6\n", false},
		{"for (;;) { print \"infinite\"; break; }", "infinite\n", false}, // Infinite loop with `break`

		// Complex nesting
		{
			"for (var i = 1; i <= 2; i = i + 1) { " +
				"  for (var j = 1; j <= 2; j = j + 1) { " +
				"    print i * j; " +
				"  } " +
				"}",
			"1\n2\n2\n4\n", false,
		},

		{"var i = 0; while (true) { if (i == 3) break; print i; i = i + 1; }", "0\n1\n2\n", false},
		{"for (var i = 0; i < 5; i = i + 1) { if (i == 2) break; print i; }", "0\n1\n", false},
		{"while (false) { break; }", "", false}, // No iterations
		{"for (;;) { break; print 1; }", "", false}, // Break immediately

		// Errors
		{"if () print 1;", "", true},          // Missing condition
		{"while () print 1;", "", true},       // Missing condition
		{"break; ", "", true},                 //break outside a loop

		
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var output string
			var didError bool

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
