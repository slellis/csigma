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
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <arquivo.sig>")
		return
	}

	inputPath := os.Args[1]
	dir := filepath.Dir(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ".sig")
	logPath := filepath.Join(dir, baseName+".log")

	logFile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("Erro ao criar log: %v\n", err)
		return
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	logPrint := func(f string, a ...interface{}) { fmt.Fprintf(mw, f, a...) }

	content, _ := os.ReadFile(inputPath)

	logPrint("======================================================================\n")
	logPrint("   CSIGMA PLATINUM - RELATORIO TECNICO DE COMPILACAO\n")
	logPrint("   Data: %s\n", time.Now().Format("02/01/2006 15:04:05"))
	logPrint("   Fonte: %s\n", inputPath)
	logPrint("======================================================================\n")

	// --- FASE 1: LEXER ---
	logPrint("\n[FASE 1] ANALISE LEXICA (Scanner):\n")
	logPrint("----------------------------------------------------------------------\n")
	l := lexer.NewLexer(string(content))
	var tokens []lexer.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		logPrint("  Token: [%-12s] | Literal: \"%s\"\n", tok.Type, tok.Literal)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}

	// --- FASE 2: PARSER ---
	logPrint("\n[FASE 2] ANALISE SINTATICA (Dump da AST):\n")
	logPrint("----------------------------------------------------------------------\n")
	p := parser.NewParser(tokens)
	statements, err := p.ParseProgram()
	if err != nil {
		logPrint("\n[ERRO SINTATICO] %v\n", err)
		return
	}

	for i, stmt := range statements {
		switch s := stmt.(type) {
		case *parser.VarDeclNode:
			logPrint("  [%02d] DECLARACAO:  Var %s = %s\n", i, s.Name, s.Value)
		case *parser.PrintNode:
			logPrint("  [%02d] PRINT:       \"%s\" (String: %v)\n", i, s.Value, s.IsString)
		case *parser.InputNode:
			logPrint("  [%02d] INPUT:       Ler para variavel %s\n", i, s.VarName)
		case *parser.AssignmentNode:
			exprStr := s.First
			for _, op := range s.Ops {
				exprStr += " " + op.Operator + " " + op.Value
			}
			logPrint("  [%02d] CALCULO:     %s = %s\n", i, s.Dest, exprStr)
		}
	}

	// --- FASE 3: CODEGEN (Com listagem no log) ---
	logPrint("\n[FASE 3] GERACAO DE CODIGO (Assembly x86_64):\n")
	logPrint("----------------------------------------------------------------------\n")
	nasmCode := codegen.GenerateNASM(statements)

	// Grava no log o código gerado
	logPrint("%s\n", nasmCode)

	os.WriteFile("output.asm", []byte(nasmCode), 0644)
	logPrint("----------------------------------------------------------------------\n")
	logPrint("  > [OK] Arquivo 'output.asm' gravado no disco.\n")

	// --- FASE 4: BUILD ---
	logPrint("\n[FASE 4] MONTAGEM E LINKAGEM (NASM & GCC):\n")
	logPrint("----------------------------------------------------------------------\n")

	logPrint("  > Executando NASM... ")
	if err := exec.Command("nasm", "-f", "elf64", "output.asm", "-o", "output.o").Run(); err != nil {
		logPrint("FALHOU: %v\n", err)
		return
	}
	logPrint("OK.\n")

	logPrint("  > Executando GCC...  ")
	if err := exec.Command("gcc", "output.o", "-o", baseName, "-no-pie").Run(); err != nil {
		logPrint("FALHOU: %v\n", err)
		return
	}
	logPrint("OK.\n")

	logPrint("\n======================================================================\n")
	logPrint("   RESULTADO FINAL: ./%s\n", baseName)
	logPrint("   Log gerado em:   %s\n", logPath)
	logPrint("======================================================================\n")
}
