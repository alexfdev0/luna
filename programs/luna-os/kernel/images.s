.bits 32

.global BOOT_IMG
.global BAYACHAO_IMG

BAYACHAO_IMG:
    .embed "kernel/images/bayachao.raw"

BOOT_IMG:
    .embed "kernel/images/boot.332"
