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
        puts32((char*) PROMPTBUF, COLOR_GRAY, COLOR_BLACK);
        readin((char*) buf, 1, 0); 
        if (strcmp("reboot", (char*) buf) == 1) {
            puts32("Rebooting...", COLOR_GRAY, COLOR_BLACK);
            asm ("mov r1, 0");
            asm ("int 0xf"); 
        }

        if (strcmp("about", (char*) buf) == 1) {
            puts32("LunaOS 1.0.0\nBy Alexander Flax\n", COLOR_GRAY, COLOR_BLACK);
            puts32("Network adapter: ", COLOR_GRAY, COLOR_BLACK);
            puts32((char*) 0x7001A65A, COLOR_GRAY, COLOR_BLACK);

            puts32("\n\n", COLOR_GRAY, COLOR_BLACK);
            continue;
        }

        if (strcmp("promptedit", (char*) buf) == 1) {
            puts32("Enter terminal prompt: ", COLOR_GRAY, COLOR_BLACK);
            readin((char*) PROMPTBUF, 0, 0);
            save_buffer((char*) PROMPTBUF, 0);

            puts32("\n", COLOR_GRAY, COLOR_BLACK);
            continue;
        }
        
        if (strcmp("notepad", (char*) buf) == 1) {
            long int size = fgetsize((char*) fntf(notepad_file));
            long int* buf = malloc(size); 
            long int* file = fopen((char*) fntf(notepad_file), 1);
            if (file == 0x00000000) {
                continue;
            }

            strcpy((char*) file, (char*) buf);

            readin((char*) buf, 0, 0);
            fwrite((char*) fntf(notepad_file), (char*) buf);

            puts32("\n", COLOR_GRAY, COLOR_BLACK);
            continue;
        }

        if (strcmp("files", (char*) buf) == 1) {
            flist();
            puts32("\n", COLOR_GRAY, COLOR_BLACK);
            continue;
        }

        if (strcmp("shutdown", (char*) buf) == 1) {
            puts32("Shutting down...\n", COLOR_GRAY, COLOR_BLACK);
            asm ("int 0x11");
        } 

        if (strcmp("bc", (char*) buf) == 1) {
            save_graphics_buf();
            render_buf(BAYACHAO_IMG);
            wait_for_key();
            render_buf((void*) 0x30303030);

            continue;
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

        puts32("Bad command '", COLOR_LRED, COLOR_BLACK);
        puts32((char*) buf, COLOR_LRED, COLOR_BLACK);
        puts32("'\n", COLOR_LRED, COLOR_BLACK);
    }
    return;
}
