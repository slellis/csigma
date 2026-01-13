package main

import (
	"csigma/lexer"
	"csigma/parser"
	"csigma/semantic" // Importa o novo módulo de análise
	"fmt"
	"os"
)

func main() {
	// 1. Simulação de um código fonte Sigma para teste
	// Comentário didático: Testamos uma declaração e uma atribuição.
	input := `
    	var x = 10
		var y = 20
        y = x / 2.5
	`

	// 2. Fase do Lexer: Transforma o texto em Tokens
	l := lexer.NewLexer(input)
	tokens := l.Tokenize()

	// 3. Fase do Parser: Constrói a Árvore de Sintaxe Abstrata (AST)
	p := parser.NewParser(tokens)
	ast, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("Erro de Sintaxe: %v\n", err)
		os.Exit(1)
	}

	// 4. Fase do Semantic Analyzer: Valida tipos e variáveis (ARCHITECTURE.md item 4 e 5)
	// Comentário didático: Aqui o Sigma verifica se 'y' foi declarado antes de ser usado.
	analyzer := semantic.NewAnalyzer()
	errSem := analyzer.Analisar(ast)

	if errSem != nil {
		// Se houver erros acumulados, o programa para aqui.
		// O método Analisar já imprime a lista detalhada.
		os.Exit(1)
	}

	// 5. Sucesso! Pronto para o próximo passo: CodeGen
	fmt.Println("\n[Sigma] Análise concluída com sucesso! Nenhuma falha detectada.")
}
