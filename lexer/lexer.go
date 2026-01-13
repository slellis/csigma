package lexer

import (
	"unicode"
)

// Definição dos tipos de tokens suportados pelo Sigma.
type TokenType string

const (
	TokenVar     TokenType = "VAR"
	TokenPrint   TokenType = "PRINT"
	TokenInput   TokenType = "INPUT"
	TokenIdent   TokenType = "IDENT"
	TokenInt     TokenType = "INT"
	TokenFloat   TokenType = "FLOAT"  // Suporte a decimais conforme Capítulo 3
	TokenAssign  TokenType = "="
	TokenPlus    TokenType = "+"
	TokenMinus   TokenType = "-"
	TokenMult    TokenType = "*"
	TokenDiv     TokenType = "/"
	TokenString  TokenType = "STRING"
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"
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

// readChar avança um caractere no input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken analisa o caractere atual e retorna o próximo Token.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = Token{Type: TokenAssign, Literal: string(l.ch)}
	case '+':
		tok = Token{Type: TokenPlus,   Literal: string(l.ch)}
	case '-':
		tok = Token{Type: TokenMinus,  Literal: string(l.ch)}
	case '*':
		tok = Token{Type: TokenMult,   Literal: string(l.ch)}
	case '/':
		tok = Token{Type: TokenDiv,    Literal: string(l.ch)}
	case 0:
		tok = Token{Type: TokenEOF,    Literal: ""}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return Token{Type: lookupIdent(literal), Literal: literal}
		} else if isDigit(l.ch) {
			// Comentário didático: readNumber agora decide se é INT ou FLOAT
			return l.readNumber()
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// readNumber diferencia inteiros de decimais.
// Comentário didático: Implementação da lógica de ponto flutuante (Ex: 10.5).
func (l *Lexer) readNumber() Token {
	posInicial := l.position
	temPonto   := false

	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if temPonto { break } // Evita números com dois pontos (10.5.2)
			temPonto = true
		}
		l.readChar()
	}

	literal := l.input[posInicial:l.position]
	tipo    := TokenInt
	if temPonto {
		tipo = TokenFloat
	}
	
	return Token{Type: tipo, Literal: literal}
}

// Tokenize automatiza a coleta de todos os tokens para o Parser.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// --- Funções Auxiliares ---

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		"var":   TokenVar,
		"print": TokenPrint,
		"input": TokenInput,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenIdent
}