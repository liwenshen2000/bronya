package remote

import (
	"bronya/internal/pkg"
	"bronya/internal/pkg/protocol"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TCPHandler(w http.ResponseWriter, r *http.Request) {
	// id := chi.URLParamFromCtx(r.Context(), "id")
	network := chi.URLParamFromCtx(r.Context(), "network")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmp := pkg.NewTCPWebsocketConn(ws)
	defer tmp.Close()
	addr, err := protocol.ReadAddress(tmp)

	if err != nil {
		return
	}
	conn, err := net.Dial(network, addr.String())
	if err != nil {
		return
	}
	log.Printf("dial server with: %s", conn.RemoteAddr().String())
	pkg.Relay(tmp, conn)
}

// func UDPHandler(w http.ResponseWriter, r *http.Request) {
// 	// id := chi.URLParamFromCtx(r.Context(), "id")
// 	network := chi.URLParamFromCtx(r.Context(), "network")
// 	ws, err := upgrader.Upgrade(w, r, nil)

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	conn, _ := net.ListenUDP(network, nil)
// 	defer conn.Close()

// 	go func() {
// 		for {
// 			_, pr, _ := ws.NextReader()

// 			addr, _ := protocol.ReadAddress(pr)
// 			buf, _ := io.ReadAll(pr)

// 			daddr, _ := net.ResolveUDPAddr(network, addr.String())

// 			conn.WriteToUDP(buf, daddr)
// 		}
// 	}()

// 	for {
// 		buf := make([]byte, 4096)
// 		n, raddr, _ := conn.ReadFromUDP(buf)

// 		pw, _ := ws.NextWriter(websocket.BinaryMessage)

// 		if len(raddr.IP) == 4 {
// 			pw.Write([]byte{0x01})
// 		} else {
// 			pw.Write([]byte{0x04})
// 		}
// 		pw.Write(raddr.IP)
// 		binary.Write(pw, binary.BigEndian, uint16(raddr.Port))

// 		pw.Write(buf[:n])

// 		pw.Close()
// 	}
// }
