.bits 32

.global readin
.global strcmp
.global puts32
.global writeout
.global xor_cycle
.global getdrive
.global PASSBUF
.global TEMPBUF
.global PROMPTBUF
.global setup_copy
.global checkpass
.global save_buffer
.global render
.global screen_fill
.global render_buf
.global save_graphics_buf
.global mouse_move
.global key_click
.global wait_for_key
.global modulo
.global itoa
.global malloc
.global free
.global ASLR_generate_address
.global lexec_core
.global sectorize
.global syscall_handler
.global sleep
.global pit_nxt
.global renderbuf_loc
.global sleep_loc

readin:
    pop e11
    pop e1
    pop e9
    pop r4
    push e11 

    mov r12, r4

    jnz e9, readin_rdy
    mov r5, 0
    lod r4, r1
    cmp r6, r1, r5
    jnz r6, readin_rdy

    mov e8, e11

    push r4
    push 255
    push 0
    call puts32
    mov r4, e6

    mov e11, e8
readin_rdy:
    mov r2, 255
    mov r3, 0
    mov r5, 0x0a
    mov r7, 0xc3

    mov e10, pc
    nop
    
    push r1
    push r2
    mov r1, 0x6FFF0019
    mov r2, 1
    str r1, r2 // ENABLE KEYBOARD INTERRUPT
    pop r2
    pop r1 
readin_rd:
    // AWAIT KEYBOARD INT
    mov e7, readin_ai
    hlt
    jmp readin_rd
readin_ai:
    push r1
    push r2
    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT
    pop r2
    pop r1

    cmp r8, r1, r7
    jnz r8, readin_bksp
    jnz e1, readin_bt
    int 1
    jmp readin_nxt
readin_bt:
    push e8
    push e7
    mov e8, 0x0a
    cmp e7, e8, r1
    jnz e7, readin_bt_nl
    push r1
    mov r1, "*"
    int 1
    pop r1
    jmp readin_bt_done
readin_bt_nl:
    mov r1, 0x0a
    int 1
readin_bt_done:
    pop e7
    pop e8
readin_nxt:
    cmp r6, r1, r5
    jnz r6, readin_done

    str r4, r1
    inc r4
    jmp e10
readin_bksp:
    cmp e7, r4, r12
    jnz e7, e10

    mov r1, 0
    mov r10, r2
    mov r11, r3 

    dec r4
    str, r4, r1

    int 0x0e
    dec r1
    int 0x0c

    mov r1, 65
    mov r2, 0
    mov r3, 0
    int 1

    int 0x0e
    dec r1
    int 0x0c

    mov r2, r10
    mov r3, r11
    jmp e10
readin_done:
    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT

    str r4, r3
    pop e11
    ret

strcmp:
    pop e11
    pop r2
    pop r1

    mov r6, 0 

    mov e10, pc
    nop 

    lod r1, r3
    lod r2, r4

    inc r1
    inc r2

    cmp r5, r3, r4
    jz r5, strcmp_false

    cmp r5, r3, r6
    jz r5, e10

    mov e6, 1
    ret
strcmp_false:
    mov e6, 0
    ret

strlen:
    pop e11
    pop r1
    mov r2, 0
    mov e6, 0

    mov e10, pc
    nop

    lod r1, r3
    inc r1
    inc e6

    cmp r5, r3, r2
    jz r5, e10

    ret

writeout:
    pop e11
    pop r2
    pop r1
    
    int 0xd

    ret

puts32:
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

    cmp r6, r5, r1
    jz r6, e10

    dec r4
    mov e6, r4
    ret

xor_cycle:
    pop e11
    pop r2

    ret // see if this fixes it

    mov r3, PASS_KEY_ENCRYPT
    lod r3, r3
    jz r3, xor_cycle_gen
xor_cycle_rdy:
    mov r4, 0

    mov e10, pc
    nop

    lod r2, r1

    cmp r5, r4, r1
    jnz r5, xor_ret

    xor r1, r1, r3
    str r2, r1

    inc r2
    jmp e10
xor_cycle_gen:
    mov r11, PASS_KEY_ENCRYPT
    mov r12, 0x90000000
    lod r12, r3
    str r11, r3

    push e11
    push r1
    push r3
    push r4

    push PASS_KEY_ENCRYPT
    call save_buffer

    pop r4
    pop r3
    pop r1
    pop e11

    jmp xor_cycle_rdy
xor_ret:
    ret
    
getdrive:
    pop e11

    int 0x10
    mov e6, r1

    ret

setup_copy:
    pop e11
    pop e7 // drive number (r2 = e7)
    pop e8 // max sectors (r5 = e8)
    push e11

    // For text 

    // Copying
    
    mov e9, 0 // (r1 = e9)
    mov e10, pc
  
    mov r1, e9
    mov r2, e7
    mov r3, r1
    int 0x0d
    inc e9
    
    push e10
    push e7
    push e8
    push e9

    push e9
    push 1
    push 0x20202020
    call itoa
    
    push 0x20202020
    push 0xFC
    push 0
    call puts32 // Print first number

    push setup_copy_sectors
    push 0xFC
    push 0
    call puts32 // "sectors" 

    push setup_copy_copied
    push 0xFC
    push 0
    call puts32 // "copied"

    int 0x0e
    mov r1, 0
    int 0x0c

    pop e9
    pop e8
    pop e7
    pop e10

    igt r4, e9, e8
    jz r4, e10    

    pop e11
    ret
setup_copy_sectors:
    .asciz " blocks"
setup_copy_of:
    .asciz " of "
setup_copy_copied:
    .asciz " copied"

checkpass:
    pop e11

    mov r1, PASSBUF
    lod r1, r1
    
    mov e6, r1

    ret

save_buffer:
    pop e11 
    pop r1

    mov r5, 512
    div r4, r1, r5
    mov r1, r4

    mov r2, 0 // Drive 0
    mov r3, r1
    int 0x0d

    ret

render:
    pop e11
    pop r4 // Buffer

    mov r1, 1

    mov e10, pc
    nop

    lod r4, r2
    mov r3, r2
    int 1
    inc r4

    jnz r2, e10
    ret

screen_fill:
    pop e11
    pop r2

    mov r1, 0x70000000 // Current addr
    mov r4, 0 // total
    mov r3, 64000 // pixels constant

    mov e10, pc

    strf r1, r2

    inc r1
    inc r1
    inc r1
    inc r1
    inc r4
    inc r4
    inc r4
    inc r4

    igt r5, r4, r3
    jz r5, e10

    mov r1, 0
    mov r2, 0
    int 0xc
    
    ret 

renderbuf_loc:
render_buf:
    pop e11
    pop r1 // BUFFER

    mov r2, 0x70000000
    mov r3, 0
    mov r4, 64000

    mov e10, pc

    lodf r1, r5
    strf r2, r5 // we'll use fast here :3

    inc r1
    inc r1
    inc r1
    inc r1
    inc r3
    inc r3
    inc r3
    inc r3
    inc r2
    inc r2
    inc r2
    inc r2

    cmp r6, r3, r4
    jz r6, e10

    ret

save_graphics_buf:
    pop e11

    mov r1, 0x30303030
    mov r2, 0x70000000
    mov r3, 0
    mov r4, 64000

    mov e10, pc

    
    lodf r2, r5
    strf r1, r5

    inc r1
    inc r1
    inc r1
    inc r1
    inc r3
    inc r3
    inc r3
    inc r3
    inc r2
    inc r2
    inc r2
    inc r2

    cmp r6, r3, r4
    jz r6, e10

    ret

wait_for_key:
    pop e11
    
    mov e7, wfk_ret

    mov r1, 0x6FFF0019
    mov r2, 1
    str r1, r2 // ENABLE KEYBOARD INTERRUPT
wfk_rd:
    hlt
    jmp wfk_rd
wfk_ret:
    mov r1, 0x6FFF0019
    mov r2, 0
    str r1, r2 // DISABLE KEYBOARD INTERRUPT
    ret

key_click:
    push r2

    mov r2, 0x7000FA12
    lod r2, r1

    pop r2

    jmp e7

ASLR_generate_address:
    pop e11
    push e11

    mov r1, ASLR_addr 
    
    mov r2, 0x60
    str r1, r2 // 1
    inc r1

    mov r3, 0x90909090
    lod r3, r4

    str r1, r4 // 2
    inc r1

    lod r3, r4

    str r1, r4 // 3
    inc r1

    lod r3, r4

    str r1, r4 // 4
    inc r1

    mov r3, 0xFFFFFE00
    mov r1, ASLR_addr
    lodf r1, r2
    and r2, r2, r3
    strf r1, r2 // align to 1 sector

    mov e6, r2

    pop e11
    ret
ASLR_addr:
    .dword 57800

lexec_core:
    pop e11
    pop r1

    mov r1, ASLR_addr
    lodf r1, r1

            // at 7f
    inc r1 // l
    inc r1 // 2
    inc r1 // p
    inc r1 // i
    inc r1 // e
    inc r1 // Bitness indicator (ignore for now)
    inc r1 // 0
    inc r1 // 0
    inc r1 // 0
    inc r1 // 0
    inc r1 // first byte of program 
    
    mov r6, lexec_raddr
    strf r6, e11

    mov r3, ASLR_addr
    lodf r3, r3

    pusha
    mov e14, r3

    jmp r1
lexec_raddr:
    .dword 0x00000000
lexec_done:
    popa

    call IDT_SETUP

    jmp shell
    

sectorize:
    pop e11
    pop r1
    
    mov r2, 512
    div e6, r1, r2

    ret

syscall_handler:
    pusha
    
    mov r2, 1
    cmp r3, r1, r2
    jnz r3, syscall_proc_exit
syscall_proc_exit:
    popa
    jmp lexec_done
syscall_ret:
    popa
    jmp irv


sleep_loc:
sleep:
    pop e11
    pop r4 // Seconds
    push e11

    mov r5, 0 

    mov e10, pc

    call pit_handler
    inc r5

    cmp r6, r4, r5
    jz r6, e10

    pop e11
    ret


pit_handler:
    pop e11

    mov r1, 0x6FFF0007
    mov r2, 1
    str r1, r2 // ENABLE PIT
pit_wait:
    hlt
    jmp pit_wait
pit_nxt:
    mov r1, 0x6FFF0007
    mov r2, 0
    str r1, r2 // DISABLE PIT

    ret

// 1 - 6 destroyed
// 

itoa:
    pop e11
    pop r12 // Output buffer
    pop r2 // Capitalized?
    pop r1 // Number

    mov e12, r12

    jnz r2, itoa_c
    jmp itoa_l
itoa_c:
    mov r3, ITOA_LETTERS_U
    jmp itoa_after
itoa_l:
    mov r3, ITOA_LETTERS_L
itoa_after:
    mov r4, 28 // bits to shift by
    mov r5, 0xF
    mov r6, 0 // current number of digits done
    mov r8, 8 // max digits

    mov e10, pc

    shr r7, r1, r4
    and r7, r7, r5
    add e7, r3, r7

    lod e7, e8
    str r12, e8

    inc r6
    inc r12
    dec r4
    dec r4
    dec r4
    dec r4

    cmp r9, r6, r8
    jz r9, e10

    mov e6, e12

    ret
ITOA_LETTERS_U:
    .asciz "0123456789ABCDEF"
ITOA_LETTERS_L:
    .asciz "0123456789abcdef"

malloc:
    pop e11
    pop r1 // Size

    mov r2, MEM_PTR
    lodf r2, r3
    mov e6, r3
    sub r3, r3, r1
    strf r2, r3

    ret

free:
    pop e11
    pop r1 // Size

    mov r2, MEM_PTR
    lodf r2, r3
    add r3, r3, r1
    strf r2, r3

    ret

PASS_KEY_ENCRYPT:
    .pad 1
    
TEMPBUF:
    .pad 256

PROMPTBUF:
    .asciz "> "
    .pad 5

PASSBUF:
    .pad 32

MEM_PTR:
    .dword 0x50505050
