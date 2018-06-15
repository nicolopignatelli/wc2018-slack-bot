package slack

import "github.com/nlopes/slack"

func NewBot(t string, c string, u string, iu string) Bot {
	params := slack.NewPostMessageParameters()
	params.Username = u
	params.IconURL = iu

	b := Bot{
		channel: c,
		client: slack.New(t),
		postMessageParams: params,
	}

	return b
}

type Bot struct {
	channel string
	client *slack.Client
	postMessageParams slack.PostMessageParameters
}

func (b Bot) Say(something string) {
	b.client.PostMessage(b.channel, something, b.postMessageParams)
}
