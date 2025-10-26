void pause() {
    asm ("push e11");
    puts32("Press any key to continue...\n\n", 255, 0);
    asm ("int 0x6");
    asm ("pop e11");
    return;
}
