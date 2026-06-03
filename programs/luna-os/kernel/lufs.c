#pragma bits 32
#include "stdlib.h"

void lufs_create_file(char* name, long int size) {
    // asm ("hlt");
    // asm ("hlt");
    // Load next file pointer
    long int* nfp = 0x61C;
    long int* nfl = *nfp;

    *nfl = 0x4C465346;; // Store file header
    nfl = nfl + 4;

    strcpy(name, nfl); // Transfer name to file
    long int name_len = strlen(name);
    nfl = nfl + name_len; 
    *nfl = size;
    nfl = nfl + 4;

    for (long int i = 0; i < size; i = i + 1) {
        *nfl = 0xDE;;

        nfl = nfl + 1;
    }

    long int sector = nfl / 512;
    save_sector(sector); 
}

void lufs_write_file() {
    save_sector(1);
}
