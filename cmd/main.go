package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/VictorLowther/btree"
	"github.com/gorilla/websocket"
)

const wsendpoint = "wss://fstream.binance.com/stream?streams=btcusdt@depth"

func byBestBid(a, b *OrderBookEntry) bool {
	return a.Price >= b.Price
}

func byBestAsk(a, b *OrderBookEntry) bool {
	return a.Price < b.Price
}

type OrderBookEntry struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
type OrderBook struct {
	Asks *btree.Tree[*OrderBookEntry]
	Bids *btree.Tree[*OrderBookEntry]
}
type BinanceDepthResult struct {
	// price | volume
	Asks [][]string `json:"a"`
	Bids [][]string `json:"b"`
}
type BinanceDepthResponse struct {
	Stream string             `json:"stream"`
	Data   BinanceDepthResult `json:"data"`
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: btree.New(byBestAsk),
		Bids: btree.New(byBestBid),
	}
}
func getAskByPrice(price float64) btree.CompareAgainst[*OrderBookEntry] {
	return func(o *OrderBookEntry) int {
		switch {
		case o.Price < price:
			return -1
		case o.Price > price:
			return 1
		default:
			return 0
		}
	}
}

func getBidByPrice(price float64) btree.CompareAgainst[*OrderBookEntry] {
	return func(o *OrderBookEntry) int {
		switch {
		case o.Price > price:
			return -1
		case o.Price < price:
			return 1
		default:
			return 0
		}
	}
}
func (ob *OrderBook) handleDepthResponse(res BinanceDepthResult) {
	for _, ask := range res.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)
		if volume == 0 {
			if thing, ok := ob.Asks.Get(getAskByPrice(price)); ok {
				fmt.Printf("-- deleting level %.2f\n", price)
				ob.Asks.Delete(thing)
			}
			return
		}
		entry := &OrderBookEntry{
			Price:  price,
			Volume: volume,
		}
		ob.Asks.Insert(entry)
	}
	for _, bid := range res.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)
		if volume == 0 {
			if thing, ok := ob.Bids.Get(getBidByPrice(price)); ok {
				fmt.Printf("-- deleting level %.2f\n", price)
				ob.Bids.Delete(thing)
			}
			return
		}
		entry := &OrderBookEntry{
			Price:  price,
			Volume: volume,
		}
		ob.Bids.Insert(entry)
	}
}
func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsendpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	var (
		result BinanceDepthResponse
		ob     = NewOrderBook()
	)
	for {
		// _, msg, err := conn.ReadMessage()
		err := conn.ReadJSON(&result)
		if err != nil {
			log.Fatal(err)
		}
		ob.handleDepthResponse(result.Data)
		it := ob.Asks.Iterator(nil, nil)
		for it.Next() {
			fmt.Printf("%+v\n", it.Item())
		}
	}
}
