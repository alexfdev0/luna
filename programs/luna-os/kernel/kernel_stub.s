.bits 16
.org 1024

// LUFS header
.pad 32

set 32
.bits 32
mov sp, 0x6FFFFFFF
jmp _cstart
