.bits 16
.global load_executable
var_15:
    .dword 0x00000000

load_executable:
    pop e11
    push e11
    call ASLR_generate_address
    mov r4, e6
    mov r7, var_15
    strf r7, r4
    mov r5, 2
    mov r1, r5
    mov r7, r1
    push r7
    mov r1, var_15
    lodf r1, r2
    mov r7, r2
    push r7
    mov r5, 0
    mov r1, r5
    mov r7, r1
    push r7
    call load_sector
    mov r4, e6
    pop e11
    ret
