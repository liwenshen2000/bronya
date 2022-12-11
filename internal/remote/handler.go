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

	tmp := pkg.TCPWebsocketConn{Conn: ws}
	addr, err := protocol.ReadAddress(&tmp)

	if err != nil {
		return
	}
	conn, err := net.Dial(network, addr.String())
	if err != nil {
		return
	}
	log.Printf("dial server with: %s", conn.RemoteAddr().String())
	pkg.Relay(&tmp, conn)
}
