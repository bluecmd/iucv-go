// This is can be removed when AF_IUCV has been merged upstream
// It is here because data(offset) is private
package iucv

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

type Cmsghdr struct {
	Len   uint64
	Level int32
	Type  int32
}

func (h *Cmsghdr) data(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(h)) + uintptr(unix.CmsgLen(int(offset))))
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}
