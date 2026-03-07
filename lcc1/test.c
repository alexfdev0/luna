extern void print(char* string);
int a = 1;
int b = 1;


void _start() {
    asm ("mov sp, 0xefff");
    if (a == b) {
        print("hi\n");
    } else {
        print("no\n");
    }
    asm ("hlt");
}
