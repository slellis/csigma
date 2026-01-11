package lexer

type TokenType string

// Definição dos símbolos e palavras-chave que o CSigma entende.
const (
	TokenVar      = "var"
	TokenIdent    = "IDENT"
	TokenAssign   = "="
	TokenNumber   = "NUMBER"
	TokenPlus     = "+"
	TokenMinus    = "-"
	TokenMult     = "*"
	TokenDiv      = "/"   // Operador de divisão aritmética
	TokenPrint    = "print"
	TokenInput    = "input"
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

// keywords define a "Fonte Única da Verdade" para as palavras reservadas do Sigma.
// Comentário didático: Mapeamos a string exata que o usuário escreve para a constante do token.
var keywords = map[string]TokenType{
	"var":   TokenVar,
	"print": TokenPrint,
	"input": TokenInput,
	// Se amanhã você criar o comando "if", basta adicionar a linha abaixo:
	// "if": TokenIf, 
}


func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar move o cursor de leitura para o próximo caractere do arquivo fonte.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // Indica Fim de Arquivo (End of File)
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar permite olhar o próximo caractere sem avançar o cursor (usado para '//').
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken identifica qual o próximo símbolo válido no código.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	// SUPORTE A COMENTÁRIOS: Ignora o texto se encontrar '//'.
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
		return l.NextToken()
	}

	var tok Token
	switch l.ch {
	case '=': tok = Token{Type: TokenAssign, Literal: string(l.ch)}
	case '+': tok = Token{Type: TokenPlus, Literal: string(l.ch)}
	case '-': tok = Token{Type: TokenMinus, Literal: string(l.ch)}
	case '*': tok = Token{Type: TokenMult, Literal: string(l.ch)}
	case '/': tok = Token{Type: TokenDiv, Literal: string(l.ch)}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			// literal := l.readIdentifier()
			// tok.Type = lookupIdent(literal)
			// tok.Literal = literal
			// return tok
			palavra := l.readIdentifier() // ou o nome da sua função que lê a palavra
    		// AQUI é onde a mágica acontece:
    		tipoDoToken := lookupIdent(palavra) 
    		return Token{Type: tipoDoToken, Literal: palavra}
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

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) { l.readChar() }
	return l.input[position:l.position]
}

func isLetter(ch byte) bool { return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' }
func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }

// func lookupIdent(ident string) TokenType {
// 	switch ident {
// 	case "var": return TokenVar
// 	case "print": return TokenPrint
// 	case "input": return TokenInput
// 	default: return TokenIdent
// 	}
// }

// Substitua sua função antiga por esta:
func lookupIdent(ident string) TokenType{
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return TokenIdent
}