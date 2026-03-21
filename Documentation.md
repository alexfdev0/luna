# Documentation
[Jump to preamble](#preamble)<br>
[Jump to registers](#registers)<br>
[Jump to instructions](#instructions)<br>
[Jump to interrupts](#interrupts)<br>
[Jump to assembly](#assembly)<br>
[Jump to frontend](#frontend)<br><br>

# Preamble
The Luna L2 is a simple, lightweight, RISC CPU that aims to be clean while also leveraging some luxuries from CISC, with the ultimate end goal of being easy to teach and learn.<br><br>

# Registers
L2 has 33 total registers for the storage and manipulation of data and information:<br><br>
R0-R12: general purpose registers, all can be written to and read from<br>
E0-E12: extra registers, also general purpose. The standard calling convention uses registers E0-E6<br>
E13-E14: Assembler reserved for PIE mode and PIE macros, do not use<br>
SP: stack pointer<br>
PC: program counter/instruction pointer<br>
IRV: interrupt return address storage<br>
IR: interrupt register<br>
B: bank register.<br>

# Instructions
The Luna L2 has 29 unique instructions that allow the CPU to interact with registers, memory, and the BIOS<br><br>

1. MOV: moves a value from the source to the destination; source can be register or immediate.<br>
2. HLT: stops the CPU from executing instructions.<br>
3. JMP: sets the program counter to the specified address; address can be register or immediate.<br>
4. INT: calls a BIOS interrupt. [Jump to interrupts](#interrupts)<br> 
5. JNZ: sets the program counter to the specified address if the register is not zero; address can be immediate or register.<br>
6. NOP: no-op; does nothing.<br>
7. CMP: sets the specified register to 1 if the other two registers are the same; otherwise sets to 0.<br>
8. JZ: sets the program counter to the specified address if the register is zero.<br>
9. INC: increments a register by 1.<br>
10. DEC: decrements a register by 1.<br>
11. PUSH: Pushes a word to the stack and increments the stack pointer by 2; word can be in a register or an immediate.<br>
12. POP: Pops a word off the stack to the specified register and decrements the stack pointer by 2.<br>
13. ADD: Puts the sum of 2 registers into a register.<br>
14. SUB: Puts the subtraction result of 2 registers into a register.<br>
15. MUL: Puts the product of 2 registers into a register.<br>
16. DIV: Puts the quotient of 2 registers into a register.<br>
17. IGT: Sets a register to 1 if the second register is greater than the third register; otherwise 0.<br>
18. ILT: Sets a register to 1 if the second register is less than the third register; otherwise 0.<br>
19. AND: performs bitwise AND on two registers and puts the result to a register.<br>
20. OR:  performs bitwise OR on two registers and puts the result to a register.<br>
21. NOT: performs bitwise NOT on a registers and puts the result to a register.<br> 
22. XOR: performs bitwise XOR on two registers and puts the result to a register.<br>
23. LOD: loads a byte from memory to a register.<br>
24. STR: stores a value to a memory address from a register. (bytewise)<br>
25. STRF: stores 2 bytes (16-bit mode) or 4 bytes(32-bit mode) from a register to memory.<br>
26. LODF: loads 2 bytes (16-bit mode) or 4 bytes (32-bit mode) from memory to a register.<br>
27. SET: if the operand is 00, it sets the CPU to 16 bit mode. If the operand is 01, it sets the CPU to 32 bit mode.<br>
28. SHL: Shifts the value in a register to the left by a certain value.<br>
29. SHR: Shifts the value in a register to the right by a certain value.<br>
(bytewise: scheme where storing a register value to memory stores the low byte in `address` and the high byte in `address + 1`; in 32 bit mode, it's from `address` to `address + 3`)<br><br>

# Interrupts
Luna L2 has several instructions necessary for operation and/or communicating to other devices, listed below:<br><br>

1. Print character to the screen (char in R1, foreground in R2, background in R3)<br>
2. Reserved for the programmable interval timer<br>
3. Returns drive attached/inserted status in R1 (drive in R1)<br>
4. Syscall reserved; (if using PIE and host OS supports it; refer to your OS' syscall implementation)<br>
5. Reserved for the keyboard<br>
6. Unmapped<br>
7. Reserved for illegal instruction trap<br>
8. Unmapped<br>
9. Unmapped<br>
10. Memory query (returns in R1)<br>
11. Load sector from disk (destination sector number in R1, drive in R2, real sector number in R3)<br>
12. Set video cursor (X in R1, Y in R2)<br>
13. Write a sector to disk (RAM sector number in R1, drive in R2, real sector number in R3)<br>
14. Load video cursor (returns X in R1, Y in R2)<br>
15. Reset the machine<br>
16. Drive query (returns in R1)<br>
17. Shut down machine<br><br>

# Memory Map
The standard Luna L2 memory map is as follows:<br>
## 16-bit mode
0x0000 - 0xEFFF: GP RAM<br>
0xF000 - 0xF009: Audio controller registers<br>
0xFA0A - 0xFA11: Mouse registers<br>
0xFA12: Keyboard register<br>
0xFA13 - 0xFA1A: PIT registers<br>
0xFA1B - 0xFA30: Network controller registers<br>
0xFA31 - 0xFA36: RTC registers<br>
0xFA37 - 0xFC36: IDT
0xFE00 - 0xFFFF (switch to higher/lower in banks of 512 bytes using B register; addresses that exceed VRAM will write/read from the final byte in VRAM): VRAM<br><br>

## 32-bit mode
0x00000000 - 0x6EFFFFFF: GP RAM<br>
0x6FFF0000 - 0x6FFFFFFF: IDT<br>
0x70000000 - 0x7000F9FF: VRAM<br>
0x7000FA00 - 0x7000FA09: Audio controller registers<br>
0x7000FA0A - 0x7000FA11: Mouse registers<br>
0x7000FA12: Keyboard register<br>
0x7000FA13 - 0x7000FA1A: PIT registers<br>
0x7001A644 - 0x7001A659: Network controller registers<br>
0x7001B65E - 0x7001B663: RTC registers<br><br>

## Device registers
### VRAM
VRAM controls the screen graphics; 1 byte in VRAM is equivalent to 1 pixel on screen with the RGB332 color scheme; this makes for a 320x200 resolution with 8 bits per pixel for 64,000 bytes of VRAM.<br>
The display refreshes at 70 hertz.<br>
### Audio controller
Byte 0: Play flag; commands the audio controller to play sound based on given parameters.<br>
Bytes 1-4: 32-bit pointer of audio start (reads from mapper so MMIO works); tells the audio controller where to start from.<br>
Bytes 5-8: 32-bit size of audio; gets added to the start pointer to determine where to stop playing audio.<br>
Byte 9: done flag; set to 1 when the audio controller is done playing audio.<br>
### Mouse registers
Bytes 0-3: 32-bit X position of mouse.<br>
Bytes 4-7: 32-bit Y position of mouse.<br>
### Keyboard register
Byte 0: character code of last key pressed.<br>
### PIT registers
Bytes 0-3: 32-bit programmed countdown value (PIT updates every millisecond)
Bytes 4-7: 32-bit actual countdown value (when this reaches 0, the PIT interrupt will be triggered and it will be reset to the programmed countdown value)<br>
### Network controller registers
All modes:<br>
Byte 0: execute register; executes a function based on parameters.<br>
Byte 1: mode; 0x00: TCP client; 0x01: TCP server.<br>
Bytes 6-7: 16 bit port number (for both modes).<br>
Client mode:<br>
Bytes 2-5: IP address; stored as number.<br>
Bytes 14-17: pointer to data to send.<br>
Bytes 18-21: size of data to send.<br>
Bytes 10-13: output pointer (outputs are always 2048 bytes).<br>
Server mode:<br>
Byte 2: Accept flag<br>
Byte 4: Connection waiting flag<br>
Byte 5: Command register; 0: close connection; 1: read from network buffer (2048 bytes); 2: write to network buffer<br>
Bytes 10-13: output pointer<br>
Bytes 14-17: pointer to data to send.<br>
Bytes 18-21: size of data to send.<br>
### RTC registers
Byte 0: Current second<br>
Byte 1: Current minute<br>
Byte 2: Current hour<br>
Byte 3: Current day<br>
Byte 4: Current month<br>
Byte 5: Current year minus 2000<br>

# Assembly
The Luna toolchain has a custom assembler (`las`) to convert programs from assembly language to object format that can then be linked and then run on L2. (Flags can be found in the [frontend](#frontend) section.)<br>
## Syntax specifications
The syntax of L2 assembly is similar to that of Intel assembly syntax. An instruction consists of a mnemonic, then the operands. Above, there were no specifications on which instructions use which registers, since every instruction that uses registers can use any register.<br>
Except for LOD/LODE, STR/STRE, LODF/LODFE, and STRF/STRFE, the destination register is always the first register in the instruction.<br>
You can use a label name followed by a colon to make a label, which gets turned into a numerical offset at assembly time. Therefore you can treat them as numbers as well. These can also be used as functions with `call` and `ret`.<br>
## Custom directives
There are some directives in LAS that do not correspond to any instruction on L2. They are as follows:<br>
`call <label>`: calculates the return address, pushes it onto the stack, and jumps to the label specified<br>
`ret`: jumps to the value in register `e11`<br>
`lode/stre/lodfe/strfe`: Same as their normal equivalents (`lod/str/lodf/strf`) but they factor in the effective address (in `e14`); useful if you are making PIEs.<br>
`pusha`: Pushes all GP registers from R0 upwards.<br>
`popa`: Pops all GP registers from E12 downwards.<br>
`.ascii <string>`: defines a sequence of ASCII bytes, wrapped in quotation marks<br>
`.asciz <string>`: defines a sequence of ASCII bytes, wrapped in quotation marks (null terminated)<br>
`.word <number>`: defines a 2-byte constant<br>
`.dword <number>`: defines a 4-byte constant<br>
`.ptr <number or label>`: defines a number with the width of the current mode; memory locations can be used with this.<br>
`.global <label>`: exposes a symbol to other object files.<br>
`.bits <16/32>` changes the mode of the assembler to the specified mode (does not change CPU; use `SET`)<br>
`.embed <file>`: includes a file in the object code the assembler at that location.<br>
`.org <location>`: tells the linker to calculate all offsets with respect to the origin.<br>
`.fill <number>`: Tells the linker to fill the resulting binary until it reaches the size specified.<br>
`.pad <number>`: inserts the specified number of null bytes at that location.<br>
`.byte <bytes, seperated by a comma>`: inserts the specified bytes at that location until a newline.<br><br>
## Examples
For assembly examples, I recommend you see the many assembly files in `programs/luna-os/` as they are much better than the examples below.<br><br>
`mov r1, 5` (destination: r1, source: 5)<br>
`pop r1` (destination: r1)<br>
`add r3, r1, r2` (destination: r3, source 1: r1, source 2: r2)<br>
`mylabel:
    mov r1, 1
    ret`<br>
`call mylabel`<br>
`.ascii "Hello world!"`<br>
`.global mysymbol` (exposes 'mysymbol' to other object code)<br>

# Linking
The Luna toolchain has a custom linker (`l2ld`) that will convert files from object format to executable format. Flags for L2LD can be found below<br><br>

# C compilation
The Luna toolchain has a custom C compiler (`lcc1`) that will convert from C (C99) to assembly. Flags for LCC1 can be found below.<br>
Please note that LCC1 is very incomplete in terms of C features supported but an effort is being made to add new features.<br>
## Attributes
In LCC1, you can use `__attribute__((<attribute(s)>))` to specify attributes for variables in the manner below:<br>
`int foo __attribute__((require_const)) = 5;`<br>
`int foo() __attribute__((noreturn)) {}`<br>
All valid attributes are listed below:<br>
`noreturn`: tells compiler to not pop a return address from the stack at the beginning of a function or adding a `ret` at the end of a function. (valid for functions)<br>
`norename`: tells compiler to not rename the function in question to `main` if the function name is `_start`. (valid for functions)<br>
`require_const`: tells compiler to require that the non-global variable in question have a constant compile time initializer. (valid for variables)<br><br>
## Extensions
In LCC1, there are a few extensions to make programming easier. They are as follows:<br>
`short short <int>`: Specifies an 8-bit integer without a custom type like `uint8_t`
`__embed__`: Embeds a file into the resulting assembly file; equivalent to `.embed` but from C. Syntax: `__embed__ <pre (puts at top of file instead at canonical location)/static (no auto .global)> <label name> (("<file path>"))`<br>
`void*` dereferences: equivalent to `char`/`short short int` dereferences; grabs 8 bits.<br><br>

# Frontend
L2 has a simple compiler frontend (`lcc`) which is similar to that of GCC or Clang. LCC will also automatically detect each file's type and use the relevant subtool for it.<br>
## Compiling a program
To compile a program, use the following: `lcc <flags> <input file(s)> -o <output file>`<br>
## Flags
`-o`: specifies output file (L2LD)<br>
`-Werror`: upgrades all warnings to errors (LAS/LCC1)<br>
`-fstdlib`: Allows linker to pull from `/usr/local/lib/l2ld/` if the corresponding label is not present (L2LD)<br>
`-fpie`: Specifies to compile in position independent executable mode (LCC1/LAS/L2LD)<br>
`-fpie-16`: Specifies to L2LD to link your PIE in 16 bit mode (L2LD)<br>
`-fpie-32`: Specifies to L2LD to link your PIE in 32 bit mode (L2LD)<br>
`-c`: do not invoke linker (`l2ld`) after assembly is complete (LAS)<br>
`-v`: shows version of LCC as well as all command invocations<br>
`-S`: do not invoke assembler (`las`) after compilation is complete (LCC1)<br>
## Supported file types
`.s/.S/.asm`: assembly<br>
`.o/.obj`: object file<br>
`.c/.h/.cpp/.hpp/.cxx/.hxx`: C (C++ is not supported; extra extensions are present for compatibility.)<br>
