package parser

import (
	"csigma/lexer"
	"fmt"
)

// Statement: Interface base. No Go, interfaces vazias permitem que
// diferentes nós (Print, Var, Assignment) sejam armazenados em uma mesma lista.
type Statement interface{}

// --- NÓS DA AST (Modelagem de Dados) ---

// VarDeclNode armazena 'var x = 10'.
type VarDeclNode struct {
	Name  string
	Value string
}

// PrintNode identifica se o que será impresso é uma constante textual ou variável.
type PrintNode struct {
	Value    string
	IsString bool
}

// InputNode mapeia o comando 'input' para o destino na memória.
type InputNode struct {
	VarName string
}

// OrderOp: Estrutura crucial para a aritmética.
// Guarda o operador e o valor seguinte, permitindo cálculos encadeados.
type OrderOp struct {
	Operator string
	Value    string
	IsVar    bool
}

// AssignmentNode: O nó mais complexo. Guarda a variável de destino (Dest),
// o primeiro valor (First) e uma fatia (slice) de todas as operações seguintes.
type AssignmentNode struct {
	Dest  string
	First string
	Ops   []OrderOp
}

// --- ESTRUTURA E LÓGICA DO PARSER ---

type Parser struct {
	tokens []lexer.Token // Lista de tokens vinda do Lexer
	pos    int           // Cursor que indica qual token estamos analisando
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// ParseProgram: O ponto de entrada. Ele decide qual "sub-parser" chamar
// baseado no tipo do token atual (Despacho de Tokens).
func (p *Parser) ParseProgram() ([]Statement, error) {
	var statements []Statement

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != lexer.TokenEOF {
		var stmt Statement
		var err error

		// Padrão de Projeto: Recursive Descent Lite
		switch p.tokens[p.pos].Type {
		case lexer.TokenVar:
			stmt, err = p.parseVarDecl()
		case lexer.TokenPrint:
			stmt, err = p.parsePrint()
		case lexer.TokenInput:
			stmt, err = p.parseInput()
		case lexer.TokenIdent:
			// Se começa com identificador, assumimos ser uma atribuição (ex: res = ...)
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

// parseAssignment: Implementa a "Aritmética Linear".
// Ele lê o destino, o '=' e então entra em um loop para capturar quantos
// operadores e valores existirem na mesma linha.
func (p *Parser) parseAssignment() (Statement, error) {
	// 1. Captura o destino (L-Value)
	dest := p.tokens[p.pos].Literal
	p.pos++

	// 2. Consome o sinal de atribuição
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != lexer.TokenAssign {
		return nil, fmt.Errorf("erro sintático: esperado '=' após identificador '%s'", dest)
	}
	p.pos++

	// 3. Captura o primeiro valor da expressão
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("erro sintático: expressão incompleta")
	}
	first := p.tokens[p.pos].Literal
	p.pos++

	stmt := &AssignmentNode{Dest: dest, First: first}

	// 4. Loop de Operadores: Enquanto houver +, -, * ou /, o parser continua montando a AST.
	// Isso permite que 'a + b * c / d' seja capturado em um único nó.
	for p.pos < len(p.tokens) && isOperator(p.tokens[p.pos].Type) {
		op := p.tokens[p.pos].Literal
		p.pos++

		if p.pos >= len(p.tokens) {
			return nil, fmt.Errorf("erro sintático: esperado valor após operador '%s'", op)
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

// --- FUNÇÕES DE CONSUMO DE TOKENS ---

// parseVarDecl: Transforma 'var x = 0' em um nó estruturado.
func (p *Parser) parseVarDecl() (Statement, error) {
	p.pos++ // pula 'var'
	name := p.tokens[p.pos].Literal
	p.pos++ // pula nome
	p.pos++ // pula '='
	val := p.tokens[p.pos].Literal
	p.pos++ // pula valor
	return &VarDeclNode{Name: name, Value: val}, nil
}

// parsePrint: Identifica o conteúdo do comando print.
func (p *Parser) parsePrint() (Statement, error) {
	p.pos++ // pula 'print'
	isStr := p.tokens[p.pos].Type == lexer.TokenString
	val := p.tokens[p.pos].Literal
	p.pos++
	return &PrintNode{Value: val, IsString: isStr}, nil
}

// isOperator: Helper para definir os limites da expressão aritmética.
func isOperator(t lexer.TokenType) bool {
	return t == lexer.TokenPlus || t == lexer.TokenMinus ||
		t == lexer.TokenMult || t == lexer.TokenDiv
}
