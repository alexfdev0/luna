#pragma bits 32
#include "stdlib.h"
#include "lufs.h"
#include "util.h"

void user_setup() {
    puts32("Setup will now initialize the user that will use this machine.\n\n", 255, 0);
    puts32("Username: root\n", 255, 0);
    puts32("Enter a password: ", 255, 0);
    readin(PASSBUF, 1, 1);
    xor_cycle(PASSBUF);
    puts32("\n\n", 255, 0);

    if (*(short short int*) 0x204 == 0x9A) {
        puts32("You have a compatible boot menu!\n", 0b00011100, 0);
    }

    // puts32("Enable drive encryption?\n\nY: Yes\nN: No\n");
    // wait_for_key();
    // handle_drive_encryption();
    return;
}

void setup() { 
    puts32("Setup LunaOS\n", 255, 0);
    puts32("Welcome to LunaOS!\n\nThis interactive setup will guide you through\nthe process of setting up LunaOS\non your computer.\n\n", 255, 0);
    pause(); 

    puts32("Detecting your drive...\n", 255, 0);
    lufs_create_file("NOTEPAD SYS     ", 256); // Create notepad file
    if (getdrive()) {
        if (query_drive_inserted(0) == 0) {
            puts32("\n\nError! ", 0xA0, 0);
            puts32("LunaOS cannot access your\nhard drive.\nPlease insert it and restart.\n\nPress any key to reboot...", 255, 0);
            wait_for_key();
            reboot();
        }
        puts32("Setup has detected you are running Setup on a USB device.\n\n", 255, 0);
        user_setup();

        puts32("Setup will copy the operating system from the USB to the hard disk.\n\n", 255, 0);
        pause();
        
        puts32("Setup is copying files to the hard disk... ", 255, 0);
        setup_copy(0x1D4, 0);
        puts32("\nSetup has completed copying files to the hard disk.\n\n", 255, 0);

        puts32("Setup will now restart this machine to complete the setup process.\n", 255, 0);
        pause(); 
        asm ("mov r1, 0");
        asm ("int 0xf");
    } else {
        puts32("Setup has detected you are running Setup on a hard disk.\n\n", 255, 0);
        user_setup();
        
        puts32("Setup will now restart this machine to complete the setup process.\n", 255, 0);
        save_buffer(PASSBUF, 0);
        asm ("int 0x6");
        asm ("mov r1, 0");
        asm ("int 0xf");
    }
    return;
}
