.bits 16
.fill 492

jmp _start

_start:
    // Setup stack
    mov sp, 0xEFFF 

    // Load next sector
    int 0x10
    mov r2, r1
    mov r1, 1
    int 11

    int 0x10
    mov r2, r1
    mov r1, 2
    int 11

    // Load sectors
    push msg_loading
    push 255
    push 0
    call write

    mov r4, num_sectors
    lodf r4, r4
    push r4
    call load_sectors

    // Check partition table
    call check_vol
    jz e6, missing_error

    jmp 512 

load_sectors:
    pop e11
    pop e1

    int 0x10
    mov r2, r1

    mov r1, 3

    mov e10, pc
    nop

    int 11
    inc r1

    igt r3, r1, e1
    jz r3, e10

    ret

wait_key_and_reboot:
    push msg_press_any_key
    push 255
    push 0
    call write

    int 6
    jmp 0

missing_error:
    push msg_missing_os
    push 255
    push 0
    call write
    jmp wait_key_and_reboot

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

check_vol:
    pop e11
    push e11

    mov r5, 492
    mov r7, 0
    mov r8, 0
    mov r9, 20
    
    mov e10, pc
    nop

    lodf r5, r6

    inc r7
    inc r7
    inc r8
    inc r8

    cmp r10, r8, r9
    jnz r10, check_vol_ret

    jz r6, e10
    inc r7
    jmp e10    
check_vol_ret:
    mov e6, r7
    pop e11
    ret

num_sectors: 
    .word 0x02CC

msg_loading:
    .asciz "Loading...\n\n"

msg_press_any_key:
    .asciz "Press any key to reboot...\n"

msg_missing_os:
    .asciz "Missing operating system...\n"
