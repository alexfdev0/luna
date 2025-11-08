_start:
    set 32
    .bits 32
   
    // IP
    mov r1, 0x496ADCFF
    mov r2, 0x7001A646
    strf r2, r1

    // PORT
    mov r1, 0x0B
    mov r2, 0x7001A64a
    str r2, r1

    mov r1, 0xB8
    mov r2, 0x7001A64b
    str r2, r1

    // TIMEOUT
    mov r1, 0x00
    mov r2, 0x7001A64c
    str r2, r1

    mov r1, 0x0A
    mov r2, 0x7001A64d
    str r2, r1

    // TCP/UDP FLAG
    mov r1, 0x00
    mov r2, 0x7001A645
    str r2, r1

    // MESSAGE
    mov r1, 0x48
    mov r2, 0x7001A64e
    str r2, r1

    mov r1, 0x49
    mov r2, 0x7001A64f
    str r2, r1

    // INIT
    mov r1, 0x01
    mov r2, 0x7001A644
    str r2, r1
