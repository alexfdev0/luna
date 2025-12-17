void user_setup() {
    puts32("Setup will now initialize the user that will use this machine.\n\n", 255, 0);
    puts32("Username: root\n", 255, 0);
    puts32("Enter a password: ", 255, 0);
    readin(PASSBUF, 1, 1);
    xor_cycle(PASSBUF);
    puts32("\n\n", 255, 0); 
    return;
}

void setup() { 
    puts32("Setup LunaOS\n", 255, 0);
    puts32("Welcome to LunaOS!\n\nThis interactive setup will guide you through\nthe process of setting up LunaOS\non your computer.\n\n", 255, 0);
    pause(); 

    puts32("Detecting your drive...\n", 255, 0);
    lufs_create_file("NOTEPAD SYS     ", 256); // Create notepad file
    if (getdrive()) {
        puts32("Setup has detected you are running Setup on a USB device.\n\n", 255, 0);
        user_setup();

        puts32("Setup will copy the operating system from the USB to the hard disk.\n\n", 255, 0);
        pause();
        
        puts32("Setup is copying files to the hard disk... ", 255, 0);
        setup_copy(0x155);
        puts32("done.\n\n", 255, 0);

        puts32("Setup will now restart this machine to complete the setup process.\n", 255, 0);
        asm ("int 0x6");
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
