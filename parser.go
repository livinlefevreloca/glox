package main

import "fmt"

type Parser struct {
	tokens        []Token
	current       int
	errorReporter func(*Token, int, int, string)
}

func reportErrorParse(token *Token, line int, where int, message string) {
	fmt.Printf("[line %d, col %d] Error at %s, %s\n", line, where, token.lexeme, message)
}

func NewParser(tokens []Token, reportError func(*Token, int, int, string)) Parser {
	return Parser{tokens: tokens, current: 0, errorReporter: reportError}
}

func (p *Parser) parse() ([]Statement, error) {
	var statements []Statement
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) declaration() (Statement, error) {
	if p.match(TOKEN_VAR) {
		if stmt, err := p.variableDeclaration(); err != nil {
			return nil, err
		} else {
			return stmt, nil
		}
	}

	return p.statement()
}

func (p *Parser) statement() (Statement, error) {
	if p.match(TOKEN_PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) expressionStatement() (Statement, error) {
	expr, err := p.assignment()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_SEMICOLON, "Expected ';' after expression.")
	if err != nil {
		return nil, err
	}
	return ExpressionStatement{expr}, nil
}

func (p *Parser) printStatement() (Statement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_SEMICOLON, "Expected ';' after expression.")
	return PrintStatement{expr}, nil

}

func (p *Parser) declarationStatement() (Statement, error) {
	if p.match(TOKEN_VAR) {
		return p.variableDeclaration()
	}

	return p.statement()
}

func (p *Parser) variableDeclaration() (Statement, error) {
	name, err := p.consume(TOKEN_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var expr Expr
	if p.match(TOKEN_EQUAL) {
		expr, err = p.assignment()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return VarDeclarationStatement{name: *name, initializer: expr}, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.match(TOKEN_EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(Variable); ok {
			name := varExpr.name
			return Assign{name, value}, nil
		}

		return nil, fmt.Errorf("Invalid assignment target: %s", equals)
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.ternary()
}

func (p *Parser) ternary() (Expr, error) {
	condition, err := p.block()
	if err != nil {
		return nil, err
	}
	if p.match(TOKEN_QUESTION_MARK) {
		left, err := p.ternary()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TOKEN_COLON, "Expected ':' in Ternay expression")
		if err != nil {
			return nil, err
		}
		right, err := p.ternary()
		return Ternary{condition: condition, left: left, right: right}, nil
	}

	return condition, nil
}

func (p *Parser) block() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(TOKEN_COMMA) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = Binary{left: expr, operator: Operator{operator: *operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_BANG_EQUAL, TOKEN_EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{left: expr, operator: Operator{operator: *operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = Binary{left: expr, operator: Operator{operator: *operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_MINUS, TOKEN_PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{left: expr, operator: Operator{operator: *operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_SLASH, TOKEN_STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{left: expr, operator: Operator{operator: *operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(TOKEN_BANG, TOKEN_MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{operator: *operator, right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TOKEN_FALSE) {
		return Literal{value: *p.previous()}, nil
	}
	if p.match(TOKEN_TRUE) {
		return Literal{value: *p.previous()}, nil
	}
	if p.match(TOKEN_NIL) {
		return Literal{value: *p.previous()}, nil
	}

	if p.match(TOKEN_NUMBER, TOKEN_STRING) {
		return Literal{value: *p.previous()}, nil
	}

	if p.match(TOKEN_IDENTIFIER) {
		return Variable{name: *p.previous()}, nil
	}

	if p.match(TOKEN_LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return Grouping{expression: expr}, nil
	}

	return nil, fmt.Errorf("Expected expression. got %s", p.peek().lexeme)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == TOKEN_SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case TOKEN_CLASS, TOKEN_FUN, TOKEN_VAR, TOKEN_FOR, TOKEN_IF, TOKEN_WHILE, TOKEN_PRINT, TOKEN_RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType int, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	tok := p.peek()
	if p.isAtEnd() {
		err := fmt.Errorf("Reached unexpected EOF")
		p.errorReporter(&tok, tok.line, tok.col, err.Error())
		return nil, err
	}

	err := fmt.Errorf(message)
	p.errorReporter(&tok, tok.line, tok.col, err.Error())
	return nil, fmt.Errorf(message)
}

func (p *Parser) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == TOKEN_EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}
