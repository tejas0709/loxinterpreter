package main

import (
	"bytes"
	"testing"
)


func TestResolver(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		// Basic variable declaration
		{
			`var a = 1; print a;`,
			"1\n",
			false,
		},
		// Block scope variable shadowing
		{
			`var a = 1; { var a = 2; print a; } print a;`,
			"2\n1\n",
			false,
		},
		// Undefined variable access
		{
			`print b;`,
			"",
			true,
		},
		// Function declaration and usage
		{
			`fun test() { print "Hello, World!"; } test();`,
			"Hello, World!\n",
			false,
		},
		// Return from function
		{
			`fun test() { return 42; } print test();`,
			"42\n",
			false,
		},
		// Return outside function
		{
			`return 42;`,
			"",
			true,
		},
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

				// Run the resolver and interpreter
				scanner := NewScanner(tt.input, nil)
				tokens := scanner.ScanTokens()

				parser := NewParser(tokens, nil)
				statements, err := parser.ParseStatements()
				if err != nil {
					didError = true
					return
				}

				interpreter := NewInterpreter()
				resolver := NewResolver(interpreter)

				resolver.Resolve(statements)
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