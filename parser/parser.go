package parser

import (
	"csigma/lexer"
	"fmt"
)

type Statement interface{}

type Operation struct {
	Operator lexer.TokenType
	Value    string
}

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

type AssignmentNode struct {
	Dest  string
	First string
	Rest  []Operation
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

	// O laço deve rodar enquanto houver tokens e não for o fim do arquivo (EOF)
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
				return nil, fmt.Errorf("identificador fora de contexto: %s", tok.Literal)
			}
		default:
			// Se encontrar um token desconhecido (como um resto de linha), apenas avança
			p.pos++
			continue
		}

		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
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

func (p *Parser) parseAssignment() (Statement, error) {
	dest := p.tokens[p.pos].Literal
	p.pos += 2 // pula destino e '='

	first := p.tokens[p.pos].Literal
	p.pos++

	var rest []Operation

	// Consome todos os pares (Operador + Identificador)
	for p.pos < len(p.tokens) {
		tokType := p.tokens[p.pos].Type
		if tokType != lexer.TokenPlus && tokType != lexer.TokenMinus {
			break
		}

		op := tokType
		p.pos++ // pula operador

		if p.pos >= len(p.tokens) { break }
		val := p.tokens[p.pos].Literal
		p.pos++ // pula valor

		rest = append(rest, Operation{Operator: op, Value: val})
	}

	return &AssignmentNode{
		Dest:  dest,
		First: first,
		Rest:  rest,
	}, nil
}