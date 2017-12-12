package main

import (
	"log"
	"os"

	bfapi "github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
	"github.com/k0kubun/pp"
)

func main() {
	apiKey := os.Getenv("BITFLYER_API_KEY")
	apiSecret := os.Getenv("BITFLYER_API_SECRET")

	bfclient, err := bfapi.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

    orders, err := bfclient.GetChildOrdersByChildOrderId(os.Args[1])
	pp.Println(orders)
}
