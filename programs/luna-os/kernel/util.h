extern void pause();
extern void video_set_cursor(int x, int y); 
extern void tohex(long int number, short short int capitalized);
extern int query_drive_inserted(short short int drive);
extern void reboot();
extern void load_sector(short short int drive, long int dest_sector, long int real_sector);
