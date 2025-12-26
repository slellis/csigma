package parser

import (
	"csigma/lexer"
	"fmt"
)

type Statement interface{}

type VarDeclNode struct {
	Name  string
	Value string
}

type PrintNode struct {
	IsString bool
	Value    string
}

type InputNode struct {
	VarName string
}

type AssignmentNode struct {
	Dest, Left, Operator, Right string
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

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
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenAssign {
				stmt, err = p.parseAssignment()
			} else {
				p.pos++
				continue
			}
		default:
			return nil, fmt.Errorf("token inesperado na posicao %d: %s", p.pos, tok.Literal)
		}
		if err != nil { return nil, err }
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) parseVar() (*VarDeclNode, error) {
	p.pos++ // VAR
	name := p.tokens[p.pos].Literal
	p.pos += 2 // nome e =
	val := p.tokens[p.pos].Literal
	p.pos++
	return &VarDeclNode{Name: name, Value: val}, nil
}

func (p *Parser) parsePrint() (*PrintNode, error) {
	p.pos++ // PRINT
	tok := p.tokens[p.pos]
	p.pos++
	return &PrintNode{IsString: tok.Type == lexer.TokenString, Value: tok.Literal}, nil
}

func (p *Parser) parseInput() (*InputNode, error) {
	p.pos++ // INPUT
	name := p.tokens[p.pos].Literal
	p.pos++
	return &InputNode{VarName: name}, nil
}

func (p *Parser) parseAssignment() (*AssignmentNode, error) {
	dest := p.tokens[p.pos].Literal
	p.pos += 2 // nome e =
	left := p.tokens[p.pos].Literal
	p.pos++
	op := p.tokens[p.pos].Literal
	p.pos++
	right := p.tokens[p.pos].Literal
	p.pos++
	return &AssignmentNode{Dest: dest, Left: left, Operator: op, Right: right}, nil
}