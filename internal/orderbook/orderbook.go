package orderbook

import (
	"fmt"
	"strconv"

	"github.com/VictorLowther/btree"
	"github.com/leetcode-golang-classroom/golang-termdicator/internal/types"
	"github.com/nsf/termbox-go"
)

// orderbook structure
type OrderBookEntry struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
type OrderBook struct {
	Asks     *btree.Tree[*OrderBookEntry]
	Bids     *btree.Tree[*OrderBookEntry]
	renderer types.Renderer
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

func (ob *OrderBook) HandleDepthResponse(res types.BinanceDepthResult) {
	for _, ask := range res.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)
		if volume == 0 {
			if thing, ok := ob.Asks.Get(getAskByPrice(price)); ok {
				// fmt.Printf("-- deleting level %.2f\n", price)
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
				// fmt.Printf("-- deleting level %.2f\n", price)
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

func byBestBid(a, b *OrderBookEntry) bool {
	return a.Price >= b.Price
}

func byBestAsk(a, b *OrderBookEntry) bool {
	return a.Price < b.Price
}
func NewOrderBook(renderer types.Renderer) *OrderBook {
	return &OrderBook{
		Asks:     btree.New(byBestAsk),
		Bids:     btree.New(byBestBid),
		renderer: renderer,
	}
}

func (ob *OrderBook) Render(x, y int) {
	it := ob.Asks.Iterator(nil, nil)
	i := 1
	ob.renderer.RenderText(x, y, "asks", termbox.ColorRed)
	for it.Next() {
		item := it.Item()
		priceStr := fmt.Sprintf("%.2f", item.Price)
		ob.renderer.RenderText(x, y+i, priceStr, termbox.ColorRed)
		i++
	}
	it = ob.Bids.Iterator(nil, nil)
	i = 1
	x = x + 10
	ob.renderer.RenderText(x, y, "bids", termbox.ColorGreen)
	for it.Next() {
		item := it.Item()
		priceStr := fmt.Sprintf("%.2f", item.Price)
		ob.renderer.RenderText(x, y+i, priceStr, termbox.ColorGreen)
		i++
	}
}
