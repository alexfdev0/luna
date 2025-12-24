#include "bayachao.h"

void shell() {
top:
    puts32(PROMPTBUF, 255, 0);
    readin(TEMPBUF, 1, 0); 
    if (strcmp("reboot", TEMPBUF)) {
        puts32("Rebooting...", 255, 0);
        // play_sound(SHUTDOWN_SOUND, 207748, 1);
        // sleep(500);
        asm ("mov r1, 0");
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
        readin(PROMPTBUF, 0, 0);
        save_buffer(PROMPTBUF, 0);

        puts32("\n", 255, 0);
        goto top;
    }
    
    if (strcmp("notepad", TEMPBUF)) {
        readin(TEMPBUF, 0, 0);
        lufs_write_file("NOTEPAD SYS     ", TEMPBUF);

        puts32("\n", 255, 0);
        goto top;
    }

    if (strcmp("passwd", TEMPBUF)) {
        puts32("Enter old password: ", 255, 0);
        readin(TEMPBUF, 1, 1);
        xor_cycle(TEMPBUF);
        if (strcmp(TEMPBUF, PASSBUF)) {
            puts32("Enter new password: ", 255, 0);
            readin(PASSBUF, 1, 1);
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

    if (strcmp("shutdown", TEMPBUF)) {
        puts32("Shutting down...\n", 255, 0);
        // play_sound(SHUTDOWN_SOUND, 207748, 1);
        // sleep(500);
        asm ("int 0x11");
    }

    if (strcmp("ws", TEMPBUF)) {
        puts32("Listening on port 3000", 255, 0);
        serve(3000);
        await:
        serve_await_connection(500);
        serve_read();
        serve_write("HTTP/1.1\r\nContent-Type: text/html\r\nServer: FurNet\r\n\r\n<!DOCTYPE html>\r\n<html>\r\n<body>\r\n<h1>Hello world!</h1>\r\n</body>\r\n</html>\r\n\r\n");
        serve_connection_close();
        goto await;
    }

    if (strcmp("bc", TEMPBUF)) {
        save_graphics_buf();
        render_buf(BAYACHAO_IMG);
        wait_for_key();
        render_buf(GBUF);

        goto top;
    }

    if (strcmp("testfault", TEMPBUF)) {
        asm ("jmp 0x7001A644");
    }

    if (strcmp("clear", TEMPBUF)) {
        render_buf(GBUF_EMPTY);
        goto top;
    } 

    puts32("'", 255, 0);
    puts32(TEMPBUF, 255, 0);
    puts32("' is not recognized as an internal or external command.\n", 255, 0);
    goto top;
    return;
}
