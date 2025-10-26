.bits 16
.fill 512

_start:
    // Setup stack
    mov sp, 0xffff 

    // Check memory
    int 10
    mov r4, num_sectors
    mov r5, 512
    lodf r4, r4
    mul r5, r4, r5
    ilt r5, r1, r5
    jnz r5, memory_error

    // Load sectors
    push r4
    call load_sectors

    // Verify OS presence
    mov r4, 512
    mov r5, 0xAA55
    lodf r4, r4
    cmp r6, r4, r5
    jz r6, os_error

    // Jump to OS
    jmp 514

load_sectors:
    pop e11
    pop e1

    int 0x10
    mov r2, r1

    mov r1, 1

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

memory_error:
    push msg_memory_error
    push 255
    push 0
    call write
    jmp wait_key_and_reboot

os_error:
    push msg_no_os
    push 255
    push 0
    call write
    jmp wait_key_and_reboot

write:
    pop e11
    pop r3
    pop r2
    pop r4
    mov r5, 0 

    mov e10, pc
    nop

    lod r4, r1
    int 1
    inc r4

    cmp r6, r1, r5
    jz r6, e10

    ret 

num_sectors: 
    .word 0x000a

msg_memory_error:
    .asciz "There is not enough memory to load the operating system.\n"

msg_no_os:
    .asciz "No operating system found.\n"

msg_press_any_key:
    .asciz "Press any key to reboot...\n"
