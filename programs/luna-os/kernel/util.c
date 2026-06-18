#pragma bits 32

#include "stdlib.h"
#include "audio.h"

void tohex(long int number, short short int capitalized) {
    puts32("0x", 255, 0);
    puts32((char*) itoa(number, capitalized, (char*) malloc(11)), 255, 0);
    puts32("\n", 255, 0);
    free(11);
}

short short int pause() {
    puts32("Press any key to continue...\n\n", COLOR_WHITE, COLOR_BLACK);
    short short int code = wait_for_key();
    return code;
}

void kernel_panic() __attribute__((noreturn)) {
    asm ("push r2");
    asm ("push r1");

    play_sound(CRASH_SOUND, 164352, 0);
    screen_fill(0x80808080);
    puts32("System error\n\nYour PC ran into an error and needs to\nbe restarted.\n\nPress any key to reboot.\n\n\n", COLOR_WHITE, COLOR_RED);
    
    puts32("Instruction: 0x", COLOR_WHITE, COLOR_RED);
    asm ("pop e9");
    puts32((char*) itoa(_e9, 1, (char*) malloc(11)), COLOR_WHITE, COLOR_RED);
    puts32("\n", COLOR_WHITE, COLOR_RED);

    puts32("Location: 0x", COLOR_WHITE, COLOR_RED);
    asm ("pop e9");
    puts32((char*) itoa(_e9, 1, (char*) malloc(11)), COLOR_WHITE, COLOR_RED);
    puts32("\n", COLOR_WHITE, COLOR_RED);


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
    return (int) _e12;
}

void reboot() {
    asm ("int 0x10");
    asm ("int 0xf");
}

void load_sector(short short int drive, long int* dest_sector, long int real_sector) {
    asm ("mov r2, e0");
    asm ("mov r1, e1");
    asm ("mov r3, e2");
    asm ("int 0x0b");
}

void load_executable() {
    if (query_drive_inserted(2) == 0) {
        puts32("Error! ", 0xA0, 0);
        puts32("Please insert a disc into the DVD\ndrive and try again.\n", COLOR_WHITE, COLOR_BLACK);
        return;
    }

    long int* address = (long int*) ASLR_generate_address();
    address = address + 1;

    load_sector(2, address / 512, 0);
    load_sector(2, address / 512 + 1, 1);
    load_sector(2, address / 512 + 2, 2);

    if (*address != 0x4C325049) {   
        puts32("Error! ", COLOR_WHITE, COLOR_BLACK);
        puts32("Invalid executable file format.\n", COLOR_WHITE, COLOR_BLACK);
        return;
    }
    lexec_core((long int) address);
}

void app_error() __attribute__((noreturn)) {
    puts32("Error! ", COLOR_LRED, COLOR_BLACK);
    puts32("Executable automatically\nterminated due to instruction fault.\n", COLOR_WHITE, COLOR_BLACK);
    goto lexec_done; 
}

short short int* get_first_word(char* string) {
    short short int* buffer = (short short int*) malloc(1024);
    short short int* ogbuf = buffer;

    while (*string != 0x00) {
        if (*string == 0x20) {
            break;
        }
        putchar(*string, (char*) buffer);
        buffer++;
        string++;
    }
    free(1024);
    return ogbuf;
}
