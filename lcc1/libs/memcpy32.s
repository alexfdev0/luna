.bits 32
.global _builtin_lcc_memcpy32

_builtin_lcc_memcpy32:
    pop e11
    pop r5 // Number of bytes to copy
    pop r6 // Source address
    pop r7 // Destination address

    mov r4, 0
    
    mov e10, pc

    lod r6, r8
    str r7, r8

    inc r6
    inc r7
    inc r4

    cmp r9, r5, r4
    jz r9, e10

    ret

