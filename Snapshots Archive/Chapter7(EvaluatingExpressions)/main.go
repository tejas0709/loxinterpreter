package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	run(string(bytes))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		run(line)
	}
}

func run(source string) {
	scanner := NewScanner(source, os.Stderr)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens, os.Stderr)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	interpreter := NewInterpreter()
	interpreter.Interpret(expr)
}