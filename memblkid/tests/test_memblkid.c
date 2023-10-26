#include <stdio.h>
#include <stdlib.h>
#include <dirent.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include "memblkid.h"

#define DEVICE_PATH "/dev/"
#define BUFFER_SIZE 65536 // 64 KiB

void check_device_filesystem_types() {
    DIR *dir;
    struct dirent *entry;
    char path[512];
    char buffer[BUFFER_SIZE];
    int fd;

    dir = opendir(DEVICE_PATH);
    if (!dir) {
        perror("Failed to open /dev");
        exit(1);
    }

    while ((entry = readdir(dir)) != NULL) {
        if (entry->d_type != DT_BLK) { // Skipping non-block devices
            continue;
        }

        if (strncmp(entry->d_name, "nvme", 4) != 0) {
            continue;
        }

        snprintf(path, sizeof(path), "%s%s", DEVICE_PATH, entry->d_name);

        fd = open(path, O_RDONLY);
        if (fd == -1) {
            perror("Failed to open device");
            continue;
        }

        ssize_t read_size = read(fd, buffer, BUFFER_SIZE);
        if (read_size != BUFFER_SIZE) {
            perror("Failed to read 64KiB from device");
            close(fd);
            continue;
        }
        close(fd);

        blkid_probe pr = create_probe_from_buffer(buffer, read_size);
        if (!pr) {
            printf("Failed to create probe for device %s\n", path);
            continue;
        }

        const char *fs_type = get_filesystem_type(pr);
        if (fs_type) {
            printf("Device %s has filesystem type: %s\n", path, fs_type);
        } else {
            printf("Device %s has an unknown or no filesystem.\n", path);
        }
        free((void*)fs_type);
    }

    closedir(dir);
}

void list_all_supported_filesystems() {
    char** filesystems = list_supported_filesystems();
    printf("Supported Filesystems:\n");

    int i = 0;
    while (filesystems[i] != NULL) {
        printf("- %s\n", filesystems[i]);
        i++;
    }

    free_supported_filesystems();
}

int main() {
    check_device_filesystem_types();
    list_all_supported_filesystems();

    return 0;
}

