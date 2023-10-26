# go-memblkid

`go-memblkid` is a Go package that provides an interface to `memblkid`, a utility written in C for probing filesystem types in a block of memory. This package allows you to leverage the power and speed of C-based filesystem probing directly in your Go applications.

## Features

- List all filesystems supported for probing by `blkid`.
- Create a `blkid` probe from a buffer.
- Retrieve the filesystem type from a buffer.

## Prerequisites

Ensure you have the necessary dependencies:

- Go (version 1.x+)
- A C compiler like GCC
- `blkid` development files

## Installation

To install the `go-memblkid` package, run:

```bash
go get github.com/khayyamsaleem/go-memblkid
```

## Usage

### Listing Supported Filesystems

```go
import "github.com/khayyamsaleem/go-memblkid/fs"

func main() {
    filesystems := fs.SupportedFilesystems()
    for _, fsType := range filesystems {
        println(fsType)
    }
}
```

### Creating a blkid probe from a buffer

Given a buffer that contains some filesystem data, you can probe it to determine the filesystem type: 

```go
import (
    "github.com/khayyamsaleem/go-memblkid/fs"
    // ... other imports ...
)

func main() {
    // Assuming `data` is a byte slice containing filesystem data:
    fsType := fs.GetFilesystemTypeFromBuffer(data)
    println(fsType)
}
```

To capture the most supported filesystems, supply a buffer of at least 64KiB off of the head of the device / allocation.

## Behind the Scenes

The Go functions in this package serve as wrappers around the C code in `memblkid`.

## Contributions

Go nuts.

## License

This project is licensed under the MIT License.

## Acknowledgments

Thank you to the `libblkid` / `util-linux/util-linux` contributors! I stand on the shoulders of giants.
