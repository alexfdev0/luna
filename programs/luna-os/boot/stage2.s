.bits 16
.fill 1044
.org 492

PARTITION_TABLE:
    .pad 20 

_stage2:
    push 0x0F0F
    call screen_draw

    push msg_header
    push 255
    push 0x0F
    call write

    mov r1, 0
    mov r2, next_vol_num
    strf r2, r1
    // Print volumes
    push msg_select_boot_vol
    push 255
    push 0x0f
    call write

    call list_volumes

    push msg_opts
    push 255
    push 0x0f
    call write
VOL_INP:
    // Tell user to select prompt
    int 6 

    mov e1, 0x0a
    cmp e2, e1, e12
    jnz e2, REBOOT // reboot if enter

    mov r1, 0x0a
    int 1

    mov e0, "0"
    sub e1, e12, e0 // OFFSET IN E1
    mov e0, 2
    mul e1, e1, e0

    mov e2, PARTITION_TABLE
    add e2, e2, e1

    lodf e2, e3

    jz e3, vol_error

    // Jump to OS
    push 0x0000
    call screen_draw

    jmp e3 

list_volumes:
    pop e11
    push e11

    mov r4, PARTITION_TABLE
    mov r6, 24
    mov r8, 20
    mov r9, 0

    mov e10, pc

    lodf r4, r5 
    jz r5, list_volumes_ret

    push r4
    push r6
    push e10

    sub r7, r5, r6

    mov r2, 255
    mov r3, 0x0f
    
    mov e0, next_vol_num
    mov e1, "0"
    lodf e0, e0
    add r1, e0, e1
    int 1

    mov r1, 0x20
    int 1

    mov r1, "-"
    int 1

    mov r1, 0x20
    int 1

    push r7
    push 255
    push 0x0f
    call write

    mov r1, 0x0a
    int 1

    inc e0
    mov e1, next_vol_num
    strf e1, e0

    pop e10
    pop r6
    pop r4

    inc r4
    inc r4
    inc r9
    inc r9 

    cmp r10, r9, r8
    jnz r10, list_volumes_ret
    
    jmp e10
list_volumes_ret:
    pop e11
    ret

vol_error:
    push msg_incorrect_vol
    push 255
    push 0x0F
    call write
    jmp VOL_INP

write:
    pop e11
    pop r3
    pop r2
    pop r4

    mov e10, pc
    nop

    lod r4, r1
    int 1
    inc r4

    jnz r1, e10

    ret

screen_draw:
    pop e11
    pop r2

    mov r1, 0x0
    mov r4, 0
    mov r3, 64000

    mov e10, pc
    int 3

    inc r4
    inc r4
    inc r1
    inc r1

    igt r5, r4, r3
    jz r5, e10

    mov r1, 0
    mov r2, 0
    int 0xc

    ret

REBOOT:
    int 0x10 
    int 0xf

next_vol_num:
    .pad 2

msg_select_boot_vol:
    .asciz "Select boot volume:\n"

msg_incorrect_vol:
    .asciz "Invalid volume\n"

msg_header:
    .asciz "Luna Boot Menu\n\n"

msg_opts:
    .asciz "\nEnter: Reboot"
