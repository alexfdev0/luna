#pragma bits 32

#include "stdlib.h"
#include "util.h"

long int* ffnt(char* filename) {
    short short int* buffer = (short short int*) malloc((long int) strlen(filename));
    long int* bufptr = (long int*) buffer;

    long int seen = 0;
    while (*filename != 0) { 
        if (*filename != 0x20) {
            putchar((char) *filename, (char*) buffer); 
            buffer = buffer + 1;
        } else {
            if (seen == 0) {
                seen = 1;
                putchar(46, (char*) buffer); // .
                buffer = buffer + 1;
            }
        }
        filename = filename + 1;
    }
    
    return bufptr;
}

void fcreate(char* name, long int size) {
    // Load next file pointer
    long int** nfl = (long int**) 0x61C;


    **nfl = 0x4C465346; // Store file header
    *nfl = *nfl + 4;

    strcpy(name, (char*) *nfl); // Transfer name to file
    long int name_len = (long int) strlen(name);
    *nfl = *nfl + name_len; 
    **nfl = size;
    *nfl = *nfl + 4;

    for (long int i = 0; i < size; i = i + 1) {
        **nfl = 0x00;
        *nfl = *nfl + 1;
    }

    long int sector = (long int) *nfl / 512;
    save_sector(sector);

    *nfl = *nfl + size;
    save_sector(3);
}

long int* find_file(char* name) {
    long int* fsp = (long int*) 0x618;
    long int* fp = *fsp;

    while (1) {
        if (*fp != 0x4C465346) {
            break;
        }
        // skip over header
        fp = fp + 4;

        if (strcmp(name, (char*) fp) == 1) {
            fp = fp + 20; // skip over name and size
            return fp;
        } else {
            fp = fp + 16; // skip over name
            long int size = (long int) *fp;
            fp = fp + 4; // skip over size marker
            fp = fp + size; // skip over file contents 
        }
    }

    return 0;
}

long int* fopen(char* filename, short short int complain_on_not_found) {
    long int* faddr = find_file(filename);
    if (faddr == 0x00000000) {
        if (complain_on_not_found) {
            puts32("File '", 0xA0, 0);
            puts32((char*) ffnt(filename), 0xA0, 0);
            puts32("' not found!\n", 0xA0, 0);
            return 0;
        }
    }

    return faddr;
}

long int fgetsize(char* filename) {
    long int* fptr = fopen(filename, 1);
    fptr = fptr - 4;
    return *fptr;
}

void flist() {
    long int* fsp = (long int*) 0x618;
    long int* fp = *fsp;

    while (1) {
        if (*fp != 0x4C465346) {
            break;
        }
        // skip over header
        fp = fp + 4;
        puts32((char*) ffnt((char*) fp), 255, 0);
        puts32("\n", 255, 0); 
        fp = fp + 16; // skip over name
        long int size = (long int) *fp;
        fp = fp + 4; // skip over size marker
        fp = fp + size; // skip over file contents 
    }

    return 0;
}

void fwrite(char* name, char* content) {
    long int* cptr = fopen(name, 1);
    if (cptr == 0) { 
        return;
    }

    strcpy(content, (char*) cptr);
    long int sec = (long int) cptr / 512;
    save_sector(sec);
}

asm (".dword 0x4a41414d");
