#pragma bits 32

asm (".global play_sound_loc");
asm ("play_sound_loc:");

void play_sound(void* buffer, long int size, short short int block) {
    short short int* done_flag = (short short int*) 0x80000009;
    *done_flag = 0;

    *(long int*) 0x80000001 = buffer;
    *(long int*) 0x80000005 = size;
    *(short short int*) 0x80000000 = 1; 

    if (block) {
        while (*done_flag == 0) {}
    }
    return;
}
