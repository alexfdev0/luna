.bits 32
.global readin
.global strcmp
.global puts32
.global writeout
.global xor_cycle
.global getdrive
.global FILE
.global PASSBUF
.global TEMPBUF
.global PROMPTBUF
.global setup_copy
.global checkpass

readin:
    pop e11
    pop e9
    pop r4
    mov r12, r4

    jnz e9, readin_rdy
    mov r5, 0
    lod r4, r1
    cmp r6, r1, r5
    jnz r6, readin_rdy

    mov e8, e11

    push r4
    push 255
    push 0
    call puts32
    mov r4, e6

    mov e11, e8
readin_rdy:
    mov r2, 255
    mov r3, 0
    mov r5, 0x0a
    mov r7, 0xc3

    mov e10, pc
    nop

    int 6
    mov r1, e12

    cmp r8, r1, r7
    jnz r8, readin_bksp

    int 1

    cmp r6, r1, r5
    jnz r6, readin_done

    str r4, r1
    inc r4
    jmp e10
readin_bksp:
    cmp e7, r4, r12
    jnz e7, e10

    mov r1, 0
    mov r10, r2
    mov r11, r3 

    dec r4
    str, r4, r1

    int 0x0e
    dec r1
    int 0x0c

    mov r1, 65
    mov r2, 0
    mov r3, 0
    int 1

    int 0x0e
    dec r1
    int 0x0c

    mov r2, r10
    mov r3, r11
    jmp e10
readin_done:
    str r4, r3
    ret

strcmp:
    pop e11
    pop r2
    pop r1

    mov r6, 0 

    mov e10, pc
    nop 

    lod r1, r3
    lod r2, r4

    inc r1
    inc r2

    cmp r5, r3, r4
    jz r5, strcmp_false

    cmp r5, r3, r6
    jz r5, e10

    mov e6, 1
    ret
strcmp_false:
    mov e6, 0
    ret

strlen:
    pop e11
    pop r1
    mov r2, 0
    mov e6, 0

    mov e10, pc
    nop

    lod r1, r3
    inc r1
    inc e6

    cmp r5, r3, r2
    jz r5, e10

    ret

writeout:
    pop e11
    pop r2
    pop r1
    
    int 0xd

    ret

puts32:
    pop e11
    pop r3
    pop r2
    pop r4

    mov r5, 0

    mov e10, pc
    nop

    lod r4, r1
    int 1
    inc r4

    cmp r6, r5, r1
    jz r6, e10

    dec r4
    mov e6, r4
    ret

xor_cycle:
    pop e11
    pop r2
    mov r3, 0xFC // key
    mov r4, 0

    mov e10, pc
    nop

    lod r2, r1

    cmp r5, r4, r1
    jnz r5, _GLOBL_RET

    xor r1, r1, r3
    str r2, r1

    inc r2
    jmp e10

_GLOBL_RET:
    ret
    
getdrive:
    pop e11

    int 0x10
    mov e6, r1

    ret

setup_copy:
    pop e11
    pop r3

    mov r1, 0
    mov r2, 0

    mov e10, pc
    nop

    int 0x0d
    inc r1

    igt r4, r1, r3
    jz r4, e10

    ret

checkpass:
    pop e11

    mov r1, PASSBUF
    lod r1, r1
    
    mov e6, r1

    ret

FILE:
    .pad 512

TEMPBUF:
    .pad 512

PROMPTBUF:
    .asciz "> "
    .pad 13

PASSBUF:
    .pad 64
