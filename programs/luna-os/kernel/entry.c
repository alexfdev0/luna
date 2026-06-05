#pragma bits 32

#include "stdlib.h"
#include "shell.h"
#include "util.h"
#include "lufs.h"

#ifndef __LCC__
    #error "LunaOS must be compiled with LCC (other compilers are not supported.)"
#endif

asm (".global enterpass");

void _cstart() __attribute__((noreturn)) {
    if (fopen("NOTEPAD     SYS", 0) == 0x00000000) {
        fcreate("NOTEPAD     SYS", 256);
    }

    puts32("Welcome to ", 0xff, 0);
    puts32("Luna", 0x9b, 0);
    puts32("OS!\n", 0xff, 0);
    
    while (1) {
        shell();
    }
}

