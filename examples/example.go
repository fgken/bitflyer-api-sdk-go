package main

import (
	"log"
	"os"

	"github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
	"github.com/k0kubun/pp"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	bfclient, err := bitflyerclient.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

	page := bitflyerclient.Page{Count: 1}
	execs, err := bfclient.GetExecutions(page)
	if err != nil {
		log.Println(err)
	}
	pp.Println(execs)
}
