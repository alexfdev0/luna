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
.global send
.global NETBUF
.global serve
.global serve_write
.global serve_read
.global serve_await_connection
.global serve_connection_close
.global render
.global sleep
.global screen_fill
.global render_buf
.global save_graphics_buf
.global GBUF
.global mouse_move
.global key_click

readin:
    pop e11
    pop e1
    pop e9
    pop r4

    mov r1, 0x6FFF0019
    mov r2, 1
    // str r1, r2 // ENABLE KEYBOARD INTERRUPT

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

    int 6
    mov r1, e12

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
    // str r1, r2 // DISABLE KEYBOARD INTERRUPT

    str r4, r3
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
    mov r3, 0xFC // key
    mov r4, 0

    mov e10, pc
    nop

    lod r2, r1

    cmp r5, r4, r1
    jnz r5, _GLOBL_RET

    xor r1, r1, r3
    str r2, r1

    inc r2
    jmp e10

_GLOBL_RET:
    ret
    
getdrive:
    pop e11

    int 0x10
    mov e6, r1

    ret

setup_copy:
    pop e11
    pop r3

    mov r1, 0
    mov r2, 0

    mov e10, pc
    nop

    int 0x0d
    inc r1

    igt r4, r1, r3
    jz r4, e10

    ret

checkpass:
    pop e11

    mov r1, PASSBUF
    lod r1, r1
    
    mov e6, r1

    ret

save_buffer:
    pop e11
    pop r2
    pop r1

    mov r3, 512
    div r4, r1, r3

    mov r1, r4
   
    int 0x0d

    ret

send:
    pop e11
    pop e7 // TIMEOUT
    pop r4 // PORT
    pop r3 // IP 
    pop r9 // BUFFER

    mov r7, 0x00FF
    mov r8, 0x0100
    mov r11, 0

    // move TCP flag to NIC
    mov r1, 0
    mov r2, 0x7001A645
    str r2, r1

    // move IP addr to nic
    mov r2, 0x7001A646
    strf r2, r3

    // move port
    // r6 high
    // r5 low
    
    div r6, r4, r8
    and r5, r4, r7

    mov r2, 0x7001A64a
    str r2, r6

    mov r2, 0x7001A64b
    str r2, r5

    // Move timeout
    div e8, e7, r8 // upper
    and e9, e7, r7 // lower

    mov r2, 0x7001A64c
    str r2, e8

    mov r2, 0x7001A64d
    str r2, e9

    // Move message
    mov r10, 0x7001A64e

    mov e10, pc
    nop

    lod r9, r1
    str r10, r1
    
    cmp r12, r11, r1
    jnz r12, send_send

    inc r9
    inc r10
    jmp e10
send_send:
    // SEND FLAG
    mov r1, 0x01
    mov r2, 0x7001A644
    str r2, r1

    mov r1, e7
    int 2

    // Copy response back

    mov r4, 0x7001ae47 // Start addr
    mov r5, NETBUF // Netbuf ptr

    mov r6, 0x7001B64C // Stop addr

    mov r3, 0

    mov e10, pc
    nop

    lod r4, r1
    str r5, r1

    inc r4
    inc r5
    
    igt r7, r4, r6
    jz r7, e10

    str r5, r3
 
    ret

serve:
    pop e11
    pop r4 // PORT

    // Move port
    mov r5, 0x00FF
    mov r6, 0x0100
    div r8, r4, r6
    and r7, r4, r5

    mov r2, 0x7001A64a
    str r2, r8

    mov r2, 0x7001A64b
    str r2, r7

    // move TCP flag to NIC
    mov r1, 1
    mov r2, 0x7001A645
    str r2, r1

    // Listen
    mov r1, 0x01
    mov r2, 0x7001A644
    str r2, r1

    ret

serve_await_connection:
    pop e11
    pop r1 // TIMEOUT

    mov r2, 0x7001A648

    mov e10, pc
    nop 

    lod r2, r4
    jnz r4, serve_await_connection_accept
     
    int 2
    jmp e10
serve_await_connection_accept: 
    mov r1, 0x01
    mov r2, 0x7001A646
    str r2, r1
    ret

serve_write:
    pop e11
    pop r9 
    
    // Copy messaege
    mov r10, 0x7001A64e
    mov r11, 0

    mov e10, pc
    nop

    lod r9, r1
    str r10, r1
    
    cmp r12, r11, r1
    jnz r12, serve_write_ready

    inc r9
    inc r10
    jmp e10
serve_write_ready:
    // Tell NIC we are ready to write
    mov r1, 0x02
    mov r2, 0x7001a649
    str r2, r1

    mov r1, 0x01
    mov r2, 0x7001a646
    str r2, r1

    // Polling
    mov r1, 500
    mov r2, 0x7001A647
    mov e10, pc
    nop
    lod r2, r3
    int 2
    jz r3, e10

    ret

serve_read:
    mov r4, 0x7001ae47 // Start addr
    mov r5, NETBUF // Netbuf ptr

    mov r6, 0x7001B64C // Stop addr

    mov r1, 500
    int 2
    
    // COMMAND NIC TO READ 
    mov r1, 0x01
    mov r2, 0x7001a649
    str r2, r1

    mov r1, 0x01
    mov r2, 0x7001a646
    str r2, r1

    // Polling
    mov r1, 500
    mov r2, 0x7001A647
    mov e10, pc
    nop
    lod r2, r3
    int 2
    jz r3, e10 

    mov e10, pc
    nop

    lod r4, r1
    str r5, r1

    inc r4
    inc r5
    
    igt r7, r4, r6
    jz r7, e10
 
    ret

serve_connection_close:
    pop e11
    
    // Tell NIC to close connection
    mov r1, 0x00
    mov r2, 0x7001a649
    str r2, r1

    mov r1, 0x01
    mov r2, 0x7001a646
    str r2, r1

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

sleep:
    pop e11
    pop r1

    int 2

    ret

screen_fill:
    pop e11
    pop r2

    mov r1, 0x0
    mov r4, 0
    mov r3, 64000

    mov e10, pc
    int 3

    inc r4
    inc r4
    inc r4
    inc r4
    inc r1
    inc r1
    inc r1
    inc r1

    igt r5, r4, r3
    jz r5, e10

    mov r1, 0
    mov r2, 0
    int 0xc

    ret

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

    mov r1, GBUF
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

mouse_move:
    pusha

    call save_graphics_buf
   
    mov r1, 0x7000FA0A
    lodf r1, r2 // X

    mov r1, 0x7000FA0E
    lodf r1, r3 // Y

    mov r4, 320

    mul r5, r3, r4
    add r5, r5, r2 // Turn from 320x200 space to 64KB space

    mov r6, GBUF
    add r6, r6, r5 // Get absolute addr

    mov r7, 0xA0
    str r6, r7

    push GBUF
    call render_buf 

    popa
    jmp irv

key_click:
    mov r2, 0x7000FA12
    lod r2, r1
    jmp e7
    
    
TEMPBUF:
    .pad 256

PROMPTBUF:
    .asciz "> "
    .pad 5

PASSBUF:
    .pad 32

NETBUF:
    .pad 2048

GBUF:
    .pad 64000
