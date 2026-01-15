section .data
    fmt_in db '%ld', 0                  ; Scanf
    fmt_out_num db '%ld', 10, 0         ; Printf
    msg_0                db 'Calculadora', 10, 0
    msg_1                db 'Conta mista: (a + b) * 2 / C', 10, 0
    a                    dq 0           ; Variavel a
    b                    dq 0           ; Variavel b
    c                    dq 0           ; Variavel c
    res                  dq 0           ; Variavel res
    msg_2                db 'Valor de a:', 10, 0
    msg_3                db 'Valor de b:', 10, 0
    msg_4                db 'Valor de c:', 10, 0
    msg_5                db 'Resultado final:', 10, 0

section .text
extern printf, scanf
global main

main:
    push rbp
    mov rbp, rsp
    sub rsp, 32

    lea rdi, [msg_0]
    xor eax, eax
    call printf
    lea rdi, [msg_1]
    xor eax, eax
    call printf
    lea rdi, [msg_2]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [a]
    xor eax, eax
    call scanf
    lea rdi, [msg_3]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [b]
    xor eax, eax
    call scanf
    lea rdi, [msg_4]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [c]
    xor eax, eax
    call scanf

    ; --- Calculo de res ---
    mov rax, [a]
    add rax, [b]                      
    imul rax, 2                       
    mov rbx, [c]                      
    xor rdx, rdx
    idiv rbx
    mov [res], rax
    lea rdi, [msg_5]
    xor eax, eax
    call printf
    lea rdi, [fmt_out_num]
    mov rsi, [res]
    xor eax, eax
    call printf

    add rsp, 32
    pop rbp
    mov rax, 0
    ret
