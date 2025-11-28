.bits 32

.global lufs_create_file
.global lufs_write_file
lufs_create_file:
    pop e11 
    pop r2 // FILE SIZE (bytes)
    pop r1 // NAME BUFFER

    mov r3, 0x41c
    lodf r3, r4 // Load the next file pointer
    mov e12, r4

    mov r5, "LFSF"
    strf r4, r5 // STORE FILE HEADER
    inc r4
    inc r4
    inc r4
    inc r4
    
    mov r6, 0
    mov r9, 16
    mov e10, pc
    nop
    
    lod r1, r7
    str r4, r7 // TRANSFER NAME TO FILE

    inc r6
    inc r1
    inc r4

    cmp r10, r6, r9
    jz r10, e10

    strf r4, r2 // Store file size buffer

    inc r4
    inc r4
    inc r4
    inc r4

    mov r6, 0
    mov r8, 0
    mov e10, pc
    nop

    str r4, r8

    inc r6
    igt r3, r6, r2
    jz r3, e10 // Allocate file full size

    mov r12, 512
    div r1, e12, r12 // Sector number in r12
    mov r2, 0
    int 0x0d

    ret

lufs_write_file:
    pop e11
    pop r2 // BUFFER
    pop r1 // NAME
    push e11
    
    mov r3, 0x41c
    lodf r3, r3 // Load the next file pointer
lufs_write_file_top:
    mov r5, "LFSF"
    lodf r3, r4
    cmp r6, r4, r5
    jnz r6, lufs_find_file
    inc r3
    inc r3
    inc r3
    inc r3
    jmp lufs_write_file_top
lufs_find_file:
    inc r3
    inc r3
    inc r3
    inc r3
    
    // Call strcmp

    push r1
    push r2 // Save registers
    push r3

    push r1
    push r3
    call strcmp

    pop r3
    pop r2
    pop r1 // Restore registers

    jnz e6, lufs_find_file_match
    jz e6, lufs_find_file_fail
lufs_find_file_match:
    mov r10, 16
    add r3, r3, r10 // SKIP OVER NAME + SIZE
    
    lodf r3, r12
    push r12
    inc r3
    inc r3
    inc r3
    inc r3 // SAVE FILE SIZE PTR

    mov e10, pc
    nop

    lod r2, r9
    str r3, r9

    inc r2
    inc r3

    jnz r9, e10

    mov r12, 512
    div r1, r3, r12 // Sector number in r12
    mov r2, 0
    int 0x0d

    inc r1
    int 0x0d

    dec r1
    dec r1
    int 0x0d
    
    mov e6, 1

    // SAVE NEXT FILE PTR
    mov r1, 0x41c
    lodf r1, r2
    pop r3
    add r2, r2, r3
    strf r1, r2
    mov r1, 1
    mov r2, 0
    int 0x0d

    pop e11
    ret
lufs_find_file_fail:
    inc r3
    inc r3
    inc r3
    inc r3
    jmp lufs_write_file_top
