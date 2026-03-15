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
    play_sound(CRASH_SOUND, 164352, 0);
    screen_fill(0xA0A0A0A0);
    puts32("System error\n\nYour PC ran into an error and needs to\nbe restarted.\n\nPress any key to reboot.\n", 255, 0xA0);
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
