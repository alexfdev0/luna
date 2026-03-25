#pragma bits 32

asm (".global play_sound_loc");
asm ("play_sound_loc:");

void play_sound(void* buffer, long int size, short short int block) {
    short short int* done_flag = 0x7000FA09;
    *done_flag = 0;

    *(long int*) 0x7000FA01 = buffer;
    *(long int*) 0x7000FA05 = size;
    *(short short int*) 0x7000FA00 = 1;

    if (block) {
        while (*done_flag == 0) {}
    }
    return;
}
