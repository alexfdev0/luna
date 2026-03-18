#pragma bits 32

void play_sound(void* buffer, long int size, short short int block) {
    short short int* done_flag = 0x7000FA09;
    *done_flag = 0;

    long int* buf_ptr = 0x7000FA01;
    *buf_ptr = buffer;

    long int* size_ptr = 0x7000FA05;
    *size_ptr = size;

    short short int* play_flag = 0x7000FA00;
    *play_flag = 1;

    if (block) {
    psb_beginning:
        if (*done_flag) {
            goto psb_end;
        }
        goto psb_beginning;
    }
psb_end:
    return;
}
