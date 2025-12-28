section .data ; Area para variaveis e constantes inicializadas
    fmt_in db ' %ld', 0 ; Formato para leitura de inteiros (scanf)
    fmt_out_num db '%ld', 10, 0 ; Formato para exibir numeros com quebra de linha (printf)
    A dq 0 ; Alocacao da variavel A
    B dq 0 ; Alocacao da variavel B
    C dq 0 ; Alocacao da variavel C
    msg_0 db 'SOMA', 10, 0 ; Constante de texto
    msg_1 db 'Informe o valor de A', 10, 0 ; Constante de texto
    msg_2 db 'Informe o valor de B', 10, 0 ; Constante de texto
    msg_3 db 'Total', 10, 0 ; Constante de texto

section .text ; Area com as instrucoes executaveis
extern printf, scanf ; Declara funcoes externas da LibC
global main ; Exporta o ponto de entrada para o Linker

main:
    push rbp ; Salva o Base Pointer antigo na pilha
    mov rbp, rsp ; Define o novo Base Pointer como o topo atual da pilha
    sub rsp, 32 ; Alinha a pilha em 16 bytes e reserva espaco de rascunho

    lea rdi, [msg_0] ; Carrega o endereco da string em RDI (1o arg)
    xor eax, eax ; Indica zero argumentos de ponto flutuante
    call printf ; Chama a funcao de impressao do sistema
    lea rdi, [msg_1] ; Carrega o endereco da string em RDI (1o arg)
    xor eax, eax ; Indica zero argumentos de ponto flutuante
    call printf ; Chama a funcao de impressao do sistema
    lea rdi, [fmt_in] ; Carrega o formato de entrada em RDI
    lea rsi, [A] ; Carrega o ENDERECO de A em RSI para o scanf salvar
    xor eax, eax ; Prepara chamada de sistema
    call scanf ; Aguarda a digitacao do usuario
    lea rdi, [msg_2] ; Carrega o endereco da string em RDI (1o arg)
    xor eax, eax ; Indica zero argumentos de ponto flutuante
    call printf ; Chama a funcao de impressao do sistema
    lea rdi, [fmt_in] ; Carrega o formato de entrada em RDI
    lea rsi, [B] ; Carrega o ENDERECO de B em RSI para o scanf salvar
    xor eax, eax ; Prepara chamada de sistema
    call scanf ; Aguarda a digitacao do usuario
    mov rax, [A] ; Carrega o valor de A no acumulador RAX
    add rax, [B] ; Soma o valor de B ao acumulador
    mov [C], rax ; Salva o resultado final no endereco de C
    lea rdi, [msg_3] ; Carrega o endereco da string em RDI (1o arg)
    xor eax, eax ; Indica zero argumentos de ponto flutuante
    call printf ; Chama a funcao de impressao do sistema
    lea rdi, [fmt_out_num] ; Carrega o formato de numero em RDI
    mov rsi, [C] ; Move o VALOR da variavel C para RSI (2o arg)
    xor eax, eax ; Limpa registradores de retorno/ponto flutuante
    call printf ; Exibe o valor numerico na tela

    add rsp, 32 ; Restaura o ponteiro da pilha (limpa o rascunho)
    pop rbp ; Recupera o Base Pointer original
    mov rax, 0 ; Define o codigo de saida do programa como 0 (Sucesso)
    ret ; Retorna o controle para o Sistema Operacional
