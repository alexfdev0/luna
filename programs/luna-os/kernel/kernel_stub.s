.bits 16
.org 512
.noentry

.word 0xAA55

set 32
.bits 32
mov sp, 0x6FFFFFFF
jmp _cstart
