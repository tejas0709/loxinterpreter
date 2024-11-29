package main

import (
	"fmt"
	"io"
	"strconv"
)

// Scanner convert a source text
// into a slice of Token-s
type Scanner struct {
	start   int
	current int
	line    int
	source  string
	tokens  []Token
	stdErr  io.Writer
}

// NewScanner returns a new Scanner
func NewScanner(source string, stdErr io.Writer) *Scanner {
	return &Scanner{source: source, stdErr: stdErr}
}

// ScanTokens returns a slice of tokens representing the source text
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// we're at the beginning of the next lexeme
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{TokenType: TokenEof, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(TokenLeftParen)
	case ')':
		s.addToken(TokenRightParen)
	case '{':
		s.addToken(TokenLeftBrace)
	case '}':
		s.addToken(TokenRightBrace)
	case '[':
		s.addToken(TokenLeftBracket)
	case ']':
		s.addToken(TokenRightBracket)
	case ',':
		s.addToken(TokenComma)
	case '.':
		s.addToken(TokenDot)
	case '-':
		s.addToken(TokenMinus)
	case '+':
		s.addToken(TokenPlus)
	case ';':
		s.addToken(TokenSemicolon)
	case ':':
		s.addToken(TokenColon)
	case '*':
		s.addToken(TokenStar)
	case '?':
		s.addToken(TokenQuestionMark)
	case '|':
		s.addToken(TokenPipe)

	// with look-ahead
	case '!':
		var nextToken TokenType
		if s.match('=') {
			nextToken = TokenBangEqual
		} else {
			nextToken = TokenBang
		}
		s.addToken(nextToken)
	case '=':
		var nextToken TokenType
		if s.match('=') {
			nextToken = TokenEqualEqual
		} else {
			nextToken = TokenEqual
		}
		s.addToken(nextToken)
	case '<':
		var nextToken TokenType
		if s.match('=') {
			nextToken = TokenLessEqual
		} else {
			nextToken = TokenLess
		}
		s.addToken(nextToken)
	case '>':
		var nextToken TokenType
		if s.match('=') {
			nextToken = TokenGreaterEqual
		} else {
			nextToken = TokenGreater
		}
		s.addToken(nextToken)
	case '/':
		if s.match('/') {
			// Single-line comment
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			// Block comment
			for !s.isAtEnd() {
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance() // consume *
					s.advance() // consume /
					break
				} else if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}
			if s.isAtEnd() {
				s.error("Unterminated block comment.")
			}
		} else {
			s.addToken(TokenSlash)
		}

	// whitespace
	case ' ':
	case '\r':
	case '\t':

	case '\n':
		s.line++

	// string
	case '"':
		s.string()

	default:
		if s.isDigit(char) {
			s.number()
		} else if s.isAlpha(char) {
			s.identifier()
		} else {
			s.error("Unexpected character.")
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	curr := rune(s.source[s.current])
	s.current++
	return curr
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
		Start:     s.start}
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error("Unterminated string.")
		return
	}

	s.advance() // the closing "

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TokenString, value)
}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	val, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenWithLiteral(TokenNumber, val)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char == '_')
}

var keywords = map[string]TokenType{
	"and":      TokenAnd,
	"class":    TokenClass,
	"else":     TokenElse,
	"false":    TokenFalse,
	"for":      TokenFor,
	"fun":      TokenFun,
	"if":       TokenIf,
	"nil":      TokenNil,
	"or":       TokenOr,
	"print":    TokenPrint,
	"return":   TokenReturn,
	"super":    TokenSuper,
	"this":     TokenThis,
	"true":     TokenTrue,
	"var":      TokenVar,
	"while":    TokenWhile,
	"break":    TokenBreak,
	"continue": TokenContinue,
	"type":     TokenTypeType,
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = TokenIdentifier
	}
	s.addToken(tokenType)
}

func (s *Scanner) isAlphaNumeric(char rune) bool {
	return s.isAlpha(char) || s.isDigit(char)
}

func (s *Scanner) error(message string) {
	_, _ = s.stdErr.Write([]byte(fmt.Sprintf("[line %d] Error: %s\n", s.line, message)))
}
