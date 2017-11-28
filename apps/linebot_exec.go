package main

import (
	"log"
	"os"
	"fmt"
	"time"

	"github.com/fgken/bitflyer-api-sdk-go/bitflyerclient"
	"github.com/line/line-bot-sdk-go/linebot"
)

const POLLING_INTERVAL = 1*time.Minute

func main() {
	apiKey := os.Getenv("BITFLYER_API_KEY")
	apiSecret := os.Getenv("BITFLYER_API_SECRET")

	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	userid := os.Getenv("LINE_USERID")

	bfclient, err := bitflyerclient.New(apiKey, apiSecret)
	if err != nil {
		log.Fatal("Falied to new bitflyerclient")
	}

	page := bitflyerclient.Page{Count: 1}
	execs, err := bfclient.GetExecutions(page)
	if err != nil {
		log.Println(err)
	}

	page.After = (*execs)[0].Id

	for {
		execs, err := bfclient.GetExecutions(page)
		if err != nil {
			log.Println(err)
		}
		if 0 < len(*execs) {
			for i := len(*execs)-1; 0 <= i; i-- {
				msg := fmt.Sprintf("Date: %s\nPrice: %f (%s)\nSize: %f",
					(*execs)[i].Exec_date, (*execs)[i].Price, (*execs)[i].Side,
					(*execs)[i].Size)
				log.Println("Notify to line:\n" + msg)
				_, err = bot.PushMessage(userid, linebot.NewTextMessage(msg)).Do()
				if err != nil {
					log.Println(err)
				}
			}
			page.After = (*execs)[0].Id
		}

		time.Sleep(POLLING_INTERVAL)
	}
}
