package parser

import (
	"csigma/lexer"
	"fmt"
)

// Interface para qualquer comando (Statement) da nossa linguagem.
type Statement interface{}

// Estrutura para operações encadeadas (ex: A + B - C).
type Operation struct {
	Operator lexer.TokenType
	Value    string
}

// Representação dos nós da Árvore Sintática (AST).
// Comentário didático: O campo Value agora é interface{} para suportar 
// diferentes tipos literais que o Analisador Semântico irá validar.
type VarDeclNode struct {
	Name  string
	Value interface{}
}

type PrintNode struct {
	Value    string
	IsString bool
}

type InputNode struct {
	VarName string
}

type AssignmentNode struct {
	Dest  string      // Variável que recebe o resultado
	First string      // Primeiro operando
	Rest  []Operation // Lista de operações seguintes
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

// ParseProgram percorre os tokens e constrói a AST.
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
			// Verifica se é uma atribuição (ID = ...)
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenAssign {
				stmt, err = p.parseAssignment()
			} else {
				return nil, fmt.Errorf("identificador fora de contexto: %s", tok.Literal)
			}
		default:
			p.pos++
			continue
		}

		if err != nil { return nil, err }
		statements = append(statements, stmt)
	}
	return statements, nil
}

// parseVar lida com a declaração: var x = 10
func (p *Parser) parseVar() (Statement, error) {
	p.pos++ // pula 'var'
	name := p.tokens[p.pos].Literal
	p.pos++ // pula nome
	p.pos++ // pula '='
	
	val := p.tokens[p.pos].Literal
	p.pos++
	
	return &VarDeclNode{Name: name, Value: val}, nil
}

// parsePrint lida com a exibição de dados
func (p *Parser) parsePrint() (Statement, error) {
	p.pos++ // pula 'print'
	tok := p.tokens[p.pos]
	p.pos++
	return &PrintNode{Value: tok.Literal, IsString: tok.Type == lexer.TokenString}, nil
}

// parseInput lida com a entrada de dados
func (p *Parser) parseInput() (Statement, error) {
	p.pos++ // pula 'input'
	name := p.tokens[p.pos].Literal
	p.pos++
	return &InputNode{VarName: name}, nil
}

// parseAssignment lida com expressões matemáticas complexas
func (p *Parser) parseAssignment() (Statement, error) {
	dest := p.tokens[p.pos].Literal
	p.pos += 2 // pula nome e '='

	first := p.tokens[p.pos].Literal
	p.pos++

	var rest []Operation

	// Loop para capturar operadores e operandos na mesma linha
	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos].Type
		// Verifica se o token atual é um operador aritmético
		if op != lexer.TokenPlus && op != lexer.TokenMinus && 
		   op != lexer.TokenMult && op != lexer.TokenDiv {
			break
		}
		
		p.pos++ // pula o operador
		val := p.tokens[p.pos].Literal
		p.pos++ // pula o valor/identificador
		
		rest = append(rest, Operation{Operator: op, Value: val})
	}

	return &AssignmentNode{Dest: dest, First: first, Rest: rest}, nil
}