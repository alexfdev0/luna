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

    puts32("Welcome to ", COLOR_GRAY, COLOR_BLACK);
    puts32("Luna", COLOR_LCYAN, COLOR_BLACK);
    puts32("OS!\n", COLOR_WHITE, COLOR_BLACK);
    
    while (1) {
        shell();
    }
}

