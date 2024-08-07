package main

import (
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// Scanner helpers
func getKeywordMap() map[string]int {
	return map[string]int{
		"and":    TOKEN_AND,
		"class":  TOKEN_CLASS,
		"else":   TOKEN_ELSE,
		"false":  TOKEN_FALSE,
		"for":    TOKEN_FOR,
		"fun":    TOKEN_FUN,
		"if":     TOKEN_IF,
		"nil":    TOKEN_NIL,
		"or":     TOKEN_OR,
		"print":  TOKEN_PRINT,
		"return": TOKEN_RETURN,
		"super":  TOKEN_SUPER,
		"this":   TOKEN_THIS,
		"true":   TOKEN_TRUE,
		"var":    TOKEN_VAR,
		"while":  TOKEN_WHILE,
	}
}

func reportErrorScan(line int, where int, message string) {
	fmt.Printf("[line %d, col %d] Error: %s\n", line, where, message)
	hadError = true
}

// Scanner
type GloxScanner struct {
	source        []rune
	start         int
	current       int
	tokens        []Token
	line          int
	lineStart     int
	keywords      map[string]int
	errorReporter func(line int, col int, message string)
}

func NewGloxScanner(source string, errorReporter func(line int, col int, message string)) GloxScanner {
	return GloxScanner{
		source:        []rune(source),
		start:         0,
		current:       0,
		tokens:        make([]Token, 0),
		line:          1,
		lineStart:     0,
		keywords:      getKeywordMap(),
		errorReporter: errorReporter,
	}
}

func (s *GloxScanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(TOKEN_EOF, io.EOF.Error(), nil, s.line, s.start-s.lineStart))
	return s.tokens
}

func (s *GloxScanner) scanToken() {
	s.start = s.current
	c := s.advance()
	switch c {
	case 0:
		s.errorReporter(s.line, s.current-s.lineStart, "Unexpected end of file.")
	case '\n':
		s.line++
		s.lineStart = s.current
		break
	case '(':
		s.addToken(TOKEN_LEFT_PAREN)
		break
	case ')':
		s.addToken(TOKEN_RIGHT_PAREN)
		break
	case '{':
		s.addToken(TOKEN_LEFT_BRACE)
		break
	case '}':
		s.addToken(TOKEN_RIGHT_BRACE)
		break
	case ',':
		s.addToken(TOKEN_COMMA)
		break
	case '.':
		s.addToken(TOKEN_DOT)
		break
	case '-':
		s.addToken(TOKEN_MINUS)
		break
	case '+':
		s.addToken(TOKEN_PLUS)
		break
	case ';':
		s.addToken(TOKEN_SEMICOLON)
		break
	case '*':
		s.addToken(TOKEN_STAR)
		break
	case '?':
		s.addToken(TOKEN_QUESTION_MARK)
		break
	case ':':
		s.addToken(TOKEN_COLON)
		break
	case '!':
		if s.match('=') {
			s.addToken(TOKEN_BANG_EQUAL)
		} else {
			s.addToken(TOKEN_BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(TOKEN_EQUAL_EQUAL)
		} else {
			s.addToken(TOKEN_EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(TOKEN_LESS_EQUAL)
		} else {
			s.addToken(TOKEN_LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(TOKEN_GREATER_EQUAL)
		} else {
			s.addToken(TOKEN_GREATER)
		}
		break
	case '/':
		if s.match('/') {
			s.singleLineComment()
		} else if s.match('*') {
			s.multiLineComment()
		} else {
			s.addToken(TOKEN_SLASH)
		}
		break
	case ' ', '\r', '\t':
		break
	case '"':
		s.string()
	default:
		if unicode.IsDigit(c) {
			s.number()
		} else if unicode.IsLetter(c) {
			for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
				s.advance()
			}
			s.identifier()
		} else {
			s.errorReporter(s.line, s.current-s.lineStart, fmt.Sprintf("Unexpected character: %s.", string(s.peekPrev())))
		}
	}
}

func (s *GloxScanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *GloxScanner) peekEnd() bool {
	return s.current+1 >= len(s.source)
}

func (s *GloxScanner) advance() rune {
	s.current++
	if s.current > len(s.source) {
		return 0
	}
	return s.peekPrev()
}

func (s *GloxScanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != r {
		return false
	}
	s.advance()
	return true
}

func (s *GloxScanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *GloxScanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *GloxScanner) peekPrev() rune {
	if s.current-1 < 0 {
		return 0
	}
	return s.source[s.current-1]
}

func (s *GloxScanner) singleLineComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *GloxScanner) multiLineComment() {
	for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
			s.lineStart = s.current
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.errorReporter(s.line, s.start-s.lineStart, "Unterminated multi-line comment.")
		return
	}
	// Consume the closing '*/'
	s.advance()
	s.advance()
}

func (s *GloxScanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' && !s.peekEnd() {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.errorReporter(s.line, s.start-s.lineStart, "Unterminated string.")
		return
	}
	s.advance()
	s.addTokenLiteral(TOKEN_STRING, string(s.source[s.start+1:s.current-1]))
}

func (s *GloxScanner) number() {
	for unicode.IsDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
	}
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	num, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	// If we can't parse the number, panic. This should never happen.
	if err != nil {
		panic(err)
	}
	s.addTokenLiteral(TOKEN_NUMBER, num)
}

func (s *GloxScanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) || s.peek() == rune('_') {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	if tokenType, ok := s.keywords[text]; ok {
		if tokenType == TOKEN_TRUE {
			s.addTokenLiteral(tokenType, true)
		} else if tokenType == TOKEN_FALSE {
			s.addTokenLiteral(tokenType, false)
		} else if tokenType == TOKEN_NIL {
			s.addTokenLiteral(tokenType, nil)
		} else {
			s.addToken(tokenType)
		}
	} else {
		s.addToken(TOKEN_IDENTIFIER)
	}
}

func (s *GloxScanner) addToken(tokenType int) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *GloxScanner) addTokenLiteral(tokenType int, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, string(text), literal, s.line, s.start-s.lineStart))
}
