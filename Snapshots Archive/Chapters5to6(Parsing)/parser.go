package main

import (
	"fmt"
	"io"
)

// Parser implements a recursive descent parser for Lox
type Parser struct {
	tokens  []Token
	current int
	stdErr  io.Writer
}

// NewParser creates a new Parser instance
func NewParser(tokens []Token, stdErr io.Writer) *Parser {
	return &Parser{tokens: tokens, stdErr: stdErr}
}

// Parse starts the parsing process
func (p *Parser) Parse() (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			if perr, ok := r.(ParseError); ok {
				err = perr
			} else {
				panic(r)
			}
		}
	}()

	return p.expression(), nil
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenBangEqual, TokenEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(TokenMinus, TokenPlus) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(TokenSlash, TokenStar) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		right := p.unary()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TokenFalse) {
		return &Literal{Value: false}
	}
	if p.match(TokenTrue) {
		return &Literal{Value: true}
	}
	if p.match(TokenNil) {
		return &Literal{Value: nil}
	}

	if p.match(TokenNumber, TokenString) {
		return &Literal{Value: p.previous().Literal}
	}

	if p.match(TokenLeftParen) {
		expr := p.expression()
		p.consume(TokenRightParen, "Expect ')' after expression.")
		return &Grouping{Expression: expr}
	}

	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == TokenEof
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token Token, message string) ParseError {
	errorMessage := fmt.Sprintf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
	_, _ = p.stdErr.Write([]byte(errorMessage + "\n"))
	return ParseError{message: errorMessage}
}

// ParseError represents a parsing error
type ParseError struct {
	message string
}

func (e ParseError) Error() string {
	return e.message
}