package codegen

import (
	"csigma/parser"
	"fmt"
	"strings"
)

func GenerateNASM(statements []parser.Statement) string {
	var dataSection strings.Builder
	var textSection strings.Builder
	msgCount := 0

	dataSection.WriteString("section .data\n")
	dataSection.WriteString(fmt.Sprintf("%-40s; Scanf\n", "    fmt_in db '%ld', 0"))
	dataSection.WriteString(fmt.Sprintf("%-40s; Printf\n", "    fmt_out_num db '%ld', 10, 0"))

	textSection.WriteString("\nsection .text\nextern printf, scanf\nglobal main\n\nmain:\n")
	textSection.WriteString("    push rbp\n    mov rbp, rsp\n    sub rsp, 32\n\n")

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *parser.VarDeclNode:
			line := fmt.Sprintf("    %-20s dq %s", s.Name, s.Value)
			dataSection.WriteString(fmt.Sprintf("%-40s; Variavel %s\n", line, s.Name))

		case *parser.AssignmentNode:
			textSection.WriteString(fmt.Sprintf("\n    ; --- Calculo de %s ---\n", s.Dest))
			textSection.WriteString(fmt.Sprintf("    mov rax, [%s]\n", s.First))
			for _, op := range s.Ops {
				val := op.Value
				if op.IsVar {
					val = "[" + op.Value + "]"
				}

				switch op.Operator {
				case "+":
					textSection.WriteString(fmt.Sprintf("    add rax, %-25s\n", val))
				case "-":
					textSection.WriteString(fmt.Sprintf("    sub rax, %-25s\n", val))
				case "*":
					textSection.WriteString(fmt.Sprintf("    imul rax, %-24s\n", val))
				case "/":
					textSection.WriteString(fmt.Sprintf("    mov rbx, %-25s\n", val))
					textSection.WriteString("    xor rdx, rdx\n    idiv rbx\n")
				}
			}
			textSection.WriteString(fmt.Sprintf("    mov [%s], rax\n", s.Dest))

		case *parser.PrintNode:
			if s.IsString {
				msgName := fmt.Sprintf("msg_%d", msgCount)
				dataSection.WriteString(fmt.Sprintf("    %-20s db '%s', 10, 0\n", msgName, s.Value))
				textSection.WriteString(fmt.Sprintf("    lea rdi, [%s]\n    xor eax, eax\n    call printf\n", msgName))
				msgCount++
			} else {
				textSection.WriteString(fmt.Sprintf("    lea rdi, [fmt_out_num]\n    mov rsi, [%s]\n    xor eax, eax\n    call printf\n", s.Value))
			}

		case *parser.InputNode:
			textSection.WriteString(fmt.Sprintf("    lea rdi, [fmt_in]\n    lea rsi, [%s]\n    xor eax, eax\n    call scanf\n", s.VarName))
		}
	}

	textSection.WriteString("\n    add rsp, 32\n    pop rbp\n    mov rax, 0\n    ret\n")
	return dataSection.String() + textSection.String()
}
