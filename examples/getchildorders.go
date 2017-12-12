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

    param := bfapi.NewGetChildOrdersParam()
	param.Child_order_id = "JFX20171212-122518-393388F"
	orders, err := bfclient.GetChildOrders(param)
	if err != nil {
		log.Println(err)
	}
	pp.Println(orders)

    orders, err = bfclient.GetChildOrdersByChildOrderId("JFX20171212-122518-393388F")
	pp.Println(orders)
}
