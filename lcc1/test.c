#pragma bits 16
extern void print(char* string);

char a = "A";
char b = "L";
char c = "E";
char d = "X";
int e;

void _start() {
    asm ("mov sp, 0xEFFF");
    print(&a);
halt:
    goto halt;
}
