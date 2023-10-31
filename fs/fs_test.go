package fs

import (
	"slices"
	"testing"
)

func TestGetFilesystemTypeFromBuffer_EXT2(t *testing.T) {
	buffer := make([]byte, 64*1024)

	// Insert the ext series magic number at the appropriate offset
	copy(buffer[0x438:], []byte{0x53, 0xEF}) // ext series magic number: 0xEF53

	fsType, err := GetFilesystemTypeFromBuffer(buffer)
	if err != nil {
		t.Fatalf("Failed to get filesystem type from buffer: %s", err)
	}
	if fsType != "ext2" {
		t.Fatalf("Expected ext2, but got %s", fsType)
	}
}

func TestGetFilesystemTypeFromBuffer_WithFilter(t *testing.T) {
	buffer := make([]byte, 64*1024)

	// Insert the ext series magic number at the appropriate offset
	copy(buffer[0x438:], []byte{0x53, 0xEF}) // ext series magic number: 0xEF53

	fsType, err := GetFilesystemTypeFromBuffer(buffer, WithFilesystemFilter("xfs"))
	if err == nil {
		t.Fatalf("Expected error because reported filesystem %s is not in filesystem filter", fsType)
	}
}

func TestSupportedFilesystems(t *testing.T) {
	supportedFS := SupportedFilesystems()
	if len(supportedFS) == 0 {
		t.Fatalf("No supported filesystems returned.")
	}

	expectedFS := []string{"xfs", "ext4", "ext3", "ext2"}
	for _, efs := range expectedFS {
		if !slices.Contains(supportedFS, efs) {
			t.Errorf("Filesystem %s is not supported", efs)
		}
	}

	t.Logf("Supported filesystems:")
	for _, fs := range supportedFS {
		t.Logf("- %s", fs)
	}
}

func TestGetBlockSizeFromBuffer(t *testing.T) {
    buffer := make([]byte, 64*1024)

    // Insert the ext series magic number
    copy(buffer[0x438:], []byte{0x53, 0xEF}) // ext series magic number: 0xEF53

    // Insert s_log_block_size to represent block size of 4096
    // superblock for ext series starts at 0x400 (1024) and s_log_block_size is located at offset 0x18 (24)
    // set flag to 2 because 2^2 * 1024 = 4096
    buffer[0x400+0x18] = 2

    blocksize, err := GetBlockSizeFromBuffer(buffer)
    if err != nil {
        t.Fatalf("Failed to get block size from buffer: %s", err)
    }
    if blocksize != 4096 {
        t.Fatalf("Expected block size of 4096, but got %d", blocksize)
    }

    t.Logf("Determined block size: %d", blocksize)
}

