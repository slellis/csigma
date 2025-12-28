package codegen

import (
	"csigma/parser"
	"fmt"
)

func GenerateNASM(statements []parser.Statement) string {
	// --- SEÇÃO DE DADOS ---
	dataSection := "section .data ; Area para variaveis e constantes inicializadas\n"
	dataSection += "    fmt_in db ' %ld', 0 ; Formato para leitura de inteiros (scanf)\n"
	dataSection += "    fmt_out_num db '%ld', 10, 0 ; Formato para exibir numeros com quebra de linha (printf)\n"

	// --- SEÇÃO DE CÓDIGO ---
	textSection := "\nsection .text ; Area com as instrucoes executaveis\n"
	textSection += "extern printf, scanf ; Declara funcoes externas da LibC\n"
	textSection += "global main ; Exporta o ponto de entrada para o Linker\n\n"
	textSection += "main:\n"
	textSection += "    push rbp ; Salva o Base Pointer antigo na pilha\n"
	textSection += "    mov rbp, rsp ; Define o novo Base Pointer como o topo atual da pilha\n"
	textSection += "    sub rsp, 32 ; Alinha a pilha em 16 bytes e reserva espaco de rascunho\n\n"

	msgCount := 0

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *parser.VarDeclNode:
			// Reserva 8 bytes (dq) para a variavel com seu valor inicial
			dataSection += fmt.Sprintf("    %s dq %s ; Alocacao da variavel %s\n", s.Name, s.Value, s.Name)

		case *parser.PrintNode:
			if s.IsString {
				msgName := fmt.Sprintf("msg_%d", msgCount)
				dataSection += fmt.Sprintf("    %s db '%s', 10, 0 ; Constante de texto\n", msgName, s.Value)
				
				textSection += fmt.Sprintf("    lea rdi, [%s] ; Carrega o endereco da string em RDI (1o arg)\n", msgName)
				textSection += "    xor eax, eax ; Indica zero argumentos de ponto flutuante\n"
				textSection += "    call printf ; Chama a funcao de impressao do sistema\n"
				msgCount++
			} else {
				textSection += "    lea rdi, [fmt_out_num] ; Carrega o formato de numero em RDI\n"
				textSection += fmt.Sprintf("    mov rsi, [%s] ; Move o VALOR da variavel %s para RSI (2o arg)\n", s.Value, s.Value)
				textSection += "    xor eax, eax ; Limpa registradores de retorno/ponto flutuante\n"
				textSection += "    call printf ; Exibe o valor numerico na tela\n"
			}

		case *parser.InputNode:
			textSection += "    lea rdi, [fmt_in] ; Carrega o formato de entrada em RDI\n"
			textSection += fmt.Sprintf("    lea rsi, [%s] ; Carrega o ENDERECO de %s em RSI para o scanf salvar\n", s.VarName, s.VarName)
			textSection += "    xor eax, eax ; Prepara chamada de sistema\n"
			textSection += "    call scanf ; Aguarda a digitacao do usuario\n"

		case *parser.AssignmentNode:
			textSection += fmt.Sprintf("    mov rax, [%s] ; Carrega o valor de %s no acumulador RAX\n", s.Left, s.Left)
			textSection += fmt.Sprintf("    add rax, [%s] ; Soma o valor de %s ao acumulador\n", s.Right, s.Right)
			textSection += fmt.Sprintf("    mov [%s], rax ; Salva o resultado final no endereco de %s\n", s.Dest, s.Dest)
		}
	}

	// --- EPÍLOGO ---
	textSection += "\n    add rsp, 32 ; Restaura o ponteiro da pilha (limpa o rascunho)\n"
	textSection += "    pop rbp ; Recupera o Base Pointer original\n"
	textSection += "    mov rax, 0 ; Define o codigo de saida do programa como 0 (Sucesso)\n"
	textSection += "    ret ; Retorna o controle para o Sistema Operacional\n"

	return dataSection + textSection
}