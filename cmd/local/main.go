package main

import (
	"bronya/internal/pkg"
	"bytes"
	"log"
	"net"

	"github.com/gorilla/websocket"
	"github.com/txthinking/socks5"
)

func main() {
	// socks5.Debug = true
	s, err := socks5.NewClassicServer("127.0.0.1:8090", "127.0.0.1", "", "", 0, 60)
	if err != nil {
		// log.Fatal(err)
		return
	}

	e := s.ListenAndServe(&SocksHandler{})
	if e != nil {
		// log.Fatal(e)
	}
}

type SocksHandler struct{}

func (h *SocksHandler) TCPHandle(s *socks5.Server, conn *net.TCPConn, r *socks5.Request) (e error) {
	var p *socks5.Reply
	if r.Atyp == socks5.ATYPIPv4 || r.Atyp == socks5.ATYPDomain {
		p = socks5.NewReply(socks5.RepSuccess, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
	} else {
		p = socks5.NewReply(socks5.RepSuccess, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
	}
	p.WriteTo(conn)

	c, _, e := websocket.DefaultDialer.Dial("ws://localhost:8080/aaaa/tcp", nil)
	if e != nil {
		log.Printf("create websocket to remote failure: %v", e)
		return
	}
	defer c.Close()

	buf := bytes.Buffer{}
	buf.WriteByte(r.Atyp)
	buf.Write(r.DstAddr)
	buf.Write(r.DstPort)
	log.Print(buf.Bytes())
	c.WriteMessage(websocket.BinaryMessage, buf.Bytes())

	ws := pkg.NewTCPWebsocketConn(c)

	pkg.Relay(conn, ws)
	return
}

func (h *SocksHandler) UDPHandle(s *socks5.Server, conn *net.UDPAddr, packet *socks5.Datagram) error {
	return nil
}
