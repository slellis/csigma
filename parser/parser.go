package parser

import (
	"csigma/lexer"
	"fmt"
)

// --- DEFINIÇÃO DOS NÓS DA AST (ÁRBORE SINTÁTICA) ---

type Statement interface{}

// VarDeclNode representa a declaração: VAR NOME = VALOR
type VarDeclNode struct {
	Name  string
	Value string
}

// PrintNode representa o comando: PRINT "TEXTO" ou PRINT VAR
type PrintNode struct {
	Value    string
	IsString bool
}

// InputNode representa o comando: INPUT VAR
type InputNode struct {
	VarName string
}

// AssignmentNode representa a operação aritmética: C = A + B ou C = A - B
type AssignmentNode struct {
	Dest     string
	Left     string
	Right    string
	Operator string // Novo campo: Guardará se é TokenPlus ou TokenMinus
}

// Parser é o motor que converte Tokens em uma estrutura lógica (AST).
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// NewParser cria uma nova instância do analisador sintático.
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

// ParseProgram é o laço principal que percorre todos os tokens do arquivo .sig.
func (p *Parser) ParseProgram() ([]Statement, error) {
	var statements []Statement

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != lexer.TokenEOF {
		var stmt Statement
		var err error
		tok := p.tokens[p.pos]

		switch tok.Type {
		case lexer.TokenVar:
			stmt, err = p.parseVar()
		case lexer.TokenPrint:
			stmt, err = p.parsePrint()
		case lexer.TokenInput:
			stmt, err = p.parseInput()
		case lexer.TokenIdent:
			// Se começar com Identificador, verificamos se o próximo é um '=' para Atribuição
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenAssign {
				stmt, err = p.parseAssignment()
			} else {
				return nil, fmt.Errorf("Identificador '%s' fora de contexto na posicao %d", tok.Literal, p.pos)
			}
		default:
			// Se encontrar algo que não conhece, interrompe e avisa o erro (Rigor Sintático)
			return nil, fmt.Errorf("Token inesperado '%s' (Tipo: %s) na posicao %d", tok.Literal, tok.Type, p.pos)
		}

		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// parseVar lida com a gramática: VAR <ident> = <number>
func (p *Parser) parseVar() (Statement, error) {
	p.pos++ // pula 'VAR'
	name := p.tokens[p.pos].Literal
	p.pos++ // pula nome
	p.pos++ // pula '='
	val := p.tokens[p.pos].Literal
	p.pos++ // pula valor
	return &VarDeclNode{Name: name, Value: val}, nil
}

// parsePrint lida com a gramática: PRINT <string> ou PRINT <ident>
func (p *Parser) parsePrint() (Statement, error) {
	p.pos++ // pula 'PRINT'
	tok := p.tokens[p.pos]
	isString := tok.Type == lexer.TokenString
	p.pos++
	return &PrintNode{Value: tok.Literal, IsString: isString}, nil
}

// parseInput lida com a gramática: INPUT <ident>
func (p *Parser) parseInput() (Statement, error) {
	p.pos++ // pula 'INPUT'
	name := p.tokens[p.pos].Literal
	p.pos++
	return &InputNode{VarName: name}, nil
}

// parseAssignment lida com a gramática: <ident> = <ident> [+ ou -] <ident>
func (p *Parser) parseAssignment() (Statement, error) {
	dest := p.tokens[p.pos].Literal
	p.pos += 2 // pula o destino e o '='

	left := p.tokens[p.pos].Literal
	p.pos++

	// Captura o operador (+ ou -) para que o Codegen saiba o que fazer
	operator := p.tokens[p.pos].Type
	p.pos++

	right := p.tokens[p.pos].Literal
	p.pos++

	return &AssignmentNode{
		Dest:     dest,
		Left:     left,
		Right:    right,
		Operator: operator,
	}, nil
}