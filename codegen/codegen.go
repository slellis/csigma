package codegen

import (
	"csigma/parser"
	"fmt"
)

// GenerateNASM é o nosso "Linkage Editor" lógico. Ele traduz a AST (Árvore Sintática)
// em um código Assembly x86_64 compatível com a ABI do Linux.
func GenerateNASM(statements []parser.Statement) string {
	
	// --- SEÇÃO DE DADOS (.data) ---
	// Equivale à WORKING-STORAGE do COBOL. Aqui reservamos espaço fixo na memória.
	dataSection := "section .data\n"
	
	// 'fmt_in': Formato para o scanf. %ld lê inteiros de 64 bits.
	dataSection += "    fmt_in db ' %ld', 0\n" 
	
	// 'fmt_out_num': Formato para printf exibir números. O 10 é o caractere \n (Line Feed).
	dataSection += "    fmt_out_num db '%ld', 10, 0\n"

	// --- SEÇÃO DE CÓDIGO (.text) ---
	// Aqui ficam as instruções que o processador irá executar sequencialmente.
	textSection := "\nsection .text\n"
	
	// Importamos funções da LibC (Biblioteca padrão do C no Linux).
	// O GCC fará o 'bind' desses nomes com os endereços reais na etapa de Link-edit.
	textSection += "extern printf, scanf\n"
	
	// Define o ponto de entrada global para o Linker.
	textSection += "global main\n\n"
	
	// Início do procedimento principal.
	textSection += "main:\n"
	
	// PROLOGO: Prepara a Stack (Pilha).
	// 'push rbp' e 'mov rbp, rsp' criam o Stack Frame para que possamos
	// retornar ao Sistema Operacional corretamente ao final.
	textSection += "    push rbp\n"
	textSection += "    mov rbp, rsp\n"
	
	// ALINHAMENTO DE PILHA: O Linux x86_64 exige que a pilha esteja 
	// alinhada em 16 bytes antes de chamadas de funções externas (ABI).
	// Subtraímos 32 bytes para garantir esse espaço e segurança.
	textSection += "    sub rsp, 32\n\n"

	msgCount := 0

	// LOOP DE TRADUÇÃO: Percorre cada nó da Árvore Sintática (AST).
	for _, stmt := range statements {
		switch s := stmt.(type) {
		
		case *parser.VarDeclNode:
			// VAR A = 0 -> Define na memória 8 bytes (dq = define quadword)
			dataSection += fmt.Sprintf("    %s dq %s\n", s.Name, s.Value)

		case *parser.PrintNode:
			if s.IsString {
				// PRINT "TEXTO": Gera um rótulo de mensagem único na seção .data
				msgName := fmt.Sprintf("msg_%d", msgCount)
				dataSection += fmt.Sprintf("    %s db '%s', 10, 0\n", msgName, s.Value)
				
				// LEA (Load Effective Address): Carrega o ENDEREÇO da string no RDI.
				// RDI é o registrador padrão para o 1º argumento no Linux.
				textSection += fmt.Sprintf("    lea rdi, [%s]\n    xor eax, eax\n    call printf\n", msgName)
				msgCount++
			} else {
				// PRINT VAR: Passa o formato em RDI e o VALOR da variável em RSI.
				// RSI é o registrador padrão para o 2º argumento.
				textSection += fmt.Sprintf("    lea rdi, [fmt_out_num]\n    mov rsi, [%s]\n    xor eax, eax\n    call printf\n", s.Value)
			}

		case *parser.InputNode:
			// INPUT VAR: Passa o formato em RDI e o ENDEREÇO de destino em RSI.
			// O scanf precisa saber ONDE salvar o dado digitado.
			textSection += fmt.Sprintf("    lea rdi, [fmt_in]\n    lea rsi, [%s]\n    xor eax, eax\n    call scanf\n", s.VarName)

		case *parser.AssignmentNode:
			// C = A + B: A lógica aritmética pura.
			// 1. Movemos o VALOR de A para o registrador acumulador RAX.
			textSection += fmt.Sprintf("    mov rax, [%s]\n", s.Left)
			// 2. Somamos o VALOR de B ao que está em RAX.
			textSection += fmt.Sprintf("    add rax, [%s]\n", s.Right)
			// 3. Movemos o resultado final de RAX para o espaço de memória de C.
			textSection += fmt.Sprintf("    mov [%s], rax\n", s.Dest)
		}
	}

	// EPÍLOGO: Limpeza antes de sair.
	// Restauramos o ponteiro da pilha (rsp) somando o que subtraímos no início.
	textSection += "\n    add rsp, 32\n"
	textSection += "    pop rbp\n"
	
	// Retornamos 0 em RAX, indicando ao Linux que o programa terminou com SUCESSO.
	textSection += "    mov rax, 0\n"
	textSection += "    ret\n"

	return dataSection + textSection
}