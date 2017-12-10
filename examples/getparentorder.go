package main

import (
	"github.com/comail/colog"
	"log"
	"os"

	bfapi "github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
	"github.com/k0kubun/pp"
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

	param := bfapi.NewGetParentOrderParam()
	param.Parent_order_acceptance_id = "JRF20171210-093537-045359"
	order, err := bfclient.GetParentOrder(param)
	if err != nil {
		log.Println(err)
	}
	pp.Println(order)
}
