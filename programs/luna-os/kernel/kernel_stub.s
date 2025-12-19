.bits 16
.org 1536

// LUFS header
.pad 32

set 32
.bits 32
mov sp, 0x6FFEFFFF

// Set up IDT
mov r1, 0x6FFF0025
mov r2, 1
str r1, r2

mov r1, 0x6FFF0026
mov r2, kernel_panic
strf r1, r2

mov r1, 0x6FFF0061
mov r2, 1
str r1, r2

mov r1, 0x6FFF0062
mov r2, mouse_move
strf r1, r2

// Jump to C code
jmp _cstart
