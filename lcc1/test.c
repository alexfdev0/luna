char* string = "Hello world!";
extern void print(char* str);

void _start() {
    asm ("mov sp, 0xEFFF");
    print(string);
halt:
    goto halt;
}
