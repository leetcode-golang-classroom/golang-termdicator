package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leetcode-golang-classroom/golang-termdicator/internal/orderbook"
	"github.com/leetcode-golang-classroom/golang-termdicator/internal/render"
	"github.com/leetcode-golang-classroom/golang-termdicator/internal/websocket"

	"github.com/nsf/termbox-go"
)

const wsendpoint = "wss://fstream.binance.com/stream?streams=bnbusdt@depth"

func main() {
	termbox.Init()
	renderer := render.NewTermRender()
	ob := orderbook.NewOrderBook(renderer)
	wshandler, err := websocket.NewWebsocketHandler(wsendpoint, ob)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer func() {
		termbox.Close()
		cancel()
	}()
	go wshandler.Start(ctx)
	renderer.Start(ob)
}
