.bits 16
.global pit_nxt
.global sleep
.global putc

sleep:
    pop e11
    pop r4 // Seconds
    push e11

    mov r5, 0 

    mov e10, pc

    call pit_handler
    inc r5

    cmp r6, r4, r5
    jz r6, e10

    pop e11
    ret

pit_handler:
    pop e11

    // mov r1, 0x6FFF0007 // 0xFA3E for 16 bit
    mov r1, 0xFA3E
    mov r2, 1
    str r1, r2
pit_wait:
    hlt
    jmp pit_wait
pit_nxt:
    // mov r1, 0x6FFF0007 // 0xFA3E for 16 bit
    mov r1, 0xFA3E
    mov r2, 0
    str r1, r2

    ret

putc:
    pop e11
    pop r2
    pop r1

    mov r3, r2

    int 1

    ret
