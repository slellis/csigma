section .data                           ; --- SECAO DE DADOS (WORKING-STORAGE) ---
    fmt_in db ' %ld', 0                 ; Formato scanf
    fmt_out_num db '%ld', 10, 0         ; Formato printf
    A                    dq 0           ; Alocacao da variavel A
    B                    dq 0           ; Alocacao da variavel B
    C                    dq 0           ; Alocacao da variavel C
    RESULTADO            dq 0           ; Alocacao da variavel RESULTADO
    msg_0                db 'VALOR A:', 10, 0               ; Texto para exibicao
    msg_1                db 'VALOR B:', 10, 0               ; Texto para exibicao
    msg_2                db 'VALOR C:', 10, 0               ; Texto para exibicao
    msg_3                db 'RESULTADO FINAL:', 10, 0               ; Texto para exibicao

section .text                           ; --- SECAO DE CODIGO (PROCEDURE DIVISION) ---
extern printf, scanf                    ; Funcoes da biblioteca C padrao
global main                             ; Ponto de entrada do executavel

main:
    push rbp                            ; Salva o ponteiro de base da pilha
    mov rbp, rsp                        ; Alinha o ponteiro de base
    sub rsp, 32                         ; Reserva espaco e alinha stack em 16 bytes

    lea rdi, [msg_0]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [A]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    lea rdi, [msg_1]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [B]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    lea rdi, [msg_2]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [C]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    mov rax, [A]                        ; Carrega conteudo de A em RAX
    add rax, [B]                        ; Soma B
    imul rax, 2                         ; Multiplica 2
    mov rbx, [C]                        ; Carrega divisor em RBX
    xor rdx, rdx                        ; Limpa RDX para divisao segura
    cqo                                 ; Estende sinal de RAX para RDX:RAX
    idiv rbx                            ; Divide RDX:RAX por RBX (Quociente->RAX)
    mov [RESULTADO], rax                ; Armazena resultado em RESULTADO
    lea rdi, [msg_3]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_out_num]              ; Formato de saida em RDI
    mov rsi, [RESULTADO]                ; Valor da variavel em RSI
    xor eax, eax                        ; Limpa EAX para printf
    call printf                         ; Exibe o valor numerico

    add rsp, 32                         ; Libera espaco da pilha
    pop rbp                             ; Restaura o ponteiro de base
    mov rax, 0                          ; Status de saída zero (Sucesso)
    ret                                 ; Retorna ao Sistema Operacional
