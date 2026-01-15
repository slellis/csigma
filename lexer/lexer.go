package lexer

// TokenType define a categoria do símbolo encontrado.
// Usamos string para que o Log seja legível (ex: "PRINT" em vez de um número 7).
type TokenType string

const (
	// Palavras-Chave (Keywords)
	TokenVar   TokenType = "VAR"
	TokenPrint TokenType = "PRINT"
	TokenInput TokenType = "INPUT"

	// Identificadores e Literais
	TokenIdent  TokenType = "IDENT"  // Nomes de variáveis (ex: soma, res)
	TokenInt    TokenType = "INT"    // Números inteiros (ex: 10, 20)
	TokenFloat  TokenType = "FLOAT"  // Números decimais (ex: 10.5)
	TokenString TokenType = "STRING" // Texto entre aspas "exemplo"

	// Operadores Aritméticos
	TokenAssign TokenType = "="
	TokenPlus   TokenType = "+"
	TokenMinus  TokenType = "-"
	TokenMult   TokenType = "*"
	TokenDiv    TokenType = "/"

	// Pontuação e Delimitadores
	TokenLParen TokenType = "("
	TokenRParen TokenType = ")"
	TokenComma  TokenType = ","
	TokenColon  TokenType = ":"

	// Tokens Especiais
	TokenEOF     TokenType = "EOF"     // End Of File: Fim do arquivo
	TokenIllegal TokenType = "ILLEGAL" // Caractere desconhecido pelo compilador
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string // O código fonte completo
	position     int    // Posição atual do caractere sendo lido (ch)
	readPosition int    // Posição da "espiada" (próximo caractere)
	ch           byte   // Caractere atual sob análise
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Inicializa o lexer lendo o primeiro caractere
	return l
}

// readChar: Avança o ponteiro de leitura.
// ch recebe 0 (ASCII Nul) se chegarmos ao fim do arquivo.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken: O motor principal do Lexer.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace() // Ignora espaços, tabs e quebras de linha

	switch l.ch {
	case '=':
		tok = Token{Type: TokenAssign, Literal: string(l.ch)}
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: TokenMult, Literal: string(l.ch)}
	case '/':
		// Lógica de Comentário: Se encontrarmos '//', ignoramos o resto da linha.
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken() // Reinicia a análise após o comentário
		}
		tok = Token{Type: TokenDiv, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: TokenLParen, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: TokenRParen, Literal: string(l.ch)}
	case '"':
		// Lógica de String: Captura tudo entre aspas.
		tok.Type = TokenString
		tok.Literal = l.readString()
	case ':':
		tok = Token{Type: TokenIllegal, Literal: ":"}
	case ',':
		tok = Token{Type: TokenIllegal, Literal: ","}
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	default:
		// Se for letra, lemos a palavra inteira (pode ser comando ou variável).
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return Token{Type: lookupIdent(literal), Literal: literal}
		} else if isDigit(l.ch) {
			// Se for dígito, lemos o número inteiro ou float.
			return l.readNumber()
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// readNumber: Diferencia 10 (INT) de 10.5 (FLOAT).
// Crucial para garantir a precisão matemática da Versão Platinum.
func (l *Lexer) readNumber() Token {
	posInicial := l.position
	temPonto := false

	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if temPonto {
				break
			} // Proteção: Um número não pode ter dois pontos.
			temPonto = true
		}
		l.readChar()
	}

	literal := l.input[posInicial:l.position]
	tipo := TokenInt
	if temPonto {
		tipo = TokenFloat
	}

	return Token{Type: tipo, Literal: literal}
}

// skipComment: Avança o ponteiro até encontrar '\n', efetivamente ignorando o comentário.
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

// readString: Captura o conteúdo entre as aspas, sem incluir as próprias aspas no Literal.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// peekChar: Olha o próximo caractere sem mover o ponteiro principal (essencial para '//').
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// lookupIdent: Verifica se uma palavra é um comando do Sigma (var, print, input) ou uma variável.
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

// ... isLetter, isDigit, readIdentifier e skipWhitespace seguem lógica padrão de leitura de texto.
