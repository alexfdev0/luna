.bits 32

.global play_sound
.global SHUTDOWN_SOUND

play_sound:
    pop e11
    pop r7 // SIZE
    pop r10 // BUFFER

    mov r2, 44100 // Copy up to 44100 at a time
    mov r3, 0 // Copied so far
    mov r11, 0 // Total read so far
    mov r4, 0x7000FA00 // Audio buffer pointer

    mov e10, pc
    nop

    lodf r10, r5
    strf r4, r5

    inc r10
    inc r10
    inc r10
    inc r10

    inc r4
    inc r4
    inc r4
    inc r4

    inc r11
    inc r11
    inc r11
    inc r11

    inc r3
    inc r3
    inc r3
    inc r3

    cmp r6, r3, r2
    jz r6, e10

    int 9

    igt r8, r11, r7
    jz r8, play_audio_goback

    ret

play_audio_goback:
    mov r12, pc
    nop

    int 0x11
    jz r1, r12

    mov r3, 0
    mov r4, 0x7000FA00
    jmp e10

SHUTDOWN_SOUND:
    .embed kernel/audio/shutdown.raw
