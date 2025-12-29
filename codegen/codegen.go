package codegen

import (
	"csigma/lexer"
	"csigma/parser"
	"fmt"
	"strings"
)

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
				stringsConst = append(stringsConst, fmt.Sprintf("    %-20s db '%s', 10, 0", msgName, s.Value))
				instructions = append(instructions, fmt.Sprintf("    lea rdi, [%s]\n    xor eax, eax\n    call printf", msgName))
				msgCount++
			} else {
				instructions = append(instructions, fmt.Sprintf("    lea rdi, [fmt_out_num]\n    mov rsi, [%s]\n    xor eax, eax\n    call printf", s.Value))
			}
		case *parser.InputNode:
			instructions = append(instructions, fmt.Sprintf("    lea rdi, [fmt_in]\n    lea rsi, [%s]\n    xor eax, eax\n    call scanf", s.VarName))
		case *parser.AssignmentNode:
			instructions = append(instructions, fmt.Sprintf("    mov rax, [%s]", s.First))
			for _, op := range s.Rest {
				instr := "add"
				if op.Operator == lexer.TokenMinus { instr = "sub" }
				instructions = append(instructions, fmt.Sprintf("    %s rax, [%s]", instr, op.Value))
			}
			instructions = append(instructions, fmt.Sprintf("    mov [%s], rax", s.Dest))
		}
	}

	var output strings.Builder
	output.WriteString("section .data\n    fmt_in db ' %ld', 0\n    fmt_out_num db '%ld', 10, 0\n")
	for name, val := range variables {
		output.WriteString(fmt.Sprintf("    %-20s dq %s\n", name, val))
	}
	for _, s := range stringsConst {
		output.WriteString(s + "\n")
	}

	output.WriteString("\nsection .text\nextern printf, scanf\nglobal main\n\nmain:\n")
	output.WriteString("    push rbp\n    mov rbp, rsp\n    sub rsp, 32\n\n")

	for _, ins := range instructions {
		output.WriteString(ins + "\n")
	}

	output.WriteString("\n    add rsp, 32\n    pop rbp\n    mov rax, 0\n    ret\n")
	return output.String()
}