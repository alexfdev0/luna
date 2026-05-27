#pragma bits 16

extern void putc(char c, short short int color);
extern void sleep(int seconds);
extern void* flags_start;

void render_flags() {
    short short int* fptr = flags_start;
    while (*fptr != 0xFE) {
        for (int i = 0; i < 40; i = i + 1) {
            putc(0x20, *fptr);
        }
        fptr = fptr + 1;
        if (*fptr == 0x00) {
            sleep(1);
        }
    }
}

void _start() {
    asm ("mov sp, 0xEFFF");

    // Load the next sectors on the disk
    
    asm ("mov r1, 1"); // sector 1
    asm ("mov r3, r1");
    asm ("mov r2, 0");
    asm ("int 11");
    asm ("mov r1, 2"); // sector 2
    asm ("mov r3, r1");
    asm ("mov r2, 0");
    asm ("int 11");

    // Set up PIT
    
    // asm ("mov r1, 0x6FFF0008"); 0xFA41 for 16 bit
    asm ("mov r1, 0xFA41");

    asm ("mov r2, pit_nxt");
    asm ("strf r1, r2");

    while (1) {
        render_flags();
    }
}
