#pragma bits 16
extern void print(char* string);
int a = 512;
int b = sizeof(a);

void _start() {
    print(*a);
halt:
    goto halt;
}
