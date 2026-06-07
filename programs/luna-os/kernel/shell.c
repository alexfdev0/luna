#pragma bits 32

#include "images.h"
#include "stdlib.h"
#include "lufs.h"
#include "util.h"
#include "audio.h"

char* notepad_file = "NOTEPAD     SYS";

void shell() {
    while (1) {
        long int* buf = malloc(256);
        puts32(PROMPTBUF, 255, 0);
        readin(buf, 1, 0); 
        if (strcmp("reboot", buf) == 1) {
            puts32("Rebooting...", 255, 0);
            asm ("mov r1, 0");
            asm ("int 0xf"); 
        }

        if (strcmp("about", buf) == 1) {
            puts32("LunaOS 1.0.0\nBy Alexander Flax\n", 255, 0);
            puts32("Network adapter: ", 255, 0);
            puts32(0x7001A65A, 255, 0);

            puts32("\n\n", 255, 0);
            continue;
        }

        if (strcmp("promptedit", buf) == 1) {
            puts32("Enter terminal prompt: ", 255, 0);
            readin(PROMPTBUF, 0, 0);
            save_buffer(PROMPTBUF, 0);

            puts32("\n", 255, 0);
            continue;
        }
        
        if (strcmp("notepad", buf) == 1) {
            long int size = fgetsize(notepad_file);
            long int* buf = malloc(size); 
            long int* file = fopen(notepad_file, 1);
            if (file == 0x00000000) {
                continue;
            }

            strcpy(file, buf);

            readin(buf, 0, 0);
            fwrite(notepad_file, buf);

            puts32("\n", 255, 0);
            continue;
        }

        if (strcmp("files", buf) == 1) {
            flist();
            puts32("\n", 255, 0);
            continue;
        }

        if (strcmp("shutdown", buf) == 1) {
            puts32("Shutting down...\n", 255, 0);
            asm ("int 0x11");
        } 

        if (strcmp("bc", buf) == 1) {
            save_graphics_buf();
            render_buf(BAYACHAO_IMG);
            wait_for_key();
            render_buf(0x30303030);

            continue;
        }

        if (strcmp("testfault", buf) == 1) { 
            asm ("mov r1, 4");
            asm ("mov r2, pc");
            asm ("int 0x07");
        }
        
        if (strcmp("clear", buf) == 1) {
            render_buf(0x40404040);
            video_set_cursor(0, 0);
            continue;
        } 

        if (strcmp("exec", buf) == 1) {
            load_executable();
            continue;
        }

        puts32("Bad command '", 0xA0, 0);
        puts32(buf, 0xA0, 0);
        puts32("'\n", 0xA0, 0);
    }
    return;
}
