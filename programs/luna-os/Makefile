CC=lcc
all:
	$(CC) boot/boot.s -o boot.bin
	$(CC) kernel/kernel_stub.s kernel/kernel.c kernel/stdlib.s kernel/audio.s -o kernel.bin
	dd if=boot.bin of=os.bin bs=512
	dd if=kernel.bin of=os.bin bs=512 seek=1 conv=notrunc
	printf '\033[32mSectors to load: 0x%X\033[0m\n' $$(( $$(stat -f%z kernel.bin) / 512 ))
	rm kernel.bin
	rm boot.bin
