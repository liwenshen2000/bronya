package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strconv"
)

type Address struct {
	atyp byte
	host []byte
	port []byte
}

func ReadAddress(r io.Reader) (addr Address, err error) {
	if err = binary.Read(r, binary.BigEndian, &addr.atyp); err != nil {
		return
	}
	switch addr.atyp {
	case 0x01, 0x04:
		addr.host = make([]byte, 4*addr.atyp)
		if _, err = io.ReadFull(r, addr.host); err != nil {
			return
		}
	case 0x03:
		var buf [256]byte
		if err = binary.Read(r, binary.BigEndian, buf[:1]); err != nil {
			return
		}
		addr.host = buf[:1+buf[0]]
		if _, err = io.ReadFull(r, addr.host[1:]); err != nil {
			return
		}
	}
	addr.port = make([]byte, 2)
	if err = binary.Read(r, binary.BigEndian, &addr.port); err != nil {
		return
	}
	return
}

func (addr *Address) String() (str string) {
	switch addr.atyp {
	case 0x01, 0x04:
		host := net.IP(addr.host).String()
		port := strconv.Itoa(int(binary.BigEndian.Uint16(addr.port)))
		str = net.JoinHostPort(host, port)
	case 0x03:
		host := string(addr.host[1:])
		port := strconv.Itoa(int(binary.BigEndian.Uint16(addr.port)))
		str = net.JoinHostPort(host, port)
	}
	return
}

func (addr *Address) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	if err = b.WriteByte(addr.atyp); err != nil {
		return
	}
	if _, err = b.Write(addr.host); err != nil {
		return
	}
	if _, err = b.Write(addr.port); err != nil {
		return
	}
	return b.WriteTo(w)
}
