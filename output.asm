section .data                           ; Area para variaveis e constantes inicializadas
    fmt_in db ' %ld', 0                 ; Formato para leitura de inteiros (scanf)
    fmt_out_num db '%ld', 10, 0         ; Formato para exibir numeros com quebra de linha (printf)
    A dq 0                              ; Alocacao da variavel A
    B dq 0                              ; Alocacao da variavel B
    C dq 0                              ; Alocacao da variavel C
    msg_0 db 'DIGITE O VALOR DE A:', 10, 0; Constante de texto
    msg_1 db 'DIGITE O VALOR DE B:', 10, 0; Constante de texto
    msg_2 db 'O RESULTADO DE A - B EH:', 10, 0; Constante de texto

section .text                           ; Area com as instrucoes executaveis
extern printf, scanf                    ; Declara funcoes externas da LibC
global main                             ; Exporta o ponto de entrada para o Linker

main:
    push rbp                            ; Salva o Base Pointer antigo na pilha
    mov rbp, rsp                        ; Define o novo Base Pointer como o topo atual
    sub rsp, 32                         ; Alinha a pilha em 16 bytes e reserva espaco

    lea rdi, [msg_0]                    ; Endereco da string em RDI
    xor eax, eax                        ; Zero args de ponto flutuante
    call printf                         ; Chama printf
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [A]                        ; Endereco de A em RSI
    xor eax, eax                        ; Prepara chamada
    call scanf                          ; Aguarda digitacao
    lea rdi, [msg_1]                    ; Endereco da string em RDI
    xor eax, eax                        ; Zero args de ponto flutuante
    call printf                         ; Chama printf
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [B]                        ; Endereco de B em RSI
    xor eax, eax                        ; Prepara chamada
    call scanf                          ; Aguarda digitacao
    mov rax, [A]                        ; A -> RAX
    sub rax, [B]                        ; RAX - B
    mov [C], rax                        ; RAX -> C
    lea rdi, [msg_2]                    ; Endereco da string em RDI
    xor eax, eax                        ; Zero args de ponto flutuante
    call printf                         ; Chama printf
    lea rdi, [fmt_out_num]              ; Formato de numero em RDI
    mov rsi, [C]                        ; Valor de C em RSI
    xor eax, eax                        ; Limpa regs
    call printf                         ; Exibe valor numerico

    add rsp, 32                         ; Restaura a pilha
    pop rbp                             ; Recupera Base Pointer
    mov rax, 0                          ; Retorno 0 (Sucesso)
    ret                                 ; Volta para o SO
