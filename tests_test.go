package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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

		// Additional valid expressions
		{"1 + 2 * 3 - 4 / 5", "Binary(Binary(Literal(1), +, Binary(Literal(2), *, Literal(3))), -, Binary(Literal(4), /, Literal(5)))", false},
		{"((1 + 2) * (3 - 4)) / 5", "Binary(Grouping(Binary(Grouping(Binary(Literal(1), +, Literal(2))), *, Grouping(Binary(Literal(3), -, Literal(4))))), /, Literal(5))", false},

		// Additional invalid expressions
		{"(1 + 2", "", true}, // Missing closing parenthesis
		{"+", "", true},      // Lone operator
		{"! + 1", "", true},  // Invalid unary operator usage

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

		{"var x = 10; { var x = 20; print x; } print x;", "20\n10\n", false}, // Nested block scoping
		{"{ var a = 1; var b = 2; print a; print b; }", "1\n2\n", false},     //  Multiple variable declarations in a block
		{"var x = 10; { var x = x + 5; print x; }", "15\n", false},           //  Shadowing variables
		{"var a = 10; { a = a + 5; print a; }", "15\n", false},               //   Complex nested expressions and reassignments

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
		{"while (false) { break; }", "", false},     // No iterations
		{"for (;;) { break; print 1; }", "", false}, // Break immediately

		// Errors
		{"if () print 1;", "", true},    // Missing condition
		{"while () print 1;", "", true}, // Missing condition
		{"break; ", "", true},           //break outside a loop

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
