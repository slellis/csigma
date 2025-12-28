package lexer

import "unicode"

const (
	TokenVar    = "VAR"
	TokenPrint  = "PRINT"
	TokenInput  = "INPUT"
	TokenIdent  = "IDENT"
	TokenString = "STRING"
	TokenAssign = "ASSIGN"
	TokenPlus   = "PLUS"
	TokenMinus  = "MINUS"  // Representa o operador '-'
	TokenEOF    = "EOF"
)

type Token struct {
	Type    string
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = Token{Type: TokenAssign, Literal: string(l.ch)}
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch)}
	case '-':
    	tok = Token{Type: TokenMinus, Literal: string(l.ch)}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	default:
		if unicode.IsLetter(rune(l.ch)) {
			ident := l.readIdentifier()
			tok.Literal = ident
			switch ident {
			case "VAR": tok.Type = TokenVar
			case "PRINT": tok.Type = TokenPrint
			case "INPUT": tok.Type = TokenInput
			default: tok.Type = TokenIdent
			}
			return tok
		}
		tok = Token{Type: "INVALID", Literal: string(l.ch)}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	l.readChar() // pula " inicial
	pos := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	s := l.input[pos:l.position]
	l.readChar() // pula " final
	return s
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for unicode.IsLetter(rune(l.ch)) || unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}