## Documentation
[Jump to preamble](#preamble)<br>
[Jump to registers](#registers)<br>
[Jump to instructions](#instructions)<br>
[Jump to interrupts](#interrupts)<br>
[Jump to assembly](#assembly)<br>
[Jump to linking](#linking)<br>
[Jump to frontend](#frontend)<br><br>

## Preamble
The Luna L2 is a simple, lightweight, RISC CPU that aims to be clean while also leveraging some luxuries from CISC, with the ultimate end goal of being easy to teach and learn.<br><br>

## Registers
L2 has 31 total registers for the storage and manipulation of data and information:<br><br>
R0-R12: general purpose registers, all can be written to and read from<br>
E0-E12: extra registers, also general purpose. The standard calling convention uses registers E0-E6<br>
SP: stack pointer<br>
PC: program counter/instruction pointer<br>
IRV: interrupt return address storage
IR: interrupt register
B: bank register.

## Instructions
The Luna L2 has 29 unique instructions that allow the CPU to interact with registers, memory, and the BIOS<br><br>

1. MOV: moves a value from the source to the destination; source can be register or immediate.<br>
2. HLT: stops the CPU from executing instructions.<br>
3. JMP: sets the program counter to the specified address; address can be register or immediate.<br>
4. INT: calls a BIOS interrupt. [Jump to interrupts](#interrupts)<br> 
5. JNZ: sets the program counter to the specified address if the register is not zero; address can be immediate or register.<br>
6. NOP: stalls the CPU for 1 cycle.<br>
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
21. NOT: performs bitwise NOT on two registers and puts the result to a register.<br> 
22. XOR: performs bitwise XOR on two registers and puts the result to a register.<br>
23. LOD: loads a byte from memory to a register.<br>
24. STR: stores a value to a memory address from a register. (bytewise)<br>
25. STRF: stores 2 bytes (16-bit mode) or 4 bytes(32-bit mode) from a register to memory.<br>
26. LODF: loads 2 bytes (16-bit mode) or 4 bytes (32-bit mode) from memory to a register.<br>
27. SET: if the operand is 00, it sets the CPU to 16 bit mode. If the operand is 01, it sets the CPU to 32 bit mode.<br>
28. SHL: Shifts the value in a register to the left by a certain value.<br>
29. SHR: Shifts the value in a register to the right by a certain value.<br>
(bytewise: scheme where storing a register value to memory stores the low byte in `address` and the high byte in `address + 1`; in 32 bit mode, it's from `address` to `address + 3`)<br><br>

## Interrupts
Luna L2 has several instructions necessary for operation and/or communicating to other devices, listed below:<br><br>

1. Print character to the screen (char in R1, foreground in R2, background in R3)<br>
2. Reserved for the programmable interval timer<br>
3. Unmapped<br>
4. Unmapped<br>
5. Reserved for the keyboard<br>
6. Unmapped<br>
7. Reserved for illegal instruction trap<br>
8. Unmapped<br>
9. Unmapped<br>
10. Memory query (returns in R1)<br>
11. Load sector from disk (sector number in R1, drive in R2)<br>
12. Set video cursor (X in R1, Y in R2)<br>
13. Write a sector to disk (sector number in R1, drive in R2)<br>
14. Load video cursor (returns X in R1, Y in R2)<br>
15. Reset the machine<br>
16. Drive query (returns in R1)<br>
17. Shut down machine<br><br>

## Assembly
The L2 architecture has a custom assembler (`las`) to convert programs from assembly language (.asm, .s, .S) to machine code (.o) that can then be linked and then run on L2.<br>
# Syntax specifications
The syntax of L2 assembly is similar to that of Intel assembly syntax. An instruction consists of a mnemonic, then the operands. Above, there were no specifications on which instructions use which registers, since every instruction that uses registers can use any register.<br>
Except for LOD, STR, LODF, and STRF, the destination register is always the first register in the instruction.<br>
You can use a label name followed by a colon to make a label, which gets turned into a numerical offset at assembly time. Therefore you can treat them as numbers as well. These can also be used as functions with `call` and `ret`.<br>
# Custom directives
There are some directives in LAS that do not correspond to any instruction on L2. They are as follows:<br>
`call <label>`: calculates the return address, pushes it onto the stack, and jumps to the label specified (`call mylabel`)<br>
`ret`: jumps to the value in register `re1`<br>
`.ascii <string>`: defines a sequence of ASCII bytes, wrapped in quotation marks<br>
`.asciz <string>`: defines a sequence of ASCII bytes, wrapped in quotation marks (null terminated)<br>
`.word <number>`: defines a 2-byte constant<br>
`.double <number>`: defines a 4-byte constant<br>
`.ptr <number or label>`: defines a number with the width of the current mode; memory locations can be used with this.<br>
`.global <label>`: exposes a symbol to other object files.<br>
`.bits <16/32>` changes the mode of the assembler to the specified mode (does not change CPU; use `SET`)<br>
`.embed <file>`: includes a file in the object code the assembler at that location.<br>
`.org <location>`: tells the linker to calculate all offsets with respect to the origin.<br>
`.fill <number>`: Tells the linker to fill the resulting binary until it reaches the size specified.<br>
`.pad <number>`: inserts the specified number of null bytes at that location.<br>
`.db <bytes, seperated by a comma>`: inserts the specified bytes at that location until a newline.<br>
# Examples
`mov r1, 5` (destination: r1, source: 5)<br>
`pop r1` (destination: r1)<br>
`add r3, r1, r2` (destination: r3, source 1: r1, source 2: r2)<br>
`mylabel:
    mov r1, 1
    ret`<br>
`call mylabel`<br>
`.ascii "Hello world!"`<br>
`.global mysymbol` (exposes 'mysymbol' to other object code)<br>
# Assembling a program
To assemble a program, use the following: `las <flags> <input file(s)> -o <output file>`<br>
The flags are as follows:<br>
`-v`: shows the version of LAS and exits.<br>
`-c`: do not invoke linker (`l2ld`) after assembly is complete.<br>
Note: you may also use the Luna Compiler Collection frontend (`lcc`) with the same syntax to do this.<br><br>

## Linking
Linking is the process of resolving offsets in an object file (.o) to turn it into an executable file (.bin). L2 also has a custom linker (`l2ld`) to convert L2 object format (L2O) to L2 executable format (L2E).<br>
# Linking a program
To link a program, use the following: `l2ld <flags> <input file(s)> -o <output file>`<br>
The flags are as follows:<br>
`-v`: shows the version of L2LD and exits.<br>
Note: you may also use the Luna Compiler Collection frontend (`lcc`) with the same syntax to do this.<br><br>

## Frontend
L2 has a simple compiler frontend (`lcc`) which is similar to that of `gcc` or `clang`. LCC can also automatically detect file type and use the relevant programs to compile it to an executable.<br>
# Compiling a program
To compile a program, use the following: `lcc <flags> <input file(s)> -o <output file>`<br>
The flags are as follows:<br>
`-c`: do not invoke linker (`l2ld`) after assembly is complete.<br>
`-v`: shows the version of LCC and exits.<br>
`-S`: do not invoke assembler (`las`) after compilation is complete.<br>
Supported file types: (subject to change)<br>
`.s`: assembly<br>
`.S`: assembly<br>
`.asm`: assembly<br>
`.o`: object file<br>
`.c`: C
`.h`: C
