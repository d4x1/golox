package main

import (
	"fmt"
	"strconv"
)

const (
	// Single-character tokens.
	LEFT_PAREN  = iota
	RIGHT_PAREN // 1
	LEFT_BRACE  // 2
	RIGHT_BRACE // 3
	COMMA       // 4
	DOT         // 5
	MINUS       // 6
	PLUS        // 7
	SEMICOLON   // 8
	SLASH       // 9
	STAR        // 10

	// One or two character tokens.
	BANG          // 11
	BANG_EQUAL    // 12
	EQUAL         // 13
	EQUAL_EQUAL   // 14
	GREATER       // 15
	GREATER_EQUAL // 16
	LESS          // 17
	LESS_EQUAL    // 18

	// Literals.
	IDENTIFIER // 19
	STRING     // 20
	NUMBER     // 21

	// Keywords.
	AND    // 22
	CLASS  // 23
	ELSE   // 24
	FALSE  // 25
	FUN    // 26
	FOR    // 27
	IF     // 28
	NIL    // 29
	OR     // 30
	PRINT  // 31
	RETURN // 32
	SUPER  // 33
	THIS   // 34
	TRUE   // 35
	VAR    // 36
	WHILE  // 37

	EOF // 38
)

func typeToString(a uint) string {
	keywordMap := map[uint]string{
		AND:    "and",
		CLASS:  "class",
		ELSE:   "else",
		FALSE:  "false",
		FOR:    "for",
		FUN:    "fun",
		IF:     "if",
		NIL:    "nil",
		OR:     "or",
		PRINT:  "print",
		RETURN: "return",
		SUPER:  "super",
		THIS:   "this",
		TRUE:   "true",
		VAR:    "var",
		WHILE:  "while",
	}
	if v, ok := keywordMap[a]; ok {
		return fmt.Sprintf("[KEYWORD] %s", v)
	}

	singleCharMap := map[uint]string{
		LEFT_PAREN:  "(",
		RIGHT_PAREN: ")",
		LEFT_BRACE:  "{",
		RIGHT_BRACE: "}",
		COMMA:       ",",
		DOT:         ".",
		MINUS:       "-",
		PLUS:        "+",
		SEMICOLON:   ";",
		SLASH:       "/",
		STAR:        "*",
	}
	if v, ok := singleCharMap[a]; ok {
		return fmt.Sprintf("[SINGLE CHAR] %s", v)
	}

	identifierMap := map[uint]string{
		IDENTIFIER: "IDENTIFIER",
		STRING:     "STRING",
		NUMBER:     "NUMBER",
	}
	if v, ok := identifierMap[a]; ok {
		return fmt.Sprintf("[%s]", v)
	}

	oneOrTwoCharMap := map[uint]string{
		BANG:          "!",
		BANG_EQUAL:    "!=",
		EQUAL:         "=",
		EQUAL_EQUAL:   "==",
		GREATER:       ">",
		GREATER_EQUAL: ">=",
		LESS:          "<",
		LESS_EQUAL:    "<=",
	}
	if v, ok := oneOrTwoCharMap[a]; ok {
		return fmt.Sprintf("[ONE OR TWO CHAR] %s", v)
	}
	return "[EOF]"
}

type token struct {
	Type    uint
	Lexeme  string
	literal interface{}
	line    int
}

func (token token) String() string {
	return fmt.Sprintf("No.: %d, type: %v, lexeme: %v, literal: %v", token.Type, typeToString(token.Type), token.Lexeme, token.literal)
}

func newToken(typ uint, lexeme string, literal interface{}, line int) token {
	return token{
		Type:    typ,
		Lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func newScanner(source string) *scanner {
	return &scanner{
		source: source,
		line:   1,
	}
}

type scanner struct {
	source string
	tokens []token

	start   int
	current int
	line    int
}

func (s *scanner) scanTokens() ([]token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, newToken(EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() uint8 {
	s.current++
	return s.source[s.current-1]
}
func (s *scanner) addToken(typ uint, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, newToken(typ, text, literal, s.line))
}

func (s *scanner) match(ch uint8) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != ch {
		return false
	}
	s.current++
	return true
}
func (s *scanner) peek() uint8 {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *scanner) peekNext() uint8 {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '+':
		s.addToken(PLUS, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '*':
		s.addToken(STAR, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '/':
		if s.match('/') {
			// 一行注释
			for {
				if s.peek() != '\n' && s.isAtEnd() {
					s.advance()
				}
			}
		} else {
			s.addToken(SLASH, nil)
		}

	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case ' ', '\t', '\r':
	case '\n':
		s.line += 1
	case '"':
		s.string()
	default:
		if isDigital(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			printError(s.line, "unexpected symbol")
		}
	}
}

func isKeyword(text string) (uint, bool) {
	keywordMap := map[string]uint{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	v, ok := keywordMap[text]
	return v, ok
}

func (s *scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	if v, ok := isKeyword(text); ok {
		s.addToken(v, nil)
	} else {
		s.addToken(IDENTIFIER, nil)
	}
}

func (s *scanner) number() {
	for isDigital(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigital(s.peekNext()) {
		s.advance()
		for isDigital(s.peek()) {
			s.advance()
		}
	}
	float64Value, err := parseFloat(s.source[s.start:s.current])
	if err != nil {
		panic(err)
	}
	s.addToken(NUMBER, float64Value)
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		printError(s.line, "unterminated string")
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func isAlphaNumeric(c uint8) bool {
	return isDigital(c) || isAlpha(c)
}

func isDigital(c uint8) bool {
	return c >= '0' && c <= '9'
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func isAlpha(c uint8) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
