.bits 32

.global play_sound
.global SHUTDOWN_SOUND
.global STARTUP_SOUND

play_sound:
    pop e11
    pop r5 // BLOCK
    pop r2 // SIZE
    pop r1 // BUFFER

    mov r3, 0x7000FA09
    mov r4, 0
    str r3, r4
    
    mov r3, 0x7000FA01
    strf r3, r1
 
    mov r3, 0x7000FA05
    strf r3, r2

    mov r4, 1
    mov r3, 0x7000FA00
    str r3, r4

    jnz r5, play_sound_block

    ret
play_sound_block:
    mov r3, 0x7000FA09
    mov r1, 15

    mov e10, pc
    nop
    lod r3, r4
    int 2
    jz r4, e10

    mov r4, 0
    str r3, r4
    ret

STARTUP_SOUND:
    .embed kernel/audio/startup.pcm

SHUTDOWN_SOUND:
    .embed kernel/audio/shutdown.pcm
