.bits 16
.org 1536

.global IDT_SETUP

// LUFS header
.pad 32

set 32
.bits 32
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
