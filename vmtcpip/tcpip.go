package vmtcpip

import (
	"golang.org/x/sys/unix"
)

type tcpip struct {
	fd int
}

func (t *tcpip) Hostname() string {
	//TRGCLS
	// High-order halfword = 8
	// Low-order halfword = 0
	//DATA
	// PRMMSG
	//PRMMSG
	// Binary zeros
	//ANSLEN
	// namelen + 8

	// Offset: 0
	// Name: rc
	// Length: 4
	//  The return code from the GETHOSTNAME call. A return
	//  code of 0 indicates that the call was successful. A
	//  return code of -1 indicates that the function could not
	//  be completed and that errno contains a reason code.
	// Offset: 4
	// Name: errno
	// Length: 4
	//  When the return code is -1, this field contains a reason
	//  code.
	// Note: The rest of the reply buffer is filled only if the call was successful.
	// Offset: 8
	// Name: name
	// Length: namelen
	//  The host name, not null-terminated.
	return ""
}

func NewTCPIP(user string, subtask string) (*tcpip, error) {
	return NewTCPIPWithName(user, subtask, "")
}

func NewTCPIPWithName(user string, subtask string, name string) (*tcpip, error) {
	sa := &unix.SockaddrIUCV{UserID: user, Name: name}
	fd, err := iucvConnect(sa)
	if err != nil {
		return nil, err
	}
	// TODO: Send initilization
	// TRGCL
	//  0
	// DATA
	//   BUFFER
	// BUFLEN
	//   20
	// TYPE
	//   2WAY
	// ANSLEN
	//   8
	// PRTY
	//   NO
	// BUFFER
	//   Points to a buffer in the following format:
	// Offset: 0
	// Length: 8
	//  Constant 'IUCVAPI '. The trailing blank is required.
	// Offset: 8
	// Lenth: 2
	//  Halfword integer. Maximum number of sockets that can
	//  be established on this IUCV connection. minimum: 50,
	//  Default: 50.
	// Offset: 10
	// Name: apitype
	// Length: 2
	//  X'0002'. Provided for compatibility with prior
	//  implementations of TCP/IP. Use X'0003' instead.
	//  X'0003'. Any number of socket requests may be
	//  outstanding on this IUCV connection at the same time.
	//  For AF_INET sockets only.
	//  X'0004'. Any number of socket requests may be
	//  outstanding on this IUCV connection at the same time.
	//  For AF_INET6 sockets only.
	//  For more information, see “Overlapping Socket Requests” on page 146.
	// Offset: 12
	// Name: subtaskname
	// Length: 8
	//  Eight printable characters. The combination of your user
	//  ID and subtaskname uniquely identifies the TCP/IP client
	//  using this path. This value is displayed by the NETSTAT
	//  CLIENT command.

	// ANSBUF
	//  Points to a buffer to contain the reply from TCP/IP:
	// Offset: 0
	// Length: 4
	//  Reserved
	// Offset: 4
	// Name: maxsock
	// Length: 4
	//  The maximum socket number that your application can
	//  use on this path. The minimum socket number is always
	//  0. Your application chooses a socket number for the
	//  accept, socket, and takesocket calls.

	return &tcpip{fd: fd}, nil
}

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
