#pragma bits 16

extern void putc(char c, short short int color);
extern void sleep(int seconds);
extern void* flags_start;

void render_flags() {
    short short int* fptr = (short short int*) flags_start;
    while (*fptr != 0xFE) {
        for (int i = 0; i < 40; i++) {
            putc(0x20, *fptr);
        }
        fptr++;
        if (*fptr == 0x00) {
            sleep(1);
        }
    }
}

void _start() {
    asm ("mov sp, 0xEFFF");

    // Load the next sectors on the disk
    asm ("int 0x10");
    asm ("mov r2, r1");
    asm ("mov r1, 1");
    asm ("mov r3, r1");
    asm ("int 11");

    // Set up PIT
    
    // asm ("mov r1, 0x6FFF0008"); 0xFA41 for 16 bit
    asm ("mov r1, 0xFA41");

    asm ("mov r2, pit_nxt");
    asm ("str16 r1, r2");

    while (1) {
        render_flags();
    }
}
