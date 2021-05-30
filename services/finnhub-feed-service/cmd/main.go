package main

import (
	"fmt"

	"github.com/nazarov-pro/stock-exchange/services/finnhub-feed-service/pkg/conf"
	"golang.org/x/net/websocket"
)

type subscribeData struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

type response struct {
	Data []responseItem `json:"data"`
	Type string         `json:"type"`
}

type responseItem struct {
	LastPrice float64 `json:"p"`
	Symbol    string  `json:"s"`
	Timestamp int64   `json:"t"`
	Volume    float64 `json:"v"`
}

func main() {
	var (
		origin = "http://localhost/"
		token  = conf.Config.GetString("app.finnhub.token")
		wsUrl  = fmt.Sprintf(conf.Config.GetString("app.finnhub.ws.url"), token)
	)

	ws, err := websocket.Dial(wsUrl, "", origin)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	msg := &subscribeData{Type: "subscribe", Symbol: "BTC/USDT"}
	if err = websocket.JSON.Send(ws, msg); err != nil {
		panic(err)
	}

	rsp := &response{}
	if err = websocket.JSON.Receive(ws, rsp); err != nil {
		panic(err)
	}
	printData(rsp)
	
	if err = websocket.JSON.Receive(ws, rsp); err != nil {
		panic(err)
	}
	printData(rsp)
	
}

func printData(rsp *response) {
	fmt.Printf("Type: %s, Size of data: %d\n", rsp.Type, len(rsp.Data))
	for _, item := range rsp.Data {
		fmt.Printf("Last Price: %f, Symbol: %s, Timestamp: %d, Volume: %f\n", item.LastPrice, item.Symbol, item.Timestamp, item.Volume)
	}
}
