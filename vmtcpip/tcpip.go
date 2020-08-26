package vmtcpip

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/bluecmd/iucv-go/iucv"
	"golang.org/x/sys/unix"
)

type tcpip struct {
	fd iucv.Socket
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
	// TODO: fd.Close() on error
	fd, err := iucv.Connect(sa)
	if err != nil {
		return nil, err
	}

	type tcpipInit struct {
		Magic       [8]byte
		MaxSockets  uint16
		APIType     uint16
		SubtaskName [8]byte
	}

	i := tcpipInit{
		// EBCDIC 'IUCVAPI '
		Magic: [8]byte{0xC9, 0xE4, 0xC3, 0xE5, 0xC1, 0xD7, 0xC9, 0x40},
		MaxSockets: 50,
		APIType: 4, // 3 for AF_INET, 4 for AF_INET6
		// TODO
		SubtaskName: [8]byte{0xC9, 0xE4, 0xC3, 0xE5, 0xC1, 0xD7, 0xC9, 0x40},
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, i); err != nil {
		return nil, err
	}

	answer := make([]byte, 8)
	scm := iucv.Send2WayAnswer(answer)
	if err := unix.Sendmsg(int(fd), buf.Bytes(), scm, sa, 0); err != nil {
		return nil, err
	}

	log.Printf("DEBUG: %+v", answer)

	return &tcpip{fd: fd}, nil
}
