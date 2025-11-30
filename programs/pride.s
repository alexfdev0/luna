.bits 16

jmp _start

_start:
    // Set up low-level stuff

    mov sp, 0xffff
    mov r1, 1
    int 11
    mov r1, 2
    int 11
    mov r1, 3
    int 11
    mov r1, 4
    int 11
    mov r1, 5
    int 11
    mov r1, 6
    int 11

    // Set up loop

    mov e9, pc
    nop

    // Pride flag

    push 0xe0
    call draw
    
    push 0xf0
    call draw

    push 0xfc
    call draw

    push 0x10
    call draw

    push 0x0b
    call draw

    push 0x62
    call draw

    // Space

    push 0x00
    call draw
    
    push 0x00
    call draw

    // Trans flag

    push 0x5b
    call draw
    
    push 0xf6
    call draw

    push 0xff
    call draw

    push 0xf6
    call draw

    push 0x5b
    call draw

    push 0x00
    call draw

    // Asexual flag

    push 0x24
    call draw

    push 0xb6
    call draw

    push 0xff
    call draw

    push 0x81
    call draw

    // Clear screen

    call wait

    push 0x00
    call draw

    push 0x00
    call draw 

    // Gay flag

    push 0x11
    call draw

    push 0x56
    call draw
    
    push 0xba
    call draw

    push 0xff
    call draw

    push 0x77
    call draw

    push 0x4a
    call draw
    
    push 0x45
    call draw

    // Clear screen

    push 0x00
    call draw

    push 0x00
    call draw

    // Lesbian flag
    
    push 0xc4
    call draw

    push 0xf1
    call draw

    push 0xff
    call draw

    push 0xce
    call draw
    
    push 0xa1
    call draw

    // Clear screen

    push 0x00
    call draw

    push 0x00
    call draw

    // Bisexual flag

    push 0xc1
    call draw

    push 0xc1
    call draw

    push 0x6a
    call draw

    push 0x0a
    call draw

    push 0x0a
    call draw

    // Wait
    call wait

    // Clear screen
    push 0x00
    call draw

    push 0x00
    call draw
    
    push 0x00
    call draw

    push 0x00
    call draw

    // Pansexual flag

    push 0xe6
    call draw

    push 0xe6
    call draw

    push 0xf8
    call draw

    push 0xf8
    call draw

    push 0x57
    call draw

    push 0x57
    call draw

    // Clear screen
    push 0x00
    call draw

    push 0x00
    call draw

    // Genderfluid flag
    push 0xee
    call draw

    push 0xff
    call draw
    
    push 0xc3
    call draw

    push 0x24
    call draw

    push 0x27
    call draw

    // Clear screen
    push 0x00
    call draw

    push 0x00
    call draw

    // Aromantic flag

    push 0x14
    call draw

    push 0x99
    call draw

    push 0xff
    call draw

    push 0xb6
    call draw

    push 0x24
    call draw

    // Wait
    call wait

    // Clear screen
    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    // Agender flag

    push 0x24
    call draw

    push 0xbb
    call draw

    push 0xff
    call draw
   
    push 0xbd
    call draw

    push 0xff
    call draw

    push 0xbb
    call draw

    push 0x24
    call draw

    // Clear screen
    push 0x00
    call draw 

    // Non-binary flag

    push 0xfc
    call draw

    push 0xff
    call draw

    push 0xab
    call draw

    push 0x24
    call draw

    // Clear screen
    push 0x00
    call draw
 

    // Polysexual flag

    push 0xe2
    call draw

    push 0x19
    call draw

    push 0x13
    call draw

    // Clear screen
    push 0x00
    call draw 

    // Omnisexual flag

    push 0xf2
    call draw

    push 0xe2
    call draw

    push 0x24
    call draw

    push 0x26
    call draw

    push 0x6f
    call draw

    // Clear screen and restart

    call wait

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    push 0x00
    call draw

    jmp e9

draw:
    pop e11
    pop r2
    mov r3, r2

    mov r4, 40 // Num chars per row
    mov r5, 0

    mov e10, pc
    nop

    mov r1, 1 
    int 1
    inc r5

    cmp r6, r5, r4
    jz r6, e10

    ret

wait:
    pop e11
    mov r1, 2500
    int 2
    ret
