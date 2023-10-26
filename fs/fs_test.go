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
