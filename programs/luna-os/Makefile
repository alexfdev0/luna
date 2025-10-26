CC=lcc
all:
	$(CC) boot/boot.s -o boot.bin
	$(CC) kernel/kernel_stub.s kernel/kernel.c kernel/stdlib.s -o kernel.bin
	dd if=boot.bin of=os.bin bs=512
	dd if=kernel.bin of=os.bin bs=512 seek=1 conv=notrunc
	rm kernel.bin
	rm boot.bin
