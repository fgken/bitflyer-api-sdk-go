package main

import (
	"log"
	"os"

	"github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
)

func main() {
	apiKey := os.Getenv("BITFLYER_API_KEY")
	apiSecret := os.Getenv("BITFLYER_API_SECRET")

	bfclient, err := bitflyerclient.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

	side := bitflyerclient.BUY
	size := 0.001
	price := 1000000.0
	limit := 1000100.0
	stop := 999000.0
	resp, err := bfclient.SendParentOrder_IFDOCO(side, size, price, limit, stop)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
