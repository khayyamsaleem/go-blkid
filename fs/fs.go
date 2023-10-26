package fs

/*
#cgo CFLAGS: -I${SRCDIR}/../memblkid/include
#cgo LDFLAGS: -lblkid
#include "../memblkid/src/memblkid.c"
#include "memblkid.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// SupportedFilesystems lists all filesystems supported for probing by blkid
func SupportedFilesystems() []string {
	cFilesystems := C.list_supported_filesystems()
	var filesystems []string

	i := 0
	for {
		current := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cFilesystems)) + uintptr(i)*unsafe.Sizeof(cFilesystems)))
		if current == nil {
			break
		}
		filesystems = append(filesystems, C.GoString(current))
		i++
	}

	C.free_supported_filesystems() // Ensure you free the memory used by C

	return filesystems
}

// GetFilesystemTypeFromBuffer gets the filesystem type from a buffer
func GetFilesystemTypeFromBuffer(data []byte) (string, error) {
	cData := (*C.char)(unsafe.Pointer(&data[0]))
	cProbe := C.create_probe_from_buffer(cData, C.size_t(len(data)))

	if cProbe == nil {
		return "", errors.New("failed to create probe from buffer")
	}

	fsType := C.get_filesystem_type(cProbe)
	if fsType == nil {
		return "", errors.New("failed to determine fs type from probe")
	}

	return C.GoString(fsType), nil
}
