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

type ProbeOption func(pr C.blkid_probe) error

func WithFilesystemFilter(filesystems ...string) ProbeOption {
	return func(pr C.blkid_probe) error {
		cfilesystems := make([]*C.char, len(filesystems)+1)
		for i, fs := range filesystems {
			cfilesystems[i] = C.CString(fs)
		}
		defer func() {
			for _, fs := range cfilesystems {
				C.free(unsafe.Pointer(fs))
			}
		}()

		if C.blkid_probe_filter_superblocks_type(pr, C.BLKID_FLTR_ONLYIN, &cfilesystems[0]) < 0 {
			return errors.New("failed to set filesystem filter")
		}
		return nil
	}
}

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
func GetFilesystemTypeFromBuffer(data []byte, options ...ProbeOption) (string, error) {
	cData := (*C.char)(unsafe.Pointer(&data[0]))
	pr := C.create_probe_from_buffer(cData, C.size_t(len(data)))
	if pr == nil {
		return "", errors.New("failed to create probe from buffer")
	}
	defer C.blkid_free_probe(pr)

	for _, opt := range options {
		if err := opt(pr); err != nil {
			return "", err
		}
	}

	fsType := C.get_filesystem_type(pr)
	if fsType == nil {
		return "", errors.New("failed to determine fs type from probe")
	}

	res := C.GoString(fsType)
	C.free(unsafe.Pointer(fsType))

	return res, nil
}

// GetBlockSizeFromBuffer gets the filesystem block size from a buffer
func GetBlockSizeFromBuffer(data []byte, options ...ProbeOption) (uint64, error) {
    cData := (*C.char)(unsafe.Pointer(&data[0]))
    pr := C.create_probe_from_buffer(cData, C.size_t(len(data)))
    if pr == nil {
        return 0, errors.New("failed to create probe from buffer")
    }
    defer C.blkid_free_probe(pr)

    for _, opt := range options {
        if err := opt(pr); err != nil {
            return 0, err
        }
    }

    blocksize := C.get_blocksize(pr)
    if blocksize < 0 {
        return 0, errors.New("failed to determine block size from probe")
    }

    return uint64(blocksize), nil
}

