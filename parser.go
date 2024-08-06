package main

import "fmt"

type Parser struct {
	tokens        []Token
	current       int
	errorReporter func(*Token, int, int, string)
}

func reportErrorParse(token *Token, line int, where int, message string) {
	fmt.Printf("[line %d, col %d] Error at %s, %s\n", line, where, token.lexeme, message)
	hadError = true
}

func NewParser(tokens []Token, reportError func(*Token, int, int, string)) Parser {
	return Parser{tokens: tokens, current: 0, errorReporter: reportError}
}

func (p *Parser) parse() (Expr[string], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) expression() (Expr[string], error) {
	return p.ternary()
}

func (p *Parser) ternary() (Expr[string], error) {
	condition, err := p.block()
	if err != nil {
		return nil, err
	}
	if p.match(TOKEN_QUESTION_MARK) {
		left, err := p.ternary()
		if err != nil {
			return nil, err
		}
		err = p.consume(TOKEN_COLON, "Expected ':' in Ternay expression")
		if err != nil {
			return nil, err
		}
		right, err := p.ternary()
		return Ternary[string]{condition: condition, left: left, right: right}, nil
	}

	return condition, nil
}

func (p *Parser) block() (Expr[string], error) {
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
		expr = Binary[string]{left: expr, operator: Operator[string]{operator: operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr[string], error) {
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
		expr = Binary[string]{left: expr, operator: Operator[string]{operator: operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr[string], error) {
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

		expr = Binary[string]{left: expr, operator: Operator[string]{operator: operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr[string], error) {
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
		expr = Binary[string]{left: expr, operator: Operator[string]{operator: operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr[string], error) {
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
		expr = Binary[string]{left: expr, operator: Operator[string]{operator: operator}, right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr[string], error) {
	if p.match(TOKEN_BANG, TOKEN_MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary[string]{operator: operator, right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr[string], error) {
	if p.match(TOKEN_FALSE) {
		return Literal[string]{value: p.previous()}, nil
	}
	if p.match(TOKEN_TRUE) {
		return Literal[string]{value: p.previous()}, nil
	}
	if p.match(TOKEN_NIL) {
		return Literal[string]{value: p.previous()}, nil
	}

	if p.match(TOKEN_NUMBER, TOKEN_STRING) {
		return Literal[string]{value: p.previous()}, nil
	}

	if p.match(TOKEN_LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		err = p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return Grouping[string]{expression: expr}, nil
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

func (p *Parser) consume(tokenType int, message string) error {
	if p.check(tokenType) {
		p.advance()
		return nil
	}

	tok := p.peek()
	if p.isAtEnd() {
		err := fmt.Errorf("Reached unexpected EOF")
		p.errorReporter(&tok, tok.line, tok.col, err.Error())
		return err
	}

	err := fmt.Errorf(message)
	p.errorReporter(&tok, tok.line, tok.col, err.Error())
	return fmt.Errorf(message)
}

func (p *Parser) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
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

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
