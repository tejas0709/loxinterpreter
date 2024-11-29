package main

import "fmt"

// TokenType represents the type of a token, categorized by its role in the language.
type TokenType uint8

// TokenType constants for various tokens in the language.
const (
	// Single-character tokens
	TokenLeftParen TokenType = iota // '('
	TokenRightParen                 // ')'
	TokenLeftBrace                  // '{'
	TokenRightBrace                 // '}'
	TokenLeftBracket                // '['
	TokenRightBracket               // ']'
	TokenComma                      // ','
	TokenDot                        // '.'
	TokenMinus                      // '-'
	TokenPlus                       // '+'
	TokenSemicolon                  // ';'
	TokenSlash                      // '/'
	TokenStar                       // '*'
	TokenColon                      // ':'
	TokenQuestionMark               // '?'
	TokenPipe                       // '|'

	// One or two character tokens
	TokenBang        // '!'
	TokenBangEqual   // '!='
	TokenEqual       // '='
	TokenEqualEqual  // '=='
	TokenGreater     // '>'
	TokenGreaterEqual // '>='
	TokenLess        // '<'
	TokenLessEqual   // '<='

	// Literals
	TokenIdentifier // Identifiers (variable/function names)
	TokenString     // String literals
	TokenNumber     // Numeric literals

	// Keywords
	TokenAnd      // "and"
	TokenClass    // "class"
	TokenElse     // "else"
	TokenFalse    // "false"
	TokenFun      // "fun"
	TokenFor      // "for"
	TokenIf       // "if"
	TokenNil      // "nil"
	TokenOr       // "or"
	TokenPrint    // "print"
	TokenReturn   // "return"
	TokenSuper    // "super"
	TokenThis     // "this"
	TokenTrue     // "true"
	TokenVar      // "var"
	TokenWhile    // "while"
	TokenBreak    // "break"
	TokenContinue // "continue"

	// Special type declaration token
	TokenTypeType // Used for type declarations

	// End-of-file token
	TokenEof
)

// Token represents a single unit of lexical information in the program.
type Token struct {
	TokenType TokenType   // The type of the token.
	Lexeme    string      // The textual representation of the token.
	Literal   interface{} // The literal value (if applicable, e.g., for strings or numbers).
	Line      int         // Line number where the token appears.
	Start     int         // Index from the start of the program.
}

// String provides a readable representation of a token for debugging or logging.
func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.TokenType, t.Lexeme, t.Literal)
}
