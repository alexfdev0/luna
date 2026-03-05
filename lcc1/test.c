int arr[1];

void foo(int a) {
    asm ("mov r1, e0");
    asm ("mov r2, 255");
    asm ("mov r3, 0");
    asm ("int 1");
}

int main() {
    asm ("mov sp, 0xefff");
    arr[0] = 64 + 1;
    foo(arr[0]);
    arr[1] = 64 + 0;
halt:
    goto halt;
}
