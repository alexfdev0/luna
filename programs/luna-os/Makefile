CC=lcc
all:
	$(CC) boot/stage1.s -o stage1.bin
	$(CC) boot/stage2.s -o stage2.bin
	$(CC) kernel/kernel_stub.s kernel/kernel.c kernel/stdlib.s kernel/audio.s kernel/lufs.s kernel/bayachao.s -o kernel.bin
	dd if=stage1.bin of=os.img bs=512
	dd if=stage2.bin of=os.img bs=1 seek=492 conv=notrunc
	dd if=kernel.bin of=os.img bs=512 seek=3 conv=notrunc
	mkfs.lufs -d os.img -s 0x5F5E1000 -n "Luna" -hst 1536 -fst 0x2C000 -rebuild-bs
	- printf '\033[32mSectors to load: 0x%X\033[0m\n' $$(( $$(stat -c%s kernel.bin) / 512 + 1 ))
	- printf '\033[32mSectors to load: 0x%X\033[0m\n' $$(( $$(stat -f%z kernel.bin) / 512 + 1 ))
	rm kernel.bin
	rm stage1.bin
	rm stage2.bin
