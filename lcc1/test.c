#pragma bits 16
extern void print(char* s);

char* string = "Hello world!";
int a = 1;

void _start() {
    a = 2 * 2;
    asm ("mov sp, 0xEFFF");
    print(string);
halt:
    goto halt;
}
