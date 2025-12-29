package parser

import (
	"csigma/lexer"
	"fmt"
)

// Interface para qualquer comando (Statement) da nossa linguagem.
type Statement interface{}

// Estrutura para operações encadeadas (ex: + B - C).
type Operation struct {
	Operator lexer.TokenType
	Value    string
}

// Representação dos nós da Árvore Sintática (AST).
type VarDeclNode    struct { Name string; Value string }
type PrintNode      struct { Value string; IsString bool }
type InputNode      struct { VarName string }
type AssignmentNode struct {
	Dest  string      // Nome da variável que recebe o resultado
	First string      // Primeiro valor do cálculo
	Rest  []Operation // Lista de operações seguintes na mesma expressão
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

// ParseProgram percorre a lista de tokens e gera a estrutura lógica do programa.
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
			// Verifica se é uma atribuição matemática (VAR = ...)
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenAssign {
				stmt, err = p.parseAssignment()
			} else {
				return nil, fmt.Errorf("identificador solto no código: %s", tok.Literal)
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

func (p *Parser) parseVar() (Statement, error) {
	p.pos++ // pula VAR
	name := p.tokens[p.pos].Literal
	p.pos++ // pula nome
	p.pos++ // pula =
	val := p.tokens[p.pos].Literal
	p.pos++
	return &VarDeclNode{Name: name, Value: val}, nil
}

func (p *Parser) parsePrint() (Statement, error) {
	p.pos++ // pula PRINT
	tok := p.tokens[p.pos]
	p.pos++
	return &PrintNode{Value: tok.Literal, IsString: tok.Type == lexer.TokenString}, nil
}

func (p *Parser) parseInput() (Statement, error) {
	p.pos++ // pula INPUT
	name := p.tokens[p.pos].Literal
	p.pos++
	return &InputNode{VarName: name}, nil
}

// parseAssignment lida com expressões complexas como A + B * C / D.
func (p *Parser) parseAssignment() (Statement, error) {
	dest := p.tokens[p.pos].Literal
	p.pos += 2 // pula nome e '='

	first := p.tokens[p.pos].Literal
	p.pos++

	var rest []Operation

	// Consome todos os pares de Operador + Valor na linha.
	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos].Type
		if op != lexer.TokenPlus && op != lexer.TokenMinus && op != lexer.TokenMult && op != lexer.TokenDiv {
			break
		}
		p.pos++ // pula o operador
		val := p.tokens[p.pos].Literal
		p.pos++ // pula o valor
		rest = append(rest, Operation{Operator: op, Value: val})
	}

	return &AssignmentNode{Dest: dest, First: first, Rest: rest}, nil
}