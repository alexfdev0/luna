#pragma bits 32
#include "stdlib.h"

void lufs_create_file(char* name, long int size) {
    // asm ("hlt");
    // asm ("hlt");
    // Load next file pointer
    long int* nfp = 0x61C;
    long int* nfl = *nfp;

    *nfl = 0x4C465346; // Store file header
    nfl = nfl + 4;

    strcpy(name, nfl); // Transfer name to file
    long int name_len = strlen(name);
    nfl = nfl + name_len; 
    *nfl = size;
    nfl = nfl + 4;

    for (long int i = 0; i < size; i = i + 1) {
        *nfl = 0x00;
        nfl = nfl + 1;
    }

    long int sector = nfl / 512;
    save_sector(sector);

    *nfp = *nfp + size;
    save_sector(3);
}

long int* lufs_find_file(char* name) {
    long int* fsp = 0x618;
    long int* fp = *fsp;

    while (1) {
        if (*fp != 0x4C465346) {
            break;
        }
        // skip over header
        fp = fp + 4;
        // name portion
        if (strcmp(name, fp) == 1) {
            fp = fp + 20; // skip over name and size
            return fp;
        } else {
            fp = fp + 16; // skip over name
            long int size = *fp;
            fp = fp + 4; // skip over size marker
            fp = fp + size; // skip over file contents 
        }
    }

    return 0;
}

void lufs_write_file(char* name, char* content) {
    long int* cptr = lufs_find_file(name);
    if (cptr == 0) {
        puts32("File '", 0xA0, 0);
        puts32(name, 0xA0, 0);
        puts32("' not found!\n", 0xA0, 0);
        return;
    }
    strcpy(content, cptr);
    long int sec = cptr / 512;
    save_sector(sec);
}
