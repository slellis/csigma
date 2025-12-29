section .data
    fmt_in db ' %ld', 0
    fmt_out_num db '%ld', 10, 0
    B                    dq 0
    C                    dq 0
    D                    dq 0
    RESULTADO            dq 0
    A                    dq 0
    msg_0                db '--- INICIO DO CALCULO ---', 10, 0
    msg_1                db 'VALOR A:', 10, 0
    msg_2                db 'VALOR B:', 10, 0
    msg_3                db 'VALOR C:', 10, 0
    msg_4                db 'VALOR D:', 10, 0

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
    lea rdi, [fmt_in]
    lea rsi, [A]
    xor eax, eax
    call scanf
    lea rdi, [msg_2]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [B]
    xor eax, eax
    call scanf
    lea rdi, [msg_3]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [C]
    xor eax, eax
    call scanf
    lea rdi, [msg_4]
    xor eax, eax
    call printf
    lea rdi, [fmt_in]
    lea rsi, [D]
    xor eax, eax
    call scanf
    mov rax, [A]
    add rax, [B]
    sub rax, [C]
    sub rax, [D]
    mov [RESULTADO], rax

    add rsp, 32
    pop rbp
    mov rax, 0
    ret
