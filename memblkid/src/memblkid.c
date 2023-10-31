#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/syscall.h>
#include <linux/memfd.h>
#include <blkid/blkid.h>

#define MAX_FILESYSTEMS 100

// Create blkid probe from a buffer
blkid_probe create_probe_from_buffer(char *buf, size_t size) {
	int fd;
	blkid_probe pr = NULL;

	fd = syscall(SYS_memfd_create, "memfd_blkid_buf", MFD_CLOEXEC);
	if (fd == -1) {
		perror("Failed to create memfd for buffer");
		return NULL;
	}

	ssize_t written = write(fd, buf, size);
	if (written == -1) {
		perror("Failed to write buf to memfd");
		return NULL;
	}

	lseek(fd, 0, SEEK_SET);

	pr = blkid_new_probe();
	if (!pr) {
		close(fd);
		perror("Could not create new blkid probe.");
		return NULL;
	}

	if (blkid_probe_set_device(pr, fd, 0, 0)) {
		close(fd);
		blkid_free_probe(pr);
		perror("Unable to point blkid probe fd to in-memory buf");
		return NULL;
	}

	return pr;
}

// Perform the probe and extract filesystem type
const char* get_filesystem_type(blkid_probe pr) {
	int probe_result;
	const char *fstype = NULL;
	char *result = NULL;

	if (!pr) {
		perror("Null probe provided.");
		return NULL;
	}

	blkid_reset_probe(pr);
	probe_result = blkid_do_probe(pr);
	if (probe_result < 0) {
		perror("Unable to probe device.");
	} else {
		if (blkid_probe_has_value(pr, "TYPE")) {
			blkid_probe_lookup_value(pr, "TYPE", &fstype, NULL);
			if (fstype) {
				result = strdup(fstype); // required so blkid_free_probe doesn't throw out my important stuff!
			}
		}
	}

	return result;
}

// Perform the probe and extract blocksize
int get_blocksize(blkid_probe pr) {
	int probe_result;
	const char* value = NULL;
	int result = -1;

	if (!pr) {
		perror("Null probe provided.");
		return -1;
	}

	blkid_reset_probe(pr);
	probe_result = blkid_do_probe(pr);
	if (probe_result < 0) {
		perror("Unable to probe device.");
	} else {
		if (blkid_probe_has_value(pr, "BLOCK_SIZE")) {
			blkid_probe_lookup_value(pr, "BLOCK_SIZE", &value, NULL);
			if (value) {
				result = atoi(value);
			}
		}
	}

	return result;

}

// assuming at most 100 supported filesystems
static char* supported_fs[MAX_FILESYSTEMS + 1] = { NULL };
static int fs_cache_initialized = 0;

// list supported filesystems
char** list_supported_filesystems(void) {
	if (fs_cache_initialized) {
		return supported_fs;
	}

	int idx = 0;
	const char* name = NULL;

	while (blkid_superblocks_get_name(idx++, &name, NULL) == 0 && idx < MAX_FILESYSTEMS) {
		supported_fs[idx-1] = strdup(name);
	}

	supported_fs[idx] = NULL; // NULL terminate the list
	fs_cache_initialized = 1;

	return supported_fs;
}

// free supported fs cached values (optional)
void free_supported_filesystems(void) {
	for (int i = 0; i < MAX_FILESYSTEMS && supported_fs[i] != NULL; i++) {
		free(supported_fs[i]);
		supported_fs[i] = NULL;
	}
	fs_cache_initialized = 0;
}

