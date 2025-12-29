package lexer

type TokenType string

const (
	TokenVar      = "VAR"
	TokenIdent    = "IDENT"
	TokenAssign   = "="
	TokenNumber   = "NUMBER"
	TokenPlus     = "+"
	TokenMinus    = "-"
	TokenPrint    = "PRINT"
	TokenInput    = "INPUT"
	TokenString   = "STRING"
	TokenEOF      = "EOF"
)

type Token struct {
	Type    TokenType
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
		l.ch = 0 // Fim da string
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	// Proteção contra loop infinito no fim do arquivo
	if l.ch == 0 {
		return Token{Type: TokenEOF, Literal: ""}
	}

	// Comentários //
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		return l.NextToken()
	}

	var tok Token
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
		tok.Type = TokenEOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = lookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: "ILLEGAL", Literal: string(l.ch)}
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 { break }
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { l.readChar() }
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) { l.readChar() }
	return l.input[position:l.position]
}

func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }

func lookupIdent(ident string) TokenType {
	switch ident {
	case "VAR": return TokenVar
	case "PRINT": return TokenPrint
	case "INPUT": return TokenInput
	default: return TokenIdent
	}
}