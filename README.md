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

# Windows Support<br>
- Windows is not currently supported because the Windows API isn't friendly to the shared library model Luna L2 uses.<br><br>

# Running a disk image<br>
- Run `luna-l2 <disk image>` to execute a disk image.<br>
- Note: there are 3 hardware slots for a disk image:<br>
HDD: Hard disk drive; use `-hdd <file>` or just `<file>` to insert a file into the slot.<br>
SD: Secure Digital card/USB; use `-sd <file>` to insert a file into the slot.<br>
DVD: Optical disc slot; use `-dvd <file>` to insert a file into the slot.<br>
- To customize which one you want to boot from, you can use `-boot <hdd/sd/dvd>`<br>

# Emulator parameters<br>
- `-ram <RAM amount in bytes>` - adjust the amount of RAM available to the emulator<br>
- `-gpu <gpu device>` - change the graphics processor device (defaults to `g1x`); see the list below for valid equipment.<br>
- `-apu <apu device>` - change the audio processor devices (defaults to `s1`); see the list below for valid equipment.<br><br>

# Equipment list<br>
## Graphics<br>
### Luna G1<br>
- Emulator equipment ID: `g1`<br>
- Valid resolutions: 320x200@8bpp<br>
- Video memory: 64 KB (64,000 bytes)<br>
- Extra features: None<br><br>
### Luna G1X<br>
- Emulator equipment ID: `g1x`<br>
- Valid resolutions: 320x200@8bpp<br>
- Video memory: 64 KB (64,000 bytes)<br>
- Extra features: Simple transparency via color code 0xE3<br><br>
## Audio<br>
### Luna S1<br>
- Emulator equipment ID: `s1`<br>
- Extra features: None<br><br>
