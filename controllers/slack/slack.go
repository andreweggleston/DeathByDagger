package slack

import (
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/nlopes/slack"
)

func ListenAndResponse() {
	api := slack.New(config.Constants.SlackToken)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Recieved: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Println("Hello!")
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}
	}
}
