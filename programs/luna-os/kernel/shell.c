#pragma bits 32

#include "images.h"
#include "stdlib.h"
#include "lufs.h"
#include "util.h"
#include "audio.h"

void shell() {
    while (1) {
        puts32(PROMPTBUF, 255, 0);
        readin(TEMPBUF, 1, 0); 
        if (strcmp("reboot", TEMPBUF) == 1) {
            puts32("Rebooting...", 255, 0);
            play_sound(SHUTDOWN_SOUND, 387436, 1);
            asm ("mov r1, 0");
            asm ("int 0xf"); 
        }

        if (strcmp("about", TEMPBUF) == 1) {
            puts32("LunaOS 1.0.0\nBy Alexander Flax\n", 255, 0);
            puts32("Network adapter: ", 255, 0);
            puts32(0x7001A65A, 255, 0);

            puts32("\n\n", 255, 0);
            continue;
        }

        if (strcmp("promptedit", TEMPBUF) == 1) {
            puts32("Enter terminal prompt: ", 255, 0);
            readin(PROMPTBUF, 0, 0);
            save_buffer(PROMPTBUF, 0);

            puts32("\n", 255, 0);
            continue;
        }
        
        if (strcmp("notepad", TEMPBUF) == 1) {
            readin(TEMPBUF, 0, 0);
            lufs_write_file("NOTEPAD SYS     ", TEMPBUF);

            puts32("\n", 255, 0);
            continue;
        }

        if (strcmp("passwd", TEMPBUF) == 1) {
            puts32("Enter old password: ", 255, 0);
            readin(TEMPBUF, 1, 1);
            xor_cycle(TEMPBUF);
            if (strcmp(TEMPBUF, PASSBUF) == 1) {
                puts32("Enter new password: ", 255, 0);
                readin(PASSBUF, 1, 1);
                xor_cycle(PASSBUF);
                save_buffer(PASSBUF, 0);
            } else {
                puts32("Password is incorrect.", 255, 0);
            }
            puts32("\n", 255, 0);
            continue;
        } 

        if (strcmp("shutdown", TEMPBUF) == 1) {
            puts32("Shutting down...\n", 255, 0);
            play_sound(SHUTDOWN_SOUND, 387436, 1);
            asm ("int 0x11");
        } 

        if (strcmp("bc", TEMPBUF) == 1) {
            save_graphics_buf();
            render_buf(BAYACHAO_IMG);
            wait_for_key();
            render_buf(0x30303030);

            continue;
        }

        if (strcmp("testfault", TEMPBUF) == 1) { 
            asm ("mov r1, 4");
            asm ("mov r2, pc");
            asm ("int 0x07");
        }
        
        if (strcmp("clear", TEMPBUF) == 1) {
            render_buf(0x40404040);
            video_set_cursor(0, 0);
            continue;
        }

        if (strcmp("logout", TEMPBUF) == 1) {
            render_buf(0x30303030);
            video_set_cursor(0, 0);
            goto enterpass;
        }

        if (strcmp("exec", TEMPBUF) == 1) {
            load_executable();
            continue;
        }

        puts32("'", 255, 0);
        puts32(TEMPBUF, 255, 0);
        puts32("' is not recognized as an internal or external command.\n", 255, 0);
    }
    return;
}
