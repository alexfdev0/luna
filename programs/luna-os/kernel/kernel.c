#pragma bits 32
#include "stdlib.h"

#ifndef __LCC__
    #error "LunaOS must be compiled with LCC (other compilers are not supported.)"
#endif

extern void setup();
extern void shell();

void _cstart() __attribute__((noreturn)) {
    puts32("LunaOS\n", 255, 0);
    puts32("Copyright (c) 2025 Alexander Flax\n\n", 255, 0);
    puts32(malloc(200), 255, 0); 
 
    if (&PASSBUF == 0x00) { 
        setup(); 
    }
enterpass:
    puts32("Password: ", 255, 0);
    readin(TEMPBUF, 1, 1);
    xor_cycle(TEMPBUF);
    if (strcmp(PASSBUF, TEMPBUF) == 1) {
        shell();
    } else {
        puts32("Password is incorrect.\n", 255, 0);
        goto enterpass;
    }
    asm ("hlt");
}

