package iucv

import (
	"golang.org/x/sys/unix"
)

func iucvConnect(sa *unix.SockaddrIUCV) (int, error) {
	fd, err := unix.Socket(unix.AF_IUCV, unix.SOCK_SEQPACKET, 0)
	if err != nil {
		return -1, err
	}

	if err := unix.Connect(fd, sa); err != nil {
		return -1, err
	}
	return fd, nil
}
