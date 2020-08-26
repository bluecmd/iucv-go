package iucv

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	SCM_SND2WAY = 2
)

var (
	ErrNotRunningInVM = errors.New("system is not running under VM")
)

type Socket int

func Connect(sa *unix.SockaddrIUCV) (Socket, error) {
	fd, err := unix.Socket(unix.AF_IUCV, unix.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}
	if err := unix.Connect(fd, sa); err != nil {
		if err == unix.EBADFD {
			// This is triggered by the fact that when AF_IUCV is initialized it
			// chooses whether if to use classic IUCV or HiperTransport.
			// HiperTransport does not allow calling connect straight on the socket
			// so EBADFD means that we are not running under VM.
			_ = unix.Close(fd)
			return -1, ErrNotRunningInVM
		}
		_ = unix.Close(fd)
		return -1, err
	}
	return Socket(fd), nil
}

func (s *Socket) Close() error {
	return unix.Close(int(*s))
}

func Send2WayAnswer(answer []byte) []byte {
	datalen := 16 // answer pointer and size_t length
	b := make([]byte, unix.CmsgSpace(datalen))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = unix.SOL_IUCV
	h.Type = SCM_SND2WAY
	h.SetLen(unix.CmsgLen(datalen))
	*(*uintptr)(h.data(0)) = uintptr(unsafe.Pointer(&answer[0]))
	*(*uint64)(h.data(8)) = uint64(len(answer))
	return b
}
