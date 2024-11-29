package main

import (
	"io"
	"os"
	"testing"
)


func TestClassesAndInheritance(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		// Class declaration and instantiation
		{
			`class Breakfast { cook() { print "Eggs a-fryin'!"; } } Breakfast().cook();`,
			"Eggs a-fryin'!\n",
			false,
		},
		// Class initializer
		{
			`class Foo { init() { print "Foo initialized!"; } } Foo();`,
			"Foo initialized!\n",
			false,
		},
		// Inheritance
		{
			`class Animal { speak() { print "The animal makes a sound."; } } class Dog < Animal { speak() { print "The dog barks."; } } Dog().speak();`,
			"The dog barks.\n",
			false,
		},
		// Super call
		{
			`class A { method() { print "A method"; } } class B < A { method() { print "B method"; super.method(); } } B().method();`,
			"B method\nA method\n",
			false,
		},
		// 'this' outside class
		{
			`print this;`,
			"",
			true,
		},
		// 'super' outside class
		{
			`print super.method();`,
			"",
			true,
		},
		// Inheritance from non-class
		{
			`var NotAClass = "not a class"; class Subclass < NotAClass {}`,
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

				originalStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

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