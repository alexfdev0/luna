#pragma bits 16
extern void print(char* string);

int a = 101;
int b = 100;

int retone() {
    return 1;
}

void _start() {
    asm ("mov sp, 0xEFFF");
    if (retone()) {
        print("yes");
    } else {
        print("no");
    }
halt:
    goto halt; 
}
