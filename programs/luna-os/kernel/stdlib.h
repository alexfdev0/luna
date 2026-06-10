extern void readin(char* buffer, short short int clrdone, short short int blind);
extern void puts32(char* str, short short int fg, short short int bg);
extern void setup_copy(int sectors, short short int drive);
extern void save_buffer(char* buffer, short short int drive);
extern void render(char* buffer);
extern void screen_fill(long int col);
extern void render_buf(void* img_buf);
extern void save_graphics_buf();
extern short short int wait_for_key();
extern int strcmp(char* buf1, char* buf2);
extern int strlen(char* str);
extern long int* malloc(long int size);
extern void free(long int size);
extern char* itoa(long int num, short short int capitalized, char* location);
extern int ASLR_generate_address();
extern void lexec_core(long int address);
extern void sleep(long int seconds);
extern void strcpy(char* b1, char* b2);
extern void save_sector(long int sector);
extern void putchar(char c, char* dest);
extern void printchar(char c);
extern void* PROMPTBUF;
extern void* renderbuf_loc;
extern void* sleep_loc;

// Console colors, adapted from BIOS colors
#define COLOR_BLACK   0x00
#define COLOR_BLUE    0x02
#define COLOR_GREEN   0x10
#define COLOR_CYAN    0x12
#define COLOR_RED     0x80
#define COLOR_MAGENTA 0x82
#define COLOR_BROWN   0x88
#define COLOR_GRAY    0xFF
#define COLOR_LBLUE   0x4B
#define COLOR_LGREEN  0x5D
#define COLOR_LCYAN   0x5F
#define COLOR_LRED    0xE9
#define COLOR_LMAGENT 0xEB
#define COLOR_YELLOW  0xFD
#define COLOR_WHITE   0xFF
