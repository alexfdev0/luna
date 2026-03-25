.bits 32
.global play_sound
.global play_sound_loc
play_sound_loc:
var_2:
    .ptr 0x00
    .ptrlabel var_2
var_3:
    .dword 0x00000000
var_4:
    .byte 0x00
var_5:
    .ptrlabel var_5
    .ptr 0x00
var_6:
    .ptrlabel var_6
    .ptr 0x00
var_7:
    .ptrlabel var_7
    .ptr 0x00
var_8:
    .ptrlabel var_8
    .ptr 0x00

play_sound:
    pop e11
    pop e2
    pop e1
    pop e0
    push e11
    mov r1, var_2
    strf r1, e0
    mov r1, var_3
    strf r1, e1
    mov r1, var_4
    str r1, e2
    mov r1, 1879112201
    mov r2, r1
    mov r4, r2
    mov r7, var_5
    strf r7, r4
    mov r1, var_5
    mov r2, r1
    lodf r2, r2
    mov r4, r2
    mov r1, 0
    mov r2, r1
    mov r5, r2
    str r4, r5
    mov r1, 1879112193
    mov r2, r1
    mov r4, r2
    mov r7, var_6
    strf r7, r4
    mov r1, var_6
    mov r2, r1
    lodf r2, r2
    mov r4, r2
    mov r1, var_2
    lodf r1, r2
    mov r5, r2
    strf r4, r5
    mov r1, 1879112197
    mov r2, r1
    mov r4, r2
    mov r7, var_7
    strf r7, r4
    mov r1, var_7
    mov r2, r1
    lodf r2, r2
    mov r4, r2
    mov r1, var_3
    lodf r1, r2
    mov r5, r2
    strf r4, r5
    mov r1, 1879112192
    mov r2, r1
    mov r4, r2
    mov r7, var_8
    strf r7, r4
    mov r1, var_8
    mov r2, r1
    lodf r2, r2
    mov r4, r2
    mov r1, 1
    mov r2, r1
    mov r5, r2
    str r4, r5
    mov r1, var_4
    lod r1, r2
    mov r11, r2
    jnz r11, if_stmt_11
    jz r11, after_stmt_13
if_stmt_11:
while_stmt_14_check:
    mov r1, var_5
    lodf r1, r2
    lod r2, r2
    mov r11, r2
    mov r1, 0
    mov r2, r1
    mov r5, r2
    cmp r11, r11, r5
    jnz r11, while_stmt_14_body
    jmp while_stmt_14_after
while_stmt_14_body:
    jmp while_stmt_14_check
while_stmt_14_after:
    jmp after_stmt_13
after_stmt_13:
    pop e11
    ret
    pop e11
    ret
