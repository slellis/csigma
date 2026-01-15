package codegen

import (
	"csigma/parser"
	"fmt"
	"strings"
)

// GenerateNASM: O tradutor final que converte a AST em código de montagem (Assembly).
func GenerateNASM(statements []parser.Statement) string {
	var dataSection strings.Builder
	var textSection strings.Builder
	msgCount := 0

	// --- SEÇÃO DE DADOS (.data) ---
	// Reservada para constantes e variáveis globais.
	dataSection.WriteString("section .data\n")
	// db = Define Byte. Usado para strings e formatos de E/S.
	dataSection.WriteString(fmt.Sprintf("%-40s; Formato para leitura (long int)\n", "    fmt_in db '%ld', 0"))
	dataSection.WriteString(fmt.Sprintf("%-40s; Formato para escrita (long int + \n)\n", "    fmt_out_num db '%ld', 10, 0"))

	// --- SEÇÃO DE CÓDIGO (.text) ---
	textSection.WriteString("\nsection .text\n")
	textSection.WriteString("extern printf, scanf                    ; Declara funções da LibC\n")
	textSection.WriteString("global main                             ; Ponto de entrada para o Linker\n\nmain:\n")

	// Prólogo da Função: Prepara a base da pilha (Stack Frame)
	textSection.WriteString("    push rbp                            ; Salva o ponteiro da base da pilha anterior\n")
	textSection.WriteString("    mov rbp, rsp                        ; Define a nova base da pilha\n")
	textSection.WriteString("    sub rsp, 32                         ; Alinha a pilha (16-byte alignment)\n\n")

	// --- PROCESSAMENTO DA AST ---
	for _, stmt := range statements {
		switch s := stmt.(type) {

		case *parser.VarDeclNode:
			// dq = Define Quadword (64 bits). Reserva espaço para inteiros Sigma.
			line := fmt.Sprintf("    %-20s dq %s", s.Name, s.Value)
			dataSection.WriteString(fmt.Sprintf("%-40s; Reserva memoria para %s\n", line, s.Name))

		case *parser.AssignmentNode:
			textSection.WriteString(fmt.Sprintf("\n    ; --- Calculo Aritmetico: %s ---\n", s.Dest))
			// RAX é o acumulador principal para operações matemáticas.
			textSection.WriteString(fmt.Sprintf("    mov rax, [%s]          ; Carrega o primeiro valor em RAX\n", s.First))

			for _, op := range s.Ops {
				val := op.Value
				if op.IsVar {
					val = "[" + op.Value + "]" // Acessa o valor no endereço da memória
				}

				switch op.Operator {
				case "+":
					textSection.WriteString(fmt.Sprintf("    add rax, %-25s ; Soma ao acumulador\n", val))
				case "-":
					textSection.WriteString(fmt.Sprintf("    sub rax, %-25s ; Subtrai do acumulador\n", val))
				case "*":
					// imul: Multiplicação com sinal de 64 bits.
					textSection.WriteString(fmt.Sprintf("    imul rax, %-24s ; Multiplicacao\n", val))
				case "/":
					// idiv exige que o dividendo esteja em RDX:RAX.
					textSection.WriteString(fmt.Sprintf("    mov rbx, %-25s ; Carrega o divisor em RBX\n", val))
					textSection.WriteString("    xor rdx, rdx               ; Zera RDX para evitar 'overflow' na divisao\n")
					textSection.WriteString("    idiv rbx                   ; Divide RAX por RBX (Resultado em RAX)\n")
				}
			}
			// Move o resultado final do acumulador RAX para o endereço de destino na memória.
			textSection.WriteString(fmt.Sprintf("    mov [%s], rax          ; Armazena resultado final\n", s.Dest))

		case *parser.PrintNode:
			if s.IsString {
				msgName := fmt.Sprintf("msg_%d", msgCount)
				// 10 = Newline (\n), 0 = Null Terminator (Padrão C)
				dataSection.WriteString(fmt.Sprintf("    %-20s db '%s', 10, 0\n", msgName, s.Value))

				// lea: Load Effective Address. Passa o endereço da string para RDI.
				textSection.WriteString(fmt.Sprintf("    lea rdi, [%s]           ; RDI = Primeiro argumento (string)\n", msgName))
				textSection.WriteString("    xor eax, eax                        ; AL=0 indica que não há vetores SSE\n")
				textSection.WriteString("    call printf\n")
				msgCount++
			} else {
				textSection.WriteString("    lea rdi, [fmt_out_num]              ; RDI = Formato de saida\n")
				textSection.WriteString(fmt.Sprintf("    mov rsi, [%s]           ; RSI = Segundo argumento (valor)\n", s.Value))
				textSection.WriteString("    xor eax, eax\n")
				textSection.WriteString("    call printf\n")
			}

		case *parser.InputNode:
			textSection.WriteString("    lea rdi, [fmt_in]                       ; RDI = Formato de entrada\n")
			textSection.WriteString(fmt.Sprintf("    lea rsi, [%s]               ; RSI = Endereco onde salvar\n", s.VarName))
			textSection.WriteString("    xor eax, eax\n")
			textSection.WriteString("    call scanf\n")
		}
	}

	// Epílogo da Função: Limpa a pilha e retorna ao SO.
	textSection.WriteString("\n    add rsp, 32                       ; Restaura espaco da pilha\n")
	textSection.WriteString("    pop rbp                             ; Restaura o RBP original\n")
	textSection.WriteString("    mov rax, 0                          ; Return 0\n")
	textSection.WriteString("    ret\n")

	return dataSection.String() + textSection.String()
}
