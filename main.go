package main

import (
	"csigma/codegen"
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <arquivo.sig>")
		return
	}

	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[ERRO] Nao foi possivel ler o arquivo: %v\n", err)
		return
	}

	fmt.Printf("--- Compilador CSigma: Processando '%s' ---\n", filePath)

	// 1. Analise Lexica
	l := lexer.NewLexer(string(content))
	var tokens []lexer.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}

	// 2. Analise Sintatica (Parser)
	p := parser.NewParser(tokens)
	statements, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("[ERRO Sintatico] %v\n", err)
		return
	}

	// 3. Geracao de Codigo (Codegen)
	nasmCode := codegen.GenerateNASM(statements)

	// 4. Salva o Assembly em disco
	err = os.WriteFile("output.asm", []byte(nasmCode), 0644)
	if err != nil {
		fmt.Printf("[ERRO] Falha ao salvar output.asm: %v\n", err)
		return
	}

	// 5. Automacao: NASM -> LINKER (GCC)
	fmt.Println("[LOG] Gerando binario executavel...")
	
	// Comando: nasm -f elf64 output.asm -o output.o
	cmdNasm := exec.Command("nasm", "-f", "elf64", "output.asm", "-o", "output.o")
	if err := cmdNasm.Run(); err != nil {
		fmt.Printf("[ERRO] Falha no NASM: %v\n", err)
		return
	}

	// Comando: gcc output.o -o calculadora -no-pie (ou o nome que voce preferir)
	cmdGcc := exec.Command("gcc", "output.o", "-o", "calculadora", "-no-pie")
	if err := cmdGcc.Run(); err != nil {
		fmt.Printf("[ERRO] Falha no Linker (GCC): %v\n", err)
		return
	}

	fmt.Println("--- Processamento Concluido com Sucesso! ---")
	fmt.Println("Execute agora com: ./calculadora")
}