package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"math/rand"
	"os"
	"time"
)

func choise() string {
	s := [3]string{"二郎", "用心棒", "うち田"}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

func main() {
	api := slack.New(os.Getenv("SLACK_API_TOKEN"))
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Recevid: ")
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				fmt.Printf(ev.Text)
				if ev.Text == "お腹すいた" {
					_, _, err := api.PostMessage(
						"#dev_chatbot",
						choise(),
						slack.PostMessageParameters{Username: "gohan", IconEmoji: ":rice:"},
					)
					if err != nil {
						fmt.Printf("%s\n", err)
						return
					}
				}
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			default:
			}
		}
	}
}
