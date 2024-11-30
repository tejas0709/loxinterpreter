
# Lox Interpreter in Go

This is a Go implementation of the Lox programming language interpreter, developed by following the first 13 chapters of the book *Crafting Interpreters* by Robert Nystrom. This project demonstrates lexical analysis, parsing, and interpretation of Lox programs.

## Features

- Lexical scanning to tokenize Lox source code.
- Parsing expressions and statements.
- Evaluation of Lox programs using an interpreter.
- Includes basic support for variable resolution and environments.
- Fully written in Go, with test cases to validate functionality.

---

## Setup

### Prerequisites
Ensure you have the following installed on your system:
- [Go](https://golang.org/dl/) (version 1.16 or later)

### Clone the Repository
Clone the repository to your local machine or just extract the submitted zip folder:
```bash
git clone https://github.com/tejas0709/loxinterpreter.git
cd <repository-name>
```

---

## Installation

### Dependencies
This project uses Go modules to manage dependencies. Ensure that your environment is set up for Go modules. Install dependencies using:
```bash
go mod tidy
```

### Build
Compile the project:
```bash
go build -o lox.exe
```
This will create an executable named `lox.exe` in the project directory.

---

## Usage

Run the Lox interpreter by passing a Lox script as an argument:
```bash
./lox.exe <path-to-your-lox-script>
```
And you'll see the code's output on the console

Example(You can just add your test cases in this file):
```bash
./lox.exe print_test.lox
```

## Testing

Unit tests are included to ensure the functionality of the interpreter. Run the tests using:
```bash
go test
```
To manually test Lox scripts, you can run the provided example `print_test.lox` or write your own scripts to validate the interpreter's behavior.

---

## Project Structure

- **`main.go`**: Entry point of the interpreter.
- **`scanner.go`**: Handles lexical analysis (tokenization).
- **`parser.go`**: Parses tokens into an Abstract Syntax Tree (AST).
- **`expr.go`, `stmt.go`**: Definitions for expressions and statements in the AST.
- **`interpreter.go`**: Evaluates the AST to execute Lox programs.
- **`environment.go`**: Manages variable scopes and environments.
- **`resolver.go`**: Resolves variable bindings and handles scope checking.
- **`token.go`**: Contains token definitions and utilities.
- **`tests_test.go`**: Unit tests to validate interpreter components.
- **`print_test.lox`**: Example Lox script for manual testing.

---

## Grammar

The grammar implemented in this interpreter covers the following Lox constructs:

```
program        → declaration* EOF ;

declaration    → varDecl | statement ;

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

statement      → exprStmt | printStmt | block ;

exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
block          → "{" declaration* "}" ;

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
```

---

## References

- [Crafting Interpreters](https://craftinginterpreters.com) by Robert Nystrom
