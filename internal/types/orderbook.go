package types

type OrderBook interface {
	Render(x, y int)
	HandleDepthResponse(result BinanceDepthResult)
}
