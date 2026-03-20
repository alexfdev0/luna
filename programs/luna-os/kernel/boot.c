#pragma bits 32

#include "stdlib.h"
#include "images.h"
#include "audio.h"
#include "stub.h"

void boot() __attribute__((noreturn)) {
    targeted_load(BOOT_SOUND, 226);
    targeted_load(BOOT_IMG, 126);
    targeted_load(play_sound_loc, 3);
    targeted_load(renderbuf_loc, 2);
    targeted_load(sleep_loc, 2);

    play_sound(BOOT_SOUND, 115308, 0);
    render_buf(BOOT_IMG);

    boot_load_all_sectors(0x337);

    sleep(5);

    render_buf(0x30303030);
    asm ("jmp _cstart");
}
