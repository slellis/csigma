section .data                           ; Area de dados (WORKING-STORAGE)
    fmt_in db ' %ld', 0                 ; Formato para entrada numerica
    fmt_out_num db '%ld', 10, 0         ; Formato para saida numerica
    C                    dq 0           ; Variavel C
    RESULTADO            dq 0           ; Variavel RESULTADO
    A                    dq 0           ; Variavel A
    B                    dq 0           ; Variavel B
    msg_0                db 'VALOR A:', 10, 0               ; Constante de texto
    msg_1                db 'VALOR B:', 10, 0               ; Constante de texto
    msg_2                db 'VALOR C:', 10, 0               ; Constante de texto
    msg_3                db 'RESULTADO FINAL:', 10, 0               ; Constante de texto

section .text                           ; Area de codigo (PROCEDURE DIVISION)
extern printf, scanf
global main

main:
    push rbp                            ; Prologo
    mov rbp, rsp
    sub rsp, 32                         ; Alinhamento de pilha

    lea rdi, [msg_0]                    ; Endereco da string para printf
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [A]                        ; Endereco de A em RSI
    xor eax, eax
    call scanf
    lea rdi, [msg_1]                    ; Endereco da string para printf
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [B]                        ; Endereco de B em RSI
    xor eax, eax
    call scanf
    lea rdi, [msg_2]                    ; Endereco da string para printf
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [C]                        ; Endereco de C em RSI
    xor eax, eax
    call scanf
    mov rax, [A]                        ; Carrega variavel A
    add rax, [B]                        ; Soma
    imul rax, 2                         ; Multiplica
    mov rbx, [C]                        ; Move divisor para RBX
    cqo                                 ; Estende sinal p/ RDX
    idiv rbx                            ; Divide RAX por RBX
    mov [RESULTADO], rax                ; Salva em RESULTADO
    lea rdi, [msg_3]                    ; Endereco da string para printf
    xor eax, eax
    call printf
    lea rdi, [fmt_out_num]
    mov rsi, [RESULTADO]                ; Valor de RESULTADO em RSI
    xor eax, eax
    call printf

    add rsp, 32                         ; Epilogo
    pop rbp
    mov rax, 0
    ret
