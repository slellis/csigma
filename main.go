package main

import (
	"csigma/codegen"
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// 1. Validação de argumentos e definição de caminhos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <arquivo.sig>")
		return
	}

	inputPath := os.Args[1]
	
	// Didático: filepath.Base extrai apenas o nome do arquivo (ex: "exemplos/calculadora.sig" -> "calculadora.sig")
	// strings.TrimSuffix remove a extensão para termos o nome base do projeto
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ".sig")
	logPath  := baseName + ".log"

	// 2. Abertura do arquivo de Log com MultiWriter (Console + Arquivo)
	logFile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("[ERRO] Falha ao criar log: %v\n", err)
		return
	}
	defer logFile.Close()

	// io.MultiWriter: Tudo enviado para 'mw' vai para o console E para o arquivo de log simultaneamente
	mw := io.MultiWriter(os.Stdout, logFile)
	logPrint := func(format string, a ...interface{}) {
		fmt.Fprintf(mw, format, a...)
	}

	// 3. Início do Processamento
	content, err := os.ReadFile(inputPath)
	if err != nil {
		logPrint("[ERRO] Falha ao ler fonte: %v\n", err)
		return
	}

	logPrint("======================================================================\n")
	logPrint("   CSIGMA COMPILER - RELATORIO DE COMPILACAO\n")
	logPrint("   Data: %s\n", time.Now().Format("02/01/2006 15:04:05"))
	logPrint("   Fonte: %s\n", inputPath)
	logPrint("======================================================================\n")

	// --- FASE 1: ANALISE LEXICA (SCANNER) ---
	logPrint("\n[FASE 1] ANALISE LEXICA: Quebrando o texto em unidades (Tokens)\n")
	logPrint("----------------------------------------------------------------------\n")
	l := lexer.NewLexer(string(content))
	var tokens []lexer.Token
	count := 0
	
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		count++
		logPrint("  [%03d] Tipo: %-12s | Conteudo: \"%s\"\n", count, tok.Type, tok.Literal)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}
	logPrint("--> Sucesso: %d tokens identificados.\n", count)

	// --- FASE 2: ANALISE SINTATICA (PARSER) ---
	logPrint("\n[FASE 2] ANALISE SINTATICA: Construindo a Arvore de Decisao (AST)\n")
	logPrint("----------------------------------------------------------------------\n")
	p := parser.NewParser(tokens)
	statements, err := p.ParseProgram()
	if err != nil {
		logPrint("\n[FATAL] Erro Sintatico detectado:\n  >> %v\n", err)
		return
	}

	for i, stmt := range statements {
		logPrint("  Instrucao %02d: %-25T | Estrutura: %+v\n", i, stmt, stmt)
	}
	logPrint("--> Sucesso: AST montada com %d nos principais.\n", len(statements))

	// --- FASE 3: GERACAO DE CODIGO (CODEGEN) ---
	logPrint("\n[FASE 3] GERACAO DE CODIGO: Traduzindo para Assembly x86_64\n")
	logPrint("----------------------------------------------------------------------\n")
	nasmCode := codegen.GenerateNASM(statements)

	err = os.WriteFile("output.asm", []byte(nasmCode), 0644)
	if err != nil {
		logPrint("[FATAL] Erro ao gravar output.asm: %v\n", err)
		return
	}
	logPrint("--> Sucesso: Arquivo 'output.asm' gerado e comentado.\n")

	// --- FASE 4: MONTAGEM E LINKAGEM ---
	logPrint("\n[FASE 4] MONTAGEM E LINKAGEM: Criando o Executavel Final\n")
	logPrint("----------------------------------------------------------------------\n")
	
	logPrint("  > Rodando NASM (Assembler)... ")
	cmdNasm := exec.Command("nasm", "-f", "elf64", "output.asm", "-o", "output.o")
	if err := cmdNasm.Run(); err != nil {
		logPrint("FALHOU!\n[ERRO]: %v\n", err)
		return
	}
	logPrint("OK.\n")

	logPrint("  > Rodando GCC (Linker)...    ")
	// Dinâmico: O nome do executável agora é o mesmo do arquivo fonte (baseName)
	cmdGcc := exec.Command("gcc", "output.o", "-o", baseName, "-no-pie")
	if err := cmdGcc.Run(); err != nil {
		logPrint("FALHOU!\n[ERRO]: %v\n", err)
		return
	}
	logPrint("OK.\n")

	logPrint("\n======================================================================\n")
	logPrint("   COMPILACAO FINALIZADA COM SUCESSO!\n")
	logPrint("   Arquivo de saida: ./%s\n", baseName)
	logPrint("   Log salvo em:    %s\n", logPath)
	logPrint("======================================================================\n\n")
}