package main

import (
	"csigma/codegen"
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 1. Validar se o arquivo fonte foi passado como argumento
	if len(os.Args) < 2 {
		fmt.Println("Uso: csigma <arquivo.sig>")
		return
	}

	caminhoFonte := os.Args[1]
	nomeExecutavel := strings.TrimSuffix(filepath.Base(caminhoFonte), filepath.Ext(caminhoFonte))

	fmt.Printf("--- Compilador CSigma: Iniciando Processamento de '%s' ---\n", caminhoFonte)

	// 2. Leitura do Arquivo Fonte
	fmt.Println("[LOG] Lendo arquivo fonte...")
	conteudo, err := os.ReadFile(caminhoFonte)
	if err != nil {
		fmt.Printf("[ERRO] Nao foi possivel ler o arquivo: %v\n", err)
		return
	}

	// 3. Analise Lexica
	fmt.Println("[LOG] Executando Analise Lexica (Tokens)...")
	l := lexer.NewLexer(string(conteudo))
	var tokens []lexer.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}

	// 4. Analise Sintatica (Parser)
	fmt.Println("[LOG] Construindo Arvore Sintatica (AST)...")
	p := parser.NewParser(tokens)
	statements, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("[ERRO] Erro Sintatico: %v\n", err)
		return
	}

	// 5. Geracao de Codigo Assembly
	fmt.Println("[LOG] Traduzindo para Assembly x86_64...")
	asm := codegen.GenerateNASM(statements)
	arquivoAsm := "output.asm"
	os.WriteFile(arquivoAsm, []byte(asm), 0644)

	// --- AUTOMACAO DOS PASSOS FINAIS ---

	// 6. Executar NASM (Assemble)
	fmt.Println("[NASM] Convertendo Assembly para Objeto (output.o)...")
	cmdNasm := exec.Command("nasm", "-f", "elf64", arquivoAsm, "-o", "output.o")
	if err := rodarComando(cmdNasm); err != nil {
		fmt.Printf("[ERRO] Falha no NASM: %v\n", err)
		return
	}

	// 7. Executar GCC (Link-edit)
	fmt.Printf("[GCC] Linkando e gerando executavel final: '%s'...\n", nomeExecutavel)
	cmdGcc := exec.Command("gcc", "output.o", "-o", nomeExecutavel, "-no-pie")
	if err := rodarComando(cmdGcc); err != nil {
		fmt.Printf("[ERRO] Falha no GCC: %v\n", err)
		return
	}

	// 8. Limpeza e Finalizacao
	fmt.Println("[LOG] Limpando arquivos temporarios...")
	os.Remove("output.o")
	// os.Remove(arquivoAsm) // Descomente se quiser deletar o .asm automaticamente

	fmt.Println("--------------------------------------------------")
	fmt.Printf("SUCESSO! Programa '%s' gerado com exito.\n", nomeExecutavel)
	fmt.Printf("Para rodar, digite: ./%s\n", nomeExecutavel)
}

// Funcao auxiliar para rodar comandos externos e capturar erros
func rodarComando(cmd *exec.Cmd) error {
	saida, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Saida do erro:\n%s\n", string(saida))
		return err
	}
	return nil
}