.bits 16
.global flags_start

flags_start:
rainbow_flag:
    .byte 0xe0
    .byte 0xf0
    .byte 0xfc
    .byte 0x10
    .byte 0x0b
    .byte 0x62

    .byte 0

trans_flag:
    .byte 0x5b
    .byte 0xf6
    .byte 0xff
    .byte 0xf6
    .byte 0x5b

    .byte 0

asexual_flag:
    .byte 0x24
    .byte 0xb6
    .byte 0xff
    .byte 0x81

    .byte 0

gay_flag:
    .byte 0x11
    .byte 0x56
    .byte 0xba
    .byte 0xff
    .byte 0x77
    .byte 0x4a
    .byte 0x45
    
    .byte 0

lesbian_flag:
    .byte 0xc4
    .byte 0xf1
    .byte 0xff
    .byte 0xce
    .byte 0xa1
    
    .byte 0

bisexual_flag:
    .word 0xc1c1
    .word 0x6a6a
    .word 0x0a0a

    .byte 0

pan_flag:
    .word 0xe6e6
    .word 0xf8f8
    .word 0x5757
    
    .byte 0

genderfluid_flag:
    .byte 0xee
    .byte 0xff
    .byte 0xc3
    .byte 0x24
    .byte 0x27
    
    .byte 0

aromantic_flag:
    .byte 0x14
    .byte 0x99
    .byte 0xff
    .byte 0xb6
    .byte 0x24
    
    .byte 0

agender_flag:
    .byte 0x24
    .byte 0xbb
    .byte 0xff
    .byte 0xbd
    .byte 0xff
    .byte 0xbb
    .byte 0x24

    .byte 0

nonbinary_flag:
    .byte 0xfc
    .byte 0xff
    .byte 0xab
    .byte 0x24

    .byte 0

polysexual_flag:
    .word 0xe2e2
    .word 0x1919
    .word 0x1313

    .byte 0

omnisexual_flag:
    .byte 0xf2
    .byte 0xe2
    .byte 0x24
    .byte 0x26
    .byte 0x6f

    .byte 0

.byte 0xfe
