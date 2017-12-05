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

	board, err := bfclient.GetBoard()
	if err != nil {
		log.Println(err)
	}

    for i, bids := range board.Bids {
        if 10 < bids.Size {
            log.Printf("[%d] Price: %f, Size: %f\n",
                i, bids.Price, bids.Size)
        }
    }
    log.Println("-----------")
    for i, asks := range board.Asks {
        if 10 < asks.Size {
            log.Printf("[%d] Price: %f, Size: %f\n",
                i, asks.Price, asks.Size)
        }
    }
}
