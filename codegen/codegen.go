package codegen

import (
	"csigma/lexer"
	"csigma/parser"
	"fmt"
)

// GenerateNASM traduz a nossa AST (Árvore Sintática) para Assembly x86_64.
// Ele orquestra a memória e as instruções de CPU necessárias para o programa rodar.
func GenerateNASM(statements []parser.Statement) string {
	
	// --- SEÇÃO DE DADOS (.data) ---
	// Aqui reservamos espaço fixo na memória, equivalente à WORKING-STORAGE do COBOL.
	dataSection := "section .data                           ; Area para variaveis e constantes inicializadas\n"
	
	// 'fmt_in' é o molde que o scanf usa para capturar números do teclado.
	dataSection += "    fmt_in db ' %ld', 0                 ; Formato para leitura de inteiros (scanf)\n"
	
	// 'fmt_out_num' é o molde do printf para exibir números seguidos de nova linha (10).
	dataSection += "    fmt_out_num db '%ld', 10, 0         ; Formato para exibir numeros com quebra de linha (printf)\n"

	// --- SEÇÃO DE CÓDIGO (.text) ---
	// Aqui ficam as instruções executáveis que o processador irá processar.
	textSection := "\nsection .text                           ; Area com as instrucoes executaveis\n"
	
	// 'extern' avisa ao Linker que printf/scanf estão nas bibliotecas do Linux (LibC).
	textSection += "extern printf, scanf                    ; Declara funcoes externas da LibC\n"
	textSection += "global main                             ; Exporta o ponto de entrada para o Linker\n\n"
	
	// Início do programa principal.
	textSection += "main:\n"
	
	// PRÓLOGO: Prepara a Stack (Pilha) salvando a base anterior (RBP).
	textSection += "    push rbp                            ; Salva o Base Pointer antigo na pilha\n"
	textSection += "    mov rbp, rsp                        ; Define o novo Base Pointer como o topo atual\n"
	
	// ALINHAMENTO: Subtraímos 32 bytes para garantir que a pilha esteja 
	// alinhada em 16 bytes, exigência do Linux para chamar funções da LibC.
	textSection += "    sub rsp, 32                         ; Alinha a pilha em 16 bytes e reserva espaco\n\n"

	msgCount := 0

	// LOOP DE TRADUÇÃO: Percorre cada nó gerado pelo Parser.
	for _, stmt := range statements {
		switch s := stmt.(type) {
		
		case *parser.VarDeclNode:
			// VAR A = 0: Reserva 8 bytes (dq = define quadword) na memória.
			line := fmt.Sprintf("    %s dq %s", s.Name, s.Value)
			dataSection += fmt.Sprintf("%-40s; Alocacao da variavel %s\n", line, s.Name)

		case *parser.PrintNode:
			if s.IsString {
				// PRINT "TEXTO": Cria uma constante de string terminada em nulo (0).
				msgName := fmt.Sprintf("msg_%d", msgCount)
				lineMsg := fmt.Sprintf("    %s db '%s', 10, 0", msgName, s.Value)
				dataSection += fmt.Sprintf("%-40s; Constante de texto\n", lineMsg)
				
				// LEA (Load Effective Address): Coloca o endereço da string no RDI (1º arg do printf).
				textSection += fmt.Sprintf("%-40s; Endereco da string em RDI\n", "    lea rdi, ["+msgName+"]")
				textSection += fmt.Sprintf("%-40s; Zero args de ponto flutuante\n", "    xor eax, eax")
				textSection += fmt.Sprintf("%-40s; Chama printf\n", "    call printf")
				msgCount++
			} else {
				// PRINT VAR: Passa o formato em RDI e o valor da variável em RSI (2º arg).
				textSection += fmt.Sprintf("%-40s; Formato de numero em RDI\n", "    lea rdi, [fmt_out_num]")
				textSection += fmt.Sprintf("%-40s; Valor de %s em RSI\n", "    mov rsi, ["+s.Value+"]", s.Value)
				textSection += fmt.Sprintf("%-40s; Limpa regs\n", "    xor eax, eax")
				textSection += fmt.Sprintf("%-40s; Exibe valor numerico\n", "    call printf")
			}

		case *parser.InputNode:
			// INPUT VAR: O scanf precisa do formato (RDI) e do endereço da variável (RSI).
			textSection += fmt.Sprintf("%-40s; Formato de entrada em RDI\n", "    lea rdi, [fmt_in]")
			textSection += fmt.Sprintf("%-40s; Endereco de %s em RSI\n", "    lea rsi, ["+s.VarName+"]", s.VarName)
			textSection += fmt.Sprintf("%-40s; Prepara chamada\n", "    xor eax, eax")
			textSection += fmt.Sprintf("%-40s; Aguarda digitacao\n", "    call scanf")

		case *parser.AssignmentNode:
			// LÓGICA ARITMÉTICA: Uso do registrador acumulador RAX.
			textSection += fmt.Sprintf("%-40s; %s -> RAX\n", "    mov rax, ["+s.Left+"]", s.Left)
			
			// Decisão do Operador: Se for TokenPlus faz ADD, senão faz SUB (Subtração)
			if s.Operator == lexer.TokenPlus {
				textSection += fmt.Sprintf("%-40s; RAX + %s\n", "    add rax, ["+s.Right+"]", s.Right)
			} else {
				textSection += fmt.Sprintf("%-40s; RAX - %s\n", "    sub rax, ["+s.Right+"]", s.Right)
			}
			
			textSection += fmt.Sprintf("%-40s; RAX -> %s\n", "    mov ["+s.Dest+"], rax", s.Dest)
		}
	}

	// --- EPÍLOGO ---
	// Devolvemos o espaço da pilha e restauramos o RBP original.
	textSection += "\n"
	textSection += fmt.Sprintf("%-40s; Restaura a pilha\n", "    add rsp, 32")
	textSection += fmt.Sprintf("%-40s; Recupera Base Pointer\n", "    pop rbp")
	
	// RAX = 0 indica sucesso ao Sistema Operacional.
	textSection += fmt.Sprintf("%-40s; Retorno 0 (Sucesso)\n", "    mov rax, 0")
	textSection += fmt.Sprintf("%-40s; Volta para o SO\n", "    ret")

	return dataSection + textSection
}