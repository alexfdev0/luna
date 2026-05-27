.bits 16
.fill 492

jmp _start

_start:
    // Setup stack
    mov sp, 0xEFFF

    // Check battery
    hlt
    hlt // allow time for the battery controller to initialize
    call bat_check

    // Check partition table
    call check_vol

    // Print loading message
    push msg_loading
    push 255
    call write

    // Load next sectors
    int 0x10
    mov r2, r1
    mov r1, 1
    mov r3, r1
    int 11

    int 0x10
    mov r2, r1
    mov r1, 2
    mov r3, r1
    int 11

    jmp 512

wait_key_and_reboot:
    push msg_press_any_key
    push 255
    call write

    call wait_key
    int 0x10
    int 0xf

missing_error:
    push msg_missing_os
    push 255
    call write
    jmp wait_key_and_reboot

write:
    pop e11
    pop r2
    pop r4

    xor r3, r3, r3

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

bat_check:
    mov r1, 0xFC3A
    lod r1, r2
    mov r3, 3
    igt r4, r2, r3

    jnz r4, bc_ret

    push msg_battery_dead
    push 0xA0
    call write
    jmp pc
bc_ret:
    pop e11
    ret

msg_loading:
    .asciz "Loading\n"

msg_press_any_key:
    .asciz "Press any key to reboot\n"

msg_missing_os:
    .asciz "Missing operating system\n"

msg_battery_dead:
    .asciz "Charge battery"
