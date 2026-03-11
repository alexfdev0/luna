#include "stdio.h"

char* string = "Hello world";
int* a = 100;

int foo(char* string) {
    printf("0x%x", string);
}

int main() {
    foo(string);
    return 0;
}
