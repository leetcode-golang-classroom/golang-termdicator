package websocket

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/leetcode-golang-classroom/golang-termdicator/internal/types"
)

type WebsocketHandler struct {
	ob    types.OrderBook
	conn  *websocket.Conn
	wsurl string
	sync.RWMutex
}

func NewWebsocketHandler(wsurl string, ob types.OrderBook) (*WebsocketHandler, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		return nil, err
	}
	handler := &WebsocketHandler{
		ob:    ob,
		wsurl: wsurl,
		conn:  conn,
	}
	return handler, nil
}

func (wshandler *WebsocketHandler) Start(ctx context.Context) error {
	// log.Printf("start receive data from %s\n", wshandler.wsurl)
	errCh := make(chan error, 1)
	go func() {
		for {
			var result types.BinanceDepthResponse
			if !wshandler.IsConnect() {
				wshandler.Reconncet(ctx)
			}
			err := wshandler.ReadJSON(&result)
			if err != nil {
				errCh <- err
				close(errCh)
			}
			wshandler.ob.HandleDepthResponse(result.Data)
		}
	}()
	defer wshandler.conn.Close()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}
func (wsHandler *WebsocketHandler) ReadJSON(v interface{}) error {
	wsHandler.Lock()
	defer wsHandler.Unlock()
	return wsHandler.conn.ReadJSON(&v)
}
func (wsHandler *WebsocketHandler) IsConnect() bool {
	wsHandler.Lock()
	defer wsHandler.Unlock()
	return wsHandler.conn == nil
}
func (wsHandler *WebsocketHandler) Reconncet(ctx context.Context) error {
	wsHandler.Lock()
	defer wsHandler.Unlock()
	if wsHandler.conn != nil {
		wsHandler.conn.Close()
		wsHandler.conn = nil
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsHandler.wsurl, nil)
	if err != nil {
		return err
	}
	wsHandler.conn = conn
	return nil
}
