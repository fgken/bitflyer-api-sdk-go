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

	client, err := bfapi.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

	page := bfapi.NewPage()
	page.SetCount(2)
	execs, err := client.GetExecutions(page)
	if err != nil {
		log.Println(err)
	}
	pp.Println(execs)
}
