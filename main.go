package main

import (
  "log"

  "golang.org/x/sys/unix"
)

func main() {
  sock, err := unix.Socket(unix.AF_IUCV, unix.SOCK_SEQPACKET, 0)
  if err != nil {
    log.Fatalf("IUCV socket failed: %v", err)
  }

  sa := &unix.SockaddrIUCV{UserID: "TCPIP", Name: "IUCVTEST"}
  if err := unix.Connect(sock, sa); err != nil {
    log.Fatalf("IUCV connect failed: %v", err)
  }

}
