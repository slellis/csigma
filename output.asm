section .data                           ; --- SECAO DE DADOS (WORKING-STORAGE) ---
    fmt_in db ' %ld', 0                 ; Formato scanf
    fmt_out_num db '%ld', 10, 0         ; Formato printf
    b                    dq 0           ; Alocacao da variavel b
    c                    dq 0           ; Alocacao da variavel c
    res                  dq 0           ; Alocacao da variavel res
    a                    dq 0           ; Alocacao da variavel a
    msg_0                db 'Calculadora', 10, 0               ; Texto para exibicao
    msg_1                db 'Conta mista: (a + b) * 2 / C', 10, 0               ; Texto para exibicao
    msg_2                db 'Valor de a:', 10, 0               ; Texto para exibicao
    msg_3                db 'Valor de b:', 10, 0               ; Texto para exibicao
    msg_4                db 'Valor de c:', 10, 0               ; Texto para exibicao
    msg_5                db 'Resultado final:', 10, 0               ; Texto para exibicao

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
    lea rdi, [msg_1]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [msg_2]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [a]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    lea rdi, [msg_3]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [b]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    lea rdi, [msg_4]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_in]                   ; Formato de entrada em RDI
    lea rsi, [c]                        ; Endereco de destino em RSI
    xor eax, eax                        ; Limpa EAX para scanf
    call scanf                          ; Captura entrada do teclado
    mov rax, [a]                        ; Carrega conteudo de a em RAX
    add rax, [b]                        ; Soma b
    imul rax, 2                         ; Multiplica 2
    mov rbx, [c]                        ; Carrega divisor em RBX
    xor rdx, rdx                        ; Limpa RDX para divisao segura
    cqo                                 ; Estende sinal de RAX para RDX:RAX
    idiv rbx                            ; Divide RDX:RAX por RBX (Quociente->RAX)
    mov [res], rax                      ; Armazena resultado em res
    lea rdi, [msg_5]                    ; Endereco da string para RDI
    xor eax, eax                        ; Limpa EAX (sem args de ponto flutuante)
    call printf                         ; Chama a funcao printf da LibC
    lea rdi, [fmt_out_num]              ; Formato de saida em RDI
    mov rsi, [res]                      ; Valor da variavel em RSI
    xor eax, eax                        ; Limpa EAX para printf
    call printf                         ; Exibe o valor numerico

    add rsp, 32                         ; Libera espaco da pilha
    pop rbp                             ; Restaura o ponteiro de base
    mov rax, 0                          ; Status de saída zero (Sucesso)
    ret                                 ; Retorna ao Sistema Operacional
