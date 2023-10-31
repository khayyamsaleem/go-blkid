#ifndef MEMBLKID_H
#define MEMBLKID_H

#include <blkid/blkid.h>

blkid_probe create_probe_from_buffer(char *buf, size_t size);
const char* get_filesystem_type(blkid_probe pr);
char** list_supported_filesystems(void);
void free_supported_filesystems(void);
int get_blocksize(blkid_probe pr);

#endif // MEMBLKID_H

