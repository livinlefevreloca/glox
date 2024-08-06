package main

import "fmt"

type Token struct {
	tokenType int
	lexeme    string
	literal   any
	line      int
	col       int
}

func NewToken(tokenType int, lexeme string, literal any, line int, col int) Token {
	return Token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line, col: col}
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, \"%s\", %v)", t.tokenTypeString(), t.lexeme, t.literal)
}

func (t *Token) tokenTypeString() string {
	switch t.tokenType {
	case TOKEN_LEFT_PAREN:
		return "TOKEN_LEFT_PAREN"
	case TOKEN_RIGHT_PAREN:
		return "TOKEN_RIGHT_PAREN"
	case TOKEN_LEFT_BRACE:
		return "TOKEN_LEFT_BRACE"
	case TOKEN_RIGHT_BRACE:
		return "TOKEN_RIGHT_BRACE"
	case TOKEN_COMMA:
		return "TOKEN_COMMA"
	case TOKEN_DOT:
		return "TOKEN_DOT"
	case TOKEN_MINUS:
		return "TOKEN_MINUS"
	case TOKEN_PLUS:
		return "TOKEN_PLUS"
	case TOKEN_SEMICOLON:
		return "TOKEN_SEMICOLON"
	case TOKEN_SLASH:
		return "TOKEN_SLASH"
	case TOKEN_STAR:
		return "TOKEN_STAR"
	case TOKEN_QUESTION_MARK:
		return "TOKEN_QUESTION_MARK"
	case TOKEN_COLON:
		return "TOKEN_COLON"
	case TOKEN_BANG:
		return "TOKEN_BANG"
	case TOKEN_BANG_EQUAL:
		return "TOKEN_BANG_EQUAL"
	case TOKEN_EQUAL:
		return "TOKEN_EQUAL"
	case TOKEN_EQUAL_EQUAL:
		return "TOKEN_EQUAL_EQUAL"
	case TOKEN_GREATER:
		return "TOKEN_GREATER"
	case TOKEN_GREATER_EQUAL:
		return "TOKEN_GREATER_EQUAL"
	case TOKEN_LESS:
		return "TOKEN_LESS"
	case TOKEN_LESS_EQUAL:
		return "TOKEN_LESS_EQUAL"
	case TOKEN_IDENTIFIER:
		return "TOKEN_IDENTIFIER"
	case TOKEN_STRING:
		return "TOKEN_STRING"
	case TOKEN_NUMBER:
		return "TOKEN_NUMBER"
	case TOKEN_AND:
		return "TOKEN_AND"
	case TOKEN_CLASS:
		return "TOKEN_CLASS"
	case TOKEN_ELSE:
		return "TOKEN_ELSE"
	case TOKEN_FALSE:
		return "TOKEN_FALSE"
	case TOKEN_FUN:
		return "TOKEN_FUN"
	case TOKEN_FOR:
		return "TOKEN_FOR"
	case TOKEN_IF:
		return "TOKEN_IF"
	case TOKEN_NIL:
		return "TOKEN_NIL"
	case TOKEN_OR:
		return "TOKEN_OR"
	case TOKEN_PRINT:
		return "TOKEN_PRINT"
	case TOKEN_RETURN:
		return "TOKEN_RETURN"
	case TOKEN_SUPER:
		return "TOKEN_SUPER"
	case TOKEN_THIS:
		return "TOKEN_THIS"
	case TOKEN_TRUE:
		return "TOKEN_TRUE"
	case TOKEN_VAR:
		return "TOKEN_VAR"
	case TOKEN_WHILE:
		return "TOKEN_WHILE"
	case TOKEN_EOF:
		return "TOKEN_EOF"
	default:
		return "Unknown token type"
	}
}

// define token types
const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR
	TOKEN_QUESTION_MARK
	TOKEN_COLON

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FUN
	TOKEN_FOR
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_EOF
)
