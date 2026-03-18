#pragma bits 32
#include "stdlib.h"
#include "audio.h"

void tohex(long int number, short short int capitalized) {
    puts32("0x", 255, 0);
    puts32(itoa(number, capitalized, malloc(11)), 255, 0);
    puts32("\n", 255, 0);
}

void pause() {
    puts32("Press any key to continue...\n\n", 255, 0);
    wait_for_key();
    return;
}

void serve_close() {
    puts32("Server closed.\n", 255, 0);
    short short int* NETWORK_COMMAND = 0x7001A644;
    *NETWORK_COMMAND = 0x00;
    return;
}

void kernel_panic() __attribute__((noreturn)) {
    asm ("push r2");
    asm ("push r1");

    play_sound(CRASH_SOUND, 164352, 0);
    screen_fill(0xA0A0A0A0);
    puts32("System error\n\nYour PC ran into an error and needs to\nbe restarted.\n\nPress any key to reboot.\n\n\n", 255, 0xA0);
    
    puts32("Instruction: 0x", 255, 0xA0);
    asm ("pop e9");
    puts32(itoa(_e9, 1, malloc(11)), 255, 0xA0);
    puts32("\n", 255, 0xA0);

    puts32("Location: 0x", 255, 0xA0);
    asm ("pop e9");
    puts32(itoa(_e9, 1, malloc(11)), 255, 0xA0);
    puts32("\n", 255, 0xA0);


    wait_for_key();
    asm ("int 0x10");
    asm ("int 0xf");
}

void video_set_cursor(int x, int y) {
    // Arguments in e0, e1
    asm ("mov r1, e0");
    asm ("mov r2, e1");
    asm ("int 0x0c");
}

int query_drive_inserted(short short int drive) {
    asm ("mov r1, e0"); // Move drive number to r1
    asm ("int 0x3"); // Query drive inserted
    asm ("mov e12, r1");
    return _e12;
}

void reboot() {
    asm ("int 0x10");
    asm ("int 0xf");
}

void load_sector(short short int drive, long int dest_sector, long int real_sector) {
    asm ("mov r2, e0");
    asm ("mov r1, e1");
    asm ("mov r3, e2");
    asm ("int 0x0b");
}
