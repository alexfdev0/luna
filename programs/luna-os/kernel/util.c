void pause() {
    puts32("Press any key to continue...\n\n", 255, 0);
    wait_for_key();
    return;
}

void serve_close() {
    puts32("Server closed.\n", 255, 0);
    asm ("mov r1, 0x7001a644");
    asm ("mov r2, 0");
    asm ("str r1, r2");
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
