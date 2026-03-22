.bits 32

.global BOOT_IMG
.global BAYACHAO_IMG

BAYACHAO_IMG:
    .embed "images/bayachao.raw"

BOOT_IMG:
    .embed "images/boot.332"
