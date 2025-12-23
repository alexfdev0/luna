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

mov r1, 0x6FFF0025
mov r2, 1
str r1, r2

mov r1, 0x6FFF0026
mov r2, kernel_panic
strf r1, r2

mov r1, 0x6FFF0067
mov r2, 1
// str r1, r2

mov r1, 0x6FFF0068
mov r2, mouse_move
strf r1, r2

mov r1, 0x6FFF001A
mov r2, key_click
strf r1, r2

mov r1, 0x6FFF0020
mov r2, wait_for_key
strf r1, r2

// Set up PIT
mov r1, 0x7000FA13
mov r2, 0x1234E8
// strf r1, r2 // set PIT interval to 1 second

// Jump to C code
jmp _cstart
