# Luna<br>
A simple, lightweight RISC CPU architecture.<br><br>

# Requirements<br>
- MacOS, Linux, or FreeBSD<br>
- Go (if you are compiling manually)<br>
- GCC (if you are compiling manually)

# Automatic Installation (MacOS (amd64/arm64))<br>
- Download the relevant installer from the latest release.<br>
- Run the installer and go through installation steps.<br>
- Note: All toolchain applications will automatically be added to your system's PATH variable.<br><br>

# Manual Installation (MacOS, Linux, FreeBSD)<br>
- Clone the repository using `git clone`<br>
- Navigate into the directory<br>
- Run `make; make install` to install the Luna L2 emulator and toolchain<br>
- Run `luna-l2 <disk image>` to run an application<br>
- Note: if you would like to install the legacy Luna L1 emulator and toolchain, run `make legacy` to install it as well as the assembler and C compiler. Then run `luna-l1 <disk image>` to run an application<br>

# Windows Support<br>
- Windows is not currently supported because the Windows API isn't friendly to the shared library model Luna L2 uses.<br><br>

# Running a disk image<br>
- Run `luna-l2 <disk image>` to execute a disk image.<br>
- Note: there are 3 hardware slots for a disk image:<br>
HDD: Hard disk drive; use `-hdd <file>` or just `<file>` to insert a file into the slot.<br>
SD: Secure Digital card/USB; use `-sd <file>` to insert a file into the slot.<br>
DVD: Optical disc slot; use `-dvd <file>` to insert a file into the slot.<br>
- To customize which one you want to boot from, you can use `-boot <hdd/sd/dvd>`<br>
