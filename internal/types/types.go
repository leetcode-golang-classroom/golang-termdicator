package types

type BinanceDepthResult struct {
	// price | volume
	Asks [][]string `json:"a"`
	Bids [][]string `json:"b"`
}
type BinanceDepthResponse struct {
	Stream string             `json:"stream"`
	Data   BinanceDepthResult `json:"data"`
}
