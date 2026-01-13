package semantic

import (
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"strconv"
)

type Simbolo struct {
	Nome string
	Tipo string // SIGMA_INT ou SIGMA_FLT
}

type SemanticAnalyzer struct {
	TabelaSimbolos map[string]Simbolo
	Erros          []string
}

func NewAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		TabelaSimbolos: make(map[string]Simbolo),
		Erros:          []string{},
	}
}

func (a *SemanticAnalyzer) Analisar(statements []parser.Statement) error {
	for _, stmt := range statements {
		switch n := stmt.(type) {
		case *parser.VarDeclNode:
			a.validarDeclaracao(n)
		case *parser.AssignmentNode:
			a.validarAtribuicao(n)
		}
	}

	if len(a.Erros) > 0 {
		fmt.Println("\n--- ERROS SEMÂNTICOS ENCONTRADOS ---")
		for _, err := range a.Erros {
			fmt.Printf(">> %s\n", err)
		}
		return fmt.Errorf("falha na análise: %d erro(s)", len(a.Erros))
	}
	return nil
}

func (a *SemanticAnalyzer) validarDeclaracao(n *parser.VarDeclNode) {
	if _, existe := a.TabelaSimbolos[n.Name]; existe {
		a.Erros = append(a.Erros, fmt.Sprintf("variável '%s' já declarada", n.Name))
		return
	}

	// Comentário didático: Identifica o tipo do valor inicial (10 vs 10.5)
	tipo := a.inferirTipo(fmt.Sprintf("%v", n.Value))
	a.TabelaSimbolos[n.Name] = Simbolo{Nome: n.Name, Tipo: tipo}
}

// func (a *SemanticAnalyzer) validarAtribuicao(n *parser.AssignmentNode) {
// 	simboloDest, existe := a.TabelaSimbolos[n.Dest]
// 	if !existe {
// 		a.Erros = append(a.Erros, fmt.Sprintf("variável '%s' não declarada", n.Dest))
// 		return
// 	}

// 	tipoExpressao := a.obterTipoDoValor(n.First)

// 	// Regra B: Verificação de Tipagem Fixa
// 	if tipoExpressao != simboloDest.Tipo {
// 		a.Erros = append(a.Erros, fmt.Sprintf("conflito de tipos em '%s': esperado %s, recebeu %s",
// 			n.Dest, simboloDest.Tipo, tipoExpressao))
// 	}

// 	for _, op := range n.Rest {
// 		tipoOperando := a.obterTipoDoValor(op.Value)

// 		// Opção A: Divisão Estrita (Tipos devem ser iguais)
// 		if op.Operator == lexer.TokenDiv && tipoExpressao != tipoOperando {
// 			a.Erros = append(a.Erros, fmt.Sprintf("divisão ilegal: não é possível dividir %s por %s",
// 				tipoExpressao, tipoOperando))
// 		}

// 		if tipoExpressao != tipoOperando {
// 			a.Erros = append(a.Erros, fmt.Sprintf("operação inválida: %s e %s são tipos incompatíveis",
// 				tipoExpressao, tipoOperando))
// 		}
// 	}
// }

func (a *SemanticAnalyzer) validarAtribuicao(n *parser.AssignmentNode) {
	// 1. Verifica se a variável de destino existe
	simboloDest, existe := a.TabelaSimbolos[n.Dest]
	if !existe {
		a.Erros = append(a.Erros, fmt.Sprintf("variável '%s' não declarada", n.Dest))
		return
	}

	// 2. Obtém o tipo do primeiro operando (ex: 'x' em 'x / 2.5')
	tipoExpressao := a.obterTipoDoValor(n.First)

	// 3. Verifica as operações seguintes (o que vem depois do operador)
	for _, op := range n.Rest {
		tipoOperando := a.obterTipoDoValor(op.Value)

		// LOG DIDÁTICO: Vamos ver o que o Sigma está comparando
		// fmt.Printf("[DEBUG] Comparando %s (%s) com %s (%s)\n", n.First, tipoExpressao, op.Value, tipoOperando)

		// Opção A: Divisão Estrita - Se for divisão, os tipos TEM que ser iguais
		if op.Operator == lexer.TokenDiv {
			if tipoExpressao != tipoOperando {
				a.Erros = append(a.Erros, fmt.Sprintf("Divisão Inválida: '%s' é %s, mas '%s' é %s",
					n.First, tipoExpressao, op.Value, tipoOperando))
			}
		}

		// Regra Geral: Não permitimos mistura de tipos em nenhuma operação aritmética no Sigma
		if tipoExpressao != tipoOperando {
			a.Erros = append(a.Erros, fmt.Sprintf("Tipo Incompatível: não pode operar %s com %s",
				tipoExpressao, tipoOperando))
		}
	}

	// 4. Regra B: O resultado da expressão deve caber no tipo da variável de destino
	if tipoExpressao != simboloDest.Tipo {
		a.Erros = append(a.Erros, fmt.Sprintf("Conflito de Atribuição: '%s' é %s, mas recebeu %s",
			n.Dest, simboloDest.Tipo, tipoExpressao))
	}
}

// inferirTipo verifica se a string é um inteiro ou decimal.
func (a *SemanticAnalyzer) inferirTipo(valor string) string {
	if _, err := strconv.Atoi(valor); err == nil {
		return "SIGMA_INT"
	}
	if _, err := strconv.ParseFloat(valor, 64); err == nil {
		return "SIGMA_FLT"
	}
	return "SIGMA_UNKNOWN"
}

func (a *SemanticAnalyzer) obterTipoDoValor(v string) string {
	if s, existe := a.TabelaSimbolos[v]; existe {
		return s.Tipo
	}
	return a.inferirTipo(v)
}
