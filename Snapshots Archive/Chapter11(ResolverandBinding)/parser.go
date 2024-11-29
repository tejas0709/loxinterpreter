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
	return p.assignment()
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

	return p.call()
}

func (p *Parser) primary() Expr {
    if p.match(TokenIdentifier) {
        return &Variable{Name: p.previous()}
    }
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

func (p *Parser) ParseStatements() ([]Stmt, error) {
	defer func() {
		if r := recover(); r != nil {
			if perr, ok := r.(ParseError); ok {
				panic(perr)
			} else {
				panic(r)
			}
		}
	}()

	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, nil
}


func (p *Parser) declaration() Stmt {
	if p.match(TokenFun) {
        return p.function("function")
    }
	if p.match(TokenVar) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(TokenIdentifier, "Expect variable name.")

	var initializer Expr
	if p.match(TokenEqual) {
		initializer = p.expression()
	}
	p.consume(TokenSemicolon, "Expect ';' after variable declaration.")
	return &VarStmt{Name: name, Initializer: initializer}
}

func (p *Parser) statement() Stmt {
    if p.match(TokenPrint) {
        return p.printStatement()
    }
    if p.match(TokenReturn) {
        return p.returnStatement()
    }
    if p.match(TokenLeftBrace) {
        return &BlockStmt{Statements: p.block()}
    }
    if p.match(TokenIf) {
        return p.ifStatement()
    }
    if p.match(TokenWhile) {
        return p.whileStatement()
    }
    if p.match(TokenFor) {
        return p.forStatement()
    }
    if p.match(TokenBreak) {
        return p.breakStatement()
    }
    return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(TokenSemicolon, "Expect ';' after value.")
	return &PrintStmt{Expression: value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(TokenSemicolon, "Expect ';' after expression.")
	return &ExpressionStmt{Expression: expr}
}

// Parse an assignment expression.
func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(TokenEqual) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*Variable); ok {
			return &Assign{Name: variable.Name, Value: value}
		}

		panic(p.error(equals, "Invalid assignment target."))
	}

	return expr
}

func (p *Parser) block() []Stmt {
	var statements []Stmt

	for !p.check(TokenRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(TokenRightBrace, "Expect '}' after block.")
	return statements
}

func (p *Parser) ifStatement() Stmt {
    p.consume(TokenLeftParen, "Expect '(' after 'if'.")
    condition := p.expression()
    p.consume(TokenRightParen, "Expect ')' after if condition.")

    thenBranch := p.statement()
    var elseBranch Stmt
    if p.match(TokenElse) {
        elseBranch = p.statement()
    }

    return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) whileStatement() Stmt {
    p.consume(TokenLeftParen, "Expect '(' after 'while'.")
    condition := p.expression()
    p.consume(TokenRightParen, "Expect ')' after condition.")
    body := p.statement()
    return &WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) forStatement() Stmt {
    p.consume(TokenLeftParen, "Expect '(' after 'for'.")
    var initializer Stmt
    if p.match(TokenSemicolon) {
        initializer = nil
    } else if p.match(TokenVar) {
        initializer = p.varDeclaration()
    } else {
        initializer = p.expressionStatement()
    }

    var condition Expr
    if !p.check(TokenSemicolon) {
        condition = p.expression()
    }
    p.consume(TokenSemicolon, "Expect ';' after loop condition.")

    var increment Expr
    if !p.check(TokenRightParen) {
        increment = p.expression()
    }
    p.consume(TokenRightParen, "Expect ')' after for clauses.")

    body := p.statement()
    if increment != nil {
        body = &BlockStmt{Statements: []Stmt{body, &ExpressionStmt{Expression: increment}}}
    }
    if condition == nil {
        condition = &Literal{Value: true}
    }
    body = &WhileStmt{Condition: condition, Body: body}
    if initializer != nil {
        body = &BlockStmt{Statements: []Stmt{initializer, body}}
    }

    return body
}

func (p *Parser) breakStatement() Stmt {
    p.consume(TokenSemicolon, "Expect ';' after 'break'.")
    return &BreakStmt{}
}

func (p *Parser) function(kind string) Stmt {
    name := p.consume(TokenIdentifier, fmt.Sprintf("Expect %s name.", kind))
    p.consume(TokenLeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
    var parameters []Token
    if !p.check(TokenRightParen) {
        for {
            if len(parameters) >= 255 {
                p.error(p.peek(), "Cannot have more than 255 parameters.")
            }
            parameters = append(parameters, p.consume(TokenIdentifier, "Expect parameter name."))
            if !p.match(TokenComma) {
                break
            }
        }
    }
    p.consume(TokenRightParen, "Expect ')' after parameters.")
    p.consume(TokenLeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
    body := p.block()
    return &FunStmt{Name: name, Params: parameters, Body: body}
}


func (p *Parser) call() Expr {
    expr := p.primary()

    for {
        if p.match(TokenLeftParen) {
            expr = p.finishCall(expr)
        } else {
            break
        }
    }

    return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
    var arguments []Expr
    if !p.check(TokenRightParen) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(TokenComma) {
				break
			}
		}
	}	
    paren := p.consume(TokenRightParen, "Expect ')' after arguments.")
    return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (p *Parser) returnStatement() Stmt {
    keyword := p.previous()
    var value Expr
    if !p.check(TokenSemicolon) {
        value = p.expression()
    }
    p.consume(TokenSemicolon, "Expect ';' after return value.")
    return &ReturnStmt{Keyword: keyword, Value: value}
}