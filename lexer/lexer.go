package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (Points EXACTLY to current char)
	readPosition int  // to look one char ahead of current position
	ch           byte // current char (where position points to)
}

// Returns the Lexer (pointer) and calls readChar to initialize the correct positions
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// Check if we reached the end of input
	// If yes ch is set to 0 which is the ASCII value for NUL
	// That means we either didnt read anything yet or reach the end of the file
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// char is the current char of the input
		l.ch = l.input[l.readPosition]
	}
	// position is the current read char
	l.position = l.readPosition
	// readPosition points to the next char of the input
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// We need to return here and NOT go until the l.readChar() because we
			// already looped and did go over the chars in the input
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	// Reads input until a NON letter occur
	// Lexer position is moving
	for isLetter(l.ch) {
		l.readChar()
	}
	// Return the identifier from position to position
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Only simple integers
// No floats, hex, octal etc.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
