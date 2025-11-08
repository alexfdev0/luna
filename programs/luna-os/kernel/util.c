void pause() {
    puts32("Press any key to continue...\n\n", 255, 0);
    asm ("int 0x6");
    return;
}

void serve_close() {
    puts32("Server closed.\n", 255, 0);
    asm ("mov r1, 0x7001a644");
    asm ("mov r2, 0");
    asm ("str r1, r2");
    return;
}
