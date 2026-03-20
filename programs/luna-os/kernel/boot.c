#pragma bits 32

#include "stdlib.h"
#include "images.h"
#include "audio.h"

void boot() __attribute__((noreturn)) {
    play_sound(BOOT_SOUND, 115308, 0);
    render_buf(BOOT_IMG);

    sleep(5);

    render_buf(0x30303030);
    asm ("jmp _cstart");
}
