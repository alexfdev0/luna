void shell() {
top:
    puts32(PROMPTBUF, 255, 0);
    readin(TEMPBUF, 1); 
    if (strcmp("reboot", TEMPBUF)) {
        puts32("Rebooting...", 255, 0);
        asm ("int 0xf");  
    }

    if (strcmp("about", TEMPBUF)) {
        puts32("LunaOS 1.0.0\nBy Alexander Flax\n", 255, 0);
        puts32("Network adapter: ", 255, 0);
        puts32(0x7001b64f, 255, 0);

        puts32("\n\n", 255, 0);
        goto top;
    }

    if (strcmp("promptedit", TEMPBUF)) {
        puts32("Enter terminal prompt: ", 255, 0);
        readin(PROMPTBUF, 0);
        save_buffer(PROMPTBUF, 0);

        puts32("\n", 255, 0);
        goto top;
    }
    
    if (strcmp("notepad", TEMPBUF)) {
        readin(FILE, 0);
        save_buffer(FILE, 0);

        puts32("\n", 255, 0);
        goto top;
    }

    if (strcmp("passwd", TEMPBUF)) {
        puts32("Enter old password: ", 255, 0);
        readin(TEMPBUF, 1);
        xor_cycle(TEMPBUF);
        if (strcmp(TEMPBUF, PASSBUF)) {
            puts32("Enter new password: ", 255, 0);
            readin(PASSBUF, 1);
            xor_cycle(PASSBUF);
            save_buffer(PASSBUF, 0);
        } else
            puts32("Password is incorrect.", 255, 0); 
        puts32("\n", 255, 0);
        goto top;
    }

    if (strcmp("imget", TEMPBUF)) {
        send("IMG_GET", 0x7F000001, 580, 500);
        render(NETBUF);
        
        puts32("\n", 255, 0);
        goto top;
    }

    puts32("'", 255, 0);
    puts32(TEMPBUF, 255, 0);
    puts32("' is not recognized as an internal or external command.\n", 255, 0);
    goto top;
    return;
}
