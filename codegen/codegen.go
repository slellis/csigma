package codegen

import (
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"strings"
	"unicode"
)

// isNumber verifica se a string é um valor numérico puro.
func isNumber(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) { return false }
	}
	return true
}

// GenerateNASM gera o arquivo .asm com 100% de cobertura de comentários.
func GenerateNASM(statements []parser.Statement) string {
	variables := make(map[string]string)
	var stringsConst []string
	var instructions []string
	msgCount := 0

	// FASE 1: Mapeamento de Instruções
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *parser.VarDeclNode:
			variables[s.Name] = s.Value

		case *parser.PrintNode:
			if s.IsString {
				msgName := fmt.Sprintf("msg_%d", msgCount)
				lineMsg := fmt.Sprintf("    %-20s db '%s', 10, 0", msgName, s.Value)
				stringsConst = append(stringsConst, lineMsg)
				
				instructions = append(instructions, fmt.Sprintf("%-40s; Endereco da string para RDI", "    lea rdi, ["+msgName+"]"))
				instructions = append(instructions, fmt.Sprintf("%-40s; Limpa EAX (sem args de ponto flutuante)", "    xor eax, eax"))
				instructions = append(instructions, fmt.Sprintf("%-40s; Chama a funcao printf da LibC", "    call printf"))
				msgCount++
			} else {
				instructions = append(instructions, fmt.Sprintf("%-40s; Formato de saida em RDI", "    lea rdi, [fmt_out_num]"))
				instructions = append(instructions, fmt.Sprintf("%-40s; Valor da variavel em RSI", "    mov rsi, ["+s.Value+"]"))
				instructions = append(instructions, fmt.Sprintf("%-40s; Limpa EAX para printf", "    xor eax, eax"))
				instructions = append(instructions, fmt.Sprintf("%-40s; Exibe o valor numerico", "    call printf"))
			}

		case *parser.InputNode:
			instructions = append(instructions, fmt.Sprintf("%-40s; Formato de entrada em RDI", "    lea rdi, [fmt_in]"))
			instructions = append(instructions, fmt.Sprintf("%-40s; Endereco de destino em RSI", "    lea rsi, ["+s.VarName+"]"))
			instructions = append(instructions, fmt.Sprintf("%-40s; Limpa EAX para scanf", "    xor eax, eax"))
			instructions = append(instructions, fmt.Sprintf("%-40s; Captura entrada do teclado", "    call scanf"))

		case *parser.AssignmentNode:
			// Carga inicial no Acumulador
			if isNumber(s.First) {
				instructions = append(instructions, fmt.Sprintf("%-40s; Carrega valor imediato %s em RAX", "    mov rax, "+s.First, s.First))
			} else {
				instructions = append(instructions, fmt.Sprintf("%-40s; Carrega conteudo de %s em RAX", "    mov rax, ["+s.First+"]", s.First))
			}
			
			for _, op := range s.Rest {
				instr, label, val := "", "", ""
				if isNumber(op.Value) { val = op.Value } else { val = "[" + op.Value + "]" }

				switch op.Operator {
				case lexer.TokenPlus:  instr = "add"; label = "Soma"
				case lexer.TokenMinus: instr = "sub"; label = "Subtrai"
				case lexer.TokenMult:  instr = "imul"; label = "Multiplica"
				case lexer.TokenDiv:
					instructions = append(instructions, fmt.Sprintf("%-40s; Carrega divisor em RBX", "    mov rbx, "+val))
					instructions = append(instructions, fmt.Sprintf("%-40s; Limpa RDX para divisao segura", "    xor rdx, rdx"))
					instructions = append(instructions, fmt.Sprintf("%-40s; Estende sinal de RAX para RDX:RAX", "    cqo"))
					instructions = append(instructions, fmt.Sprintf("%-40s; Divide RDX:RAX por RBX (Quociente->RAX)", "    idiv rbx"))
					continue
				}
				if instr != "" {
					line := fmt.Sprintf("    %s rax, %s", instr, val)
					instructions = append(instructions, fmt.Sprintf("%-40s; %s %s", line, label, op.Value))
				}
			}
			instructions = append(instructions, fmt.Sprintf("%-40s; Armazena resultado em %s", "    mov ["+s.Dest+"], rax", s.Dest))
		}
	}

	// FASE 2: Montagem Final do Arquivo
	var output strings.Builder
	output.WriteString("section .data                           ; --- SECAO DE DADOS (WORKING-STORAGE) ---\n")
	output.WriteString(fmt.Sprintf("%-40s; Formato scanf\n", "    fmt_in db ' %ld', 0"))
	output.WriteString(fmt.Sprintf("%-40s; Formato printf\n", "    fmt_out_num db '%ld', 10, 0"))
	
	for name, val := range variables {
		line := fmt.Sprintf("    %-20s dq %s", name, val)
		output.WriteString(fmt.Sprintf("%-40s; Alocacao da variavel %s\n", line, name))
	}
	for _, s := range stringsConst {
		output.WriteString(s + "               ; Texto para exibicao\n")
	}

	output.WriteString("\nsection .text                           ; --- SECAO DE CODIGO (PROCEDURE DIVISION) ---\n")
	output.WriteString("extern printf, scanf                    ; Funcoes da biblioteca C padrao\n")
	output.WriteString("global main                             ; Ponto de entrada do executavel\n\n")
	output.WriteString("main:\n")
	output.WriteString(fmt.Sprintf("%-40s; Salva o ponteiro de base da pilha\n", "    push rbp"))
	output.WriteString(fmt.Sprintf("%-40s; Alinha o ponteiro de base\n", "    mov rbp, rsp"))
	output.WriteString(fmt.Sprintf("%-40s; Reserva espaco e alinha stack em 16 bytes\n\n", "    sub rsp, 32"))

	for _, ins := range instructions {
		output.WriteString(ins + "\n")
	}

	output.WriteString(fmt.Sprintf("\n%-40s; Libera espaco da pilha\n", "    add rsp, 32"))
	output.WriteString(fmt.Sprintf("%-40s; Restaura o ponteiro de base\n", "    pop rbp"))
	output.WriteString(fmt.Sprintf("%-40s; Status de saída zero (Sucesso)\n", "    mov rax, 0"))
	output.WriteString(fmt.Sprintf("%-40s; Retorna ao Sistema Operacional\n", "    ret"))

	return output.String()
}