.bits 16
.fill 492

jmp _start

_start:
    // Setup stack
    mov sp, 0xEFFF

    // Check partition table
    call check_vol

    // Print loading message
    push msg_loading
    call write

    // Load next sectors
    int 0x10
    mov r2, r1
    mov r1, 1
    mov r3, r1
    int 11
    jz r0, read_error

    int 0x10
    mov r2, r1
    mov r1, 2
    mov r3, r1
    int 11
    jnz r0, read_error

    jmp 512

wait_key_and_reboot:
    push msg_press_any_key
    call write

    call wait_key
    int 0x10
    int 0xf

missing_error:
    push msg_missing_os
    call write
    jmp wait_key_and_reboot

read_error:
    push msg_read_error
    call write
    jmp wait_key_and_reboot

write:
    pop e11
    pop r4

    mov r2, 255
    mov r3, 0

    mov e10, pc

    lod r4, r1
    int 1
    inc r4

    jnz r1, e10

    ret

check_vol:
    pop e11

    mov r1, 492 // addr of partition table
    mov r3, 512

    mov e10, pc

    lod r1, r2
    jnz r2, check_vol_ret

    inc r1
    inc r1

    igt r4, r1, r3
    jz r4, e10
    jmp missing_error
check_vol_ret:
    ret

wait_key:
    pop e11

    mov e7, wk_after

    mov r1, 0xFA53
    mov r2, key_inp
    strf r1, r2  // SET KEY CLICK ADDR

    mov r1, 0xFA50
    mov r2, 1
    str r1, r2 // ENABLE KEY INP
wk_wait:
    hlt
    jmp wk_wait
wk_after:
    mov r1, 0xFA50
    mov r2, 0
    str r1, r2 // DISABLE KEY INP
    ret

key_inp:
    mov r1, 0xFA12
    lod r1, e12
    jmp e7

msg_loading:
    .asciz "Loading\n"

msg_press_any_key:
    .asciz "Press any key to reboot\n"

msg_missing_os:
    .asciz "Missing operating system\n"

msg_read_error:
    .asciz "Read from disk failed\n"
