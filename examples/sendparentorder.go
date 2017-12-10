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

	param := bfapi.NewSendParentOrderParam()
	param.Order_method = bfapi.SIMPLE
	parentOrder := bfapi.ParentOrder{
		Condition_type: bfapi.STOP,
		Side:           bfapi.BUY,
		Size:           0.001,
		Trigger_price:  3000000,
	}
	param.Parameters = append(param.Parameters, parentOrder)
	resp, err := bfclient.SendParentOrder(param)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

	resp, err = bfclient.SendParentOrderStop(bfapi.BUY, 3100000, 0.001)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

	var price float64 = 3000000
	resp, err = bfclient.SendParentOrderIFDOCO(bfapi.STOP, bfapi.BUY, price, price+1000, price-500, 0.001)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

	price = 1000000
	resp, err = bfclient.SendParentOrderIFDOCO(bfapi.LIMIT, bfapi.BUY, price, price+1000, price-500, 0.001)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
