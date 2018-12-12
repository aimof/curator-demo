package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	var botID, botName string
	var articles []article

	for {
		select {
		case inEvents := <-rtm.IncomingEvents:
			log.Println("Got event")
			switch event := inEvents.Data.(type) {
			case *slack.ConnectedEvent:
				botID = event.Info.User.ID
				botName = event.Info.User.Name
				log.Println(botID, botName)
			case *slack.MessageEvent:
				log.Println("Got Message")
				var msg string
				if event.Type == "message" {
					msg = event.Text
				} else {
					continue
				}
				log.Println(msg)
				log.Println(len(msg))
				if len(msg) < 12 {
					continue
				} else if len(msg) == 12 {
					log.Println("list")
					for _, a := range articles {
						rtm.SendMessage(rtm.NewOutgoingMessage(a.toString(), event.Channel))
					}
					continue
				}

				if msg[2:11] == botID {
					rows := strings.Split(msg, "\n")

					var a = article{
						curatorName: event.User,
					}
					for i, row := range rows {
						if i == 0 {
							continue
						}
						if len(row) < 5 {
							a.comment = a.comment + row + "\n"
							continue
						}
						if strings.HasPrefix(row, "<http") {
							a.url = row
						} else {
							a.comment = a.comment + row + "\n"
						}
					}
					log.Println(a)
					if a.url != "" {
						articles = append(articles, a)
						rtm.SendMessage(rtm.NewOutgoingMessage("ok", event.Channel))
						continue
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("ng", event.Channel))
					}
				}
			}
		}
	}
}

type article struct {
	curatorName string
	url         string
	comment     string
}

func (a article) toString() string {
	return fmt.Sprintf("by <@%s>,\nurl: <%s>,\ncomment: %s", a.curatorName, a.url, a.comment)
}
