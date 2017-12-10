package main

import (
	"github.com/comail/colog"
	"log"
	"os"

	bfapi "github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
)

func main() {
	/* Init logging */
	colog.Register()
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	colog.SetMinLevel(colog.LDebug)

	apiKey := os.Getenv("BITFLYER_API_KEY")
	apiSecret := os.Getenv("BITFLYER_API_SECRET")

	bfclient, err := bfapi.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

	param := bfapi.NewSendChildOrderParam()
	param.Child_order_type = bfapi.MARKET
	param.Side = bfapi.BUY
	param.Size = 0.001
	resp, err := bfclient.SendChildOrder(param)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
