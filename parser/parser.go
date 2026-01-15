package parser

import (
	"csigma/lexer"
	"fmt"
)

// Statement: Interface base para todos os nós da AST
type Statement interface{}

// --- NÓS DA AST ---

type VarDeclNode struct {
	Name  string
	Value string
}

type PrintNode struct {
	Value    string
	IsString bool
}

type InputNode struct {
	VarName string
}

// OrderOp: Representa um par (operador, valor) na expressão
type OrderOp struct {
	Operator string
	Value    string
	IsVar    bool
}

// AssignmentNode: Agora suporta a lista 'Ops' para cálculos lineares
type AssignmentNode struct {
	Dest  string
	First string
	Ops   []OrderOp
}

// --- ESTRUTURA DO PARSER ---

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// ParseProgram: O coração do Parser que percorre os tokens
func (p *Parser) ParseProgram() ([]Statement, error) {
	var statements []Statement

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != lexer.TokenEOF {
		var stmt Statement
		var err error

		switch p.tokens[p.pos].Type {
		case lexer.TokenVar:
			stmt, err = p.parseVarDecl()
		case lexer.TokenPrint:
			stmt, err = p.parsePrint()
		case lexer.TokenInput:
			stmt, err = p.parseInput()
		case lexer.TokenIdent:
			// Aqui estava o erro de mismatch; agora tratamos (stmt, err)
			stmt, err = p.parseAssignment()
		default:
			p.pos++
			continue
		}

		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// parseAssignment: Captura 'res = a + b * 2 / c'
func (p *Parser) parseAssignment() (Statement, error) {
	// 1. Destino (ex: res)
	dest := p.tokens[p.pos].Literal
	p.pos++

	// 2. Pula o '='
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != lexer.TokenAssign {
		return nil, fmt.Errorf("esperado '=' após identificador")
	}
	p.pos++

	// 3. Primeiro valor (ex: a)
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("expressão incompleta após '='")
	}
	first := p.tokens[p.pos].Literal
	p.pos++

	stmt := &AssignmentNode{Dest: dest, First: first}

	// 4. Loop de operações (+ b * 2 / c)
	for p.pos < len(p.tokens) && isOperator(p.tokens[p.pos].Type) {
		op := p.tokens[p.pos].Literal
		p.pos++

		if p.pos >= len(p.tokens) {
			return nil, fmt.Errorf("esperando valor após operador %s", op)
		}

		val := p.tokens[p.pos].Literal
		isVar := p.tokens[p.pos].Type == lexer.TokenIdent

		stmt.Ops = append(stmt.Ops, OrderOp{
			Operator: op,
			Value:    val,
			IsVar:    isVar,
		})
		p.pos++
	}

	return stmt, nil
}

// --- FUNÇÕES DIDÁTICAS AUXILIARES ---

func (p *Parser) parseVarDecl() (Statement, error) {
	p.pos++ // pula 'var'
	name := p.tokens[p.pos].Literal
	p.pos++ // pula nome
	p.pos++ // pula '='
	val := p.tokens[p.pos].Literal
	p.pos++ // pula valor
	return &VarDeclNode{Name: name, Value: val}, nil
}

func (p *Parser) parsePrint() (Statement, error) {
	p.pos++ // pula 'print'
	isStr := p.tokens[p.pos].Type == lexer.TokenString
	val := p.tokens[p.pos].Literal
	p.pos++
	return &PrintNode{Value: val, IsString: isStr}, nil
}

func (p *Parser) parseInput() (Statement, error) {
	p.pos++ // pula 'input'
	name := p.tokens[p.pos].Literal
	p.pos++
	return &InputNode{VarName: name}, nil
}

func isOperator(t lexer.TokenType) bool {
	return t == lexer.TokenPlus || t == lexer.TokenMinus ||
		t == lexer.TokenMult || t == lexer.TokenDiv
}
