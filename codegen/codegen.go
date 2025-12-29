package codegen

import (
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"strings"
	"unicode"
)

// isNumber checks if a string represents a pure numeric value
func isNumber(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) { return false }
	}
	return true
}

// GenerateNASM gera o arquivo .asm com seções e comentários didáticos.
func GenerateNASM(statements []parser.Statement) string {
	variables := make(map[string]string)
	var stringsConst []string
	var instructions []string
	msgCount := 0

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *parser.VarDeclNode:
			variables[s.Name] = s.Value

		case *parser.PrintNode:
			if s.IsString {
				msgName := fmt.Sprintf("msg_%d", msgCount)
				lineMsg := fmt.Sprintf("    %-20s db '%s', 10, 0", msgName, s.Value)
				stringsConst = append(stringsConst, lineMsg)
				instructions = append(instructions, fmt.Sprintf("%-40s; Endereco da string para printf", "    lea rdi, ["+msgName+"]"))
				instructions = append(instructions, "    xor eax, eax\n    call printf")
				msgCount++
			} else {
				instructions = append(instructions, "    lea rdi, [fmt_out_num]")
				instructions = append(instructions, fmt.Sprintf("%-40s; Valor de %s em RSI", "    mov rsi, ["+s.Value+"]", s.Value))
				instructions = append(instructions, "    xor eax, eax\n    call printf")
			}

		case *parser.InputNode:
			instructions = append(instructions, "    lea rdi, [fmt_in]")
			instructions = append(instructions, fmt.Sprintf("%-40s; Endereco de %s em RSI", "    lea rsi, ["+s.VarName+"]", s.VarName))
			instructions = append(instructions, "    xor eax, eax\n    call scanf")

		case *parser.AssignmentNode:
			// Carga inicial
			if isNumber(s.First) {
				instructions = append(instructions, fmt.Sprintf("%-40s; Carrega numero %s", "    mov rax, "+s.First, s.First))
			} else {
				instructions = append(instructions, fmt.Sprintf("%-40s; Carrega variavel %s", "    mov rax, ["+s.First+"]", s.First))
			}
			
			for _, op := range s.Rest {
				instr := ""
				label := ""
				val := ""

				// Se for número, usa direto. Se for variável, usa [nome]
				if isNumber(op.Value) {
					val = op.Value
				} else {
					val = "[" + op.Value + "]"
				}

				switch op.Operator {
				case lexer.TokenPlus:
					instr = "add"; label = "Soma"
				case lexer.TokenMinus:
					instr = "sub"; label = "Subtrai"
				case lexer.TokenMult:
					instr = "imul"; label = "Multiplica"
				case lexer.TokenDiv:
					// Divisão exige registrador para o divisor
					instructions = append(instructions, fmt.Sprintf("%-40s; Move divisor para RBX", "    mov rbx, "+val))
					instructions = append(instructions, fmt.Sprintf("%-40s; Estende sinal p/ RDX", "    cqo"))
					instructions = append(instructions, fmt.Sprintf("%-40s; Divide RAX por RBX", "    idiv rbx"))
					continue // idiv já faz o trabalho
				}

				if instr != "" {
					line := fmt.Sprintf("    %s rax, %s", instr, val)
					instructions = append(instructions, fmt.Sprintf("%-40s; %s", line, label))
				}
			}
			instructions = append(instructions, fmt.Sprintf("%-40s; Salva em %s", "    mov ["+s.Dest+"], rax", s.Dest))
		}
	}

	var output strings.Builder
	output.WriteString("section .data                           ; Area de dados (WORKING-STORAGE)\n")
	output.WriteString("    fmt_in db ' %ld', 0                 ; Formato para entrada numerica\n")
	output.WriteString("    fmt_out_num db '%ld', 10, 0         ; Formato para saida numerica\n")
	
	for name, val := range variables {
		line := fmt.Sprintf("    %-20s dq %s", name, val)
		output.WriteString(fmt.Sprintf("%-40s; Variavel %s\n", line, name))
	}
	for _, s := range stringsConst {
		output.WriteString(s + "               ; Constante de texto\n")
	}

	output.WriteString("\nsection .text                           ; Area de codigo (PROCEDURE DIVISION)\n")
	output.WriteString("extern printf, scanf\nglobal main\n\nmain:\n")
	output.WriteString("    push rbp                            ; Prologo\n")
	output.WriteString("    mov rbp, rsp\n")
	output.WriteString("    sub rsp, 32                         ; Alinhamento de pilha\n\n")

	for _, ins := range instructions {
		output.WriteString(ins + "\n")
	}

	output.WriteString("\n    add rsp, 32                         ; Epilogo\n")
	output.WriteString("    pop rbp\n")
	output.WriteString("    mov rax, 0\n")
	output.WriteString("    ret\n")

	return output.String()
}