.bits 16
.org 1536

.global IDT_SETUP
.global targeted_load
.global boot_load_all_sectors

// LUFS header
.pad 32

set 32
.bits 32
LOS_BASE:
mov sp, 0x6FFEFFFF

call IDT_SETUP
jmp boot

IDT_SETUP:
    pop e11

    // Set up IDT

    mov r1, 0x6FFF0025
    mov r2, 1
    str r1, r2

    mov r1, 0x6FFF0026
    mov r2, kernel_panic
    strf r1, r2

    mov r1, 0x6FFF0013
    mov r2, 1
    str r1, r2

    mov r1, 0x6FFF0014
    mov r2, syscall_handler
    strf r1, r2

    mov r1, 0x6FFF001A
    mov r2, key_click
    strf r1, r2

    mov r1, 0x6FFF0020
    mov r2, wait_for_key
    strf r1, r2

    mov r1, pit_nxt
    mov r2, 0x6FFF0008
    strf r2, r1

    ret

targeted_load:
    pop e11
    pop r9 // Sectors
    pop r1 // Address

    push r1

    int 0x10
    mov r2, r1

    pop r1

    mov r6, 0

    mov r4, 512
    div r1, r1, r4

    dec r1
    mov r3, r1
    int 0x0b

    inc r1

    mov e10, pc

    mov r3, r1
    int 0x0b
    inc r1
    inc r6

    igt r5, r6, r9
    jz r5, e10

    ret

boot_load_all_sectors:
    pop e11
    pop r6 // num sectors

    int 0x10
    mov r2, r1

    mov r1, LOS_BASE
    mov r4, 512
    div r1, r1, r4 // get base

    mov r7, 0

    mov e10, pc

    mov r3, r1
    int 0x0b
    inc r1
    inc r7

    igt r9, r7, r6
    jz r9, e10

    ret
    
