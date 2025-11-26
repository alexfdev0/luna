asm (".bits 32");

#include "stdlib.h"
#include "audio.h"
#include "textlibs.h"
#include "lufs.h"

#include "util.c"
#include "setup.c"
#include "shell.c"


void _cstart() __attribute__((noreturn)) { 
    puts32("LunaOS\n", 255, 0);
    puts32("Copyright (c) 2025 Alexander Flax\n\n", 255, 0); 
 
    if (checkpass()) { 
         
    } else {
        setup();
    }
enterpass:
    puts32("Password: ", 255, 0);
    readin(TEMPBUF, 1);
    xor_cycle(TEMPBUF);
    if (strcmp(PASSBUF, TEMPBUF)) {
        // play_sound(STARTUP_SOUND, 205130, 0);
        shell();
    } else {
        puts32("Password is incorrect.\n", 255, 0);
        goto enterpass;
    }
    asm ("hlt");
}

