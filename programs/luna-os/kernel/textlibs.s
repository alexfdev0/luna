.bits 32
.global atoi_int

atoi_int:
    pop e11
    pop r1

    mov r2, 48
    mov r4, 0 // Result
    mov r5, 10


    mov e10, pc
    nop

    lod r1, r3

    mul r4, r4, r5
    sub r3, r3, r2
    add r4, r4, r3

    inc r1
    lod r1, r3

    jnz r3, e10
    mov e6, r4 // E6 is the return register

    ret
