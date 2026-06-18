#pragma bits 32

#include "images.h"
#include "stdlib.h"
#include "lufs.h"
#include "util.h"
#include "audio.h"

char* notepad_file = "NOTEPAD.SYS";

void shell() {
    while (1) {
        long int* buf = malloc(256);
        puts32((char*) PROMPTBUF, COLOR_WHITE, COLOR_BLACK);
        readin((char*) buf, 1, 0);
        if (strcmp("reboot", (char*) buf) == 1) {
            puts32("Rebooting...", COLOR_WHITE, COLOR_BLACK);
            asm ("mov r1, 0");
            asm ("int 0xf"); 
        }

        if (strcmp("about", (char*) buf) == 1) {
            puts32("LunaOS 2.0.0\nBy Alexander Flax\n", COLOR_WHITE, COLOR_BLACK);
            puts32("\n\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("promptedit", (char*) buf) == 1) {
            puts32("Enter terminal prompt: ", COLOR_WHITE, COLOR_BLACK);
            readin((char*) PROMPTBUF, 0, 0);
            save_buffer((char*) PROMPTBUF, 0);

            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }
        
        if (strcmp("notepad", (char*) buf) == 1) {
            long int size = fgetsize((char*) fntf(notepad_file));
            long int* buf = malloc(size); 
            File* f = fopen((char*) fntf(notepad_file), 1);
            long int* file = f->Address;
            if (file == 0x00000000) {
                continue;
            }

            strcpy((char*) file, (char*) buf);

            readin((char*) buf, 0, 0);
            fwrite((char*) fntf(notepad_file), (char*) buf);

            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("files", (char*) buf) == 1) {
            flist();
            puts32("\n", COLOR_WHITE, COLOR_BLACK);
            continue;
        }

        if (strcmp("shutdown", (char*) buf) == 1) {
            puts32("Shutting down...\n", COLOR_WHITE, COLOR_BLACK);
            asm ("int 0x11");
        } 

        if (strcmp("testfault", (char*) buf) == 1) { 
            asm ("mov r1, 4");
            asm ("mov r2, pc");
            asm ("int 0x07");
        }
        
        if (strcmp("clear", (char*) buf) == 1) {
            render_buf((void*) 0x40404040);
            video_set_cursor(0, 0);
            continue;
        } 

        if (strcmp("exec", (char*) buf) == 1) {
            load_executable();
            continue;
        }

        if (strcmp("battery", (char*) buf) == 1) {
            puts32("Battery level: ", COLOR_WHITE, COLOR_BLACK);
            short short int* bat_ptr = 0x80000026;
            tohex((long int) *bat_ptr, 1);
            continue;
        }

        puts32("Bad command '", COLOR_LRED, COLOR_BLACK);
        puts32((char*) get_first_word((char*) buf), COLOR_LRED, COLOR_BLACK);
        puts32("'\n", COLOR_LRED, COLOR_BLACK);
    }
    return;
}
