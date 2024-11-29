package main

import (
	"io"
	"os"
	"testing"
	"strings"
)


func TestFunctions(t *testing.T) {
    tests := []struct {
        input       string
        expected    string
        shouldError bool
    }{
        // Simple function declaration and call
        {
            "fun sayHi() { print \"Hi!\"; } sayHi();",
            "Hi!\n",
            false,
        },
        // Function with parameters
        {
            "fun greet(name) { print \"Hello, \" + name + \"!\"; } greet(\"Alice\");",
            "Hello, Alice!\n",
            false,
        },
        // Returning a value
        {
            "fun add(a, b) { return a + b; } print add(3, 4);",
            "7\n",
            false,
        },
        // Nested function calls
        {
            "fun square(x) { return x * x; } fun sumOfSquares(a, b) { return square(a) + square(b); } print sumOfSquares(3, 4);",
            "25\n",
            false,
        },
        // Function without return
        {
            "fun noReturn() { 123; } print noReturn();",
            "nil\n",
            false,
        },
        // Function with no arguments
        {
            "fun doSomething() { return 42; } print doSomething();",
            "42\n",
            false,
        },
        // Recursion
        {
            `
            fun factorial(n) {
              if (n <= 1) return 1;
              return n * factorial(n - 1);
            }
            print factorial(5); // 120
            `,
            "120\n",
            false,
        },
        // Function shadowing
        {
            `
            fun outer() {
              fun inner() {
                return "inner";
              }
              return inner();
            }
            print outer(); // "inner"
            `,
            "inner\n",
            false,
        },
        // Calling a function with the wrong number of arguments
        {
            "fun oneArg(x) { print x; } oneArg();",
            "",
            true,
        },
        {
            "fun oneArg(x) { print x; } oneArg(1, 2);",
            "",
            true,
        },
        
        // Error: Undefined function
        {
            "undefinedFunction();",
            "",
            true,
        },
        // More than 255 arguments
        {
            "fun tooManyArgs(" + strings.Repeat("x,", 255) + "x) { print x; }",
            "",
            true,
        },
        // More than 255 parameters
        {
            "fun tooManyParams(" + strings.Repeat("x,", 255) + "x) { return 42; }",
            "",
            true,
        },
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

