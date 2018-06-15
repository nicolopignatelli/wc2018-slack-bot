package main

import (
	"scheduling"
	"slack"
	"wc2018"
	"fmt"
	"github.com/caarlos0/env"
	"os"
	"time"
)

type Config struct {
	PollingInterval time.Duration `env:"WC2018_POLLING_INTERVAL" envDefault:"10s"`
	CurrentMatchThreshold time.Duration `env:"WC2018_CURRENT_MATCH_THRESHOLD" envDefault:"20s"`
	SlackToken string `env:"WC2018_SLACK_TOKEN,required"`
	SlackChannel string `env:"WC2018_SLACK_CHANNEL,required"`
	SlackBotUsername string `env:"WC2018_SLACK_BOT_USERNAME" envDefault:"FIFA World Cup 2018"`
	SlackBotIconUrl string `env:"WC2018_SLACK_BOT_ICON_URL" envDefault:"https://image.ibb.co/e40U7y/avatar_bd44be5b227e_128.jpg"`
}

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)
	exitOnError(err)

	slackBot := slack.NewBot(cfg.SlackToken, cfg.SlackChannel, cfg.SlackBotUsername, cfg.SlackBotIconUrl)

	matches := wc2018.NewMatches(cfg.CurrentMatchThreshold)

	scheduling.
		NewScheduler(cfg.PollingInterval, slackBot, matches).
		Run()
}

func exitOnError(e error) {
	if e != nil {
		fmt.Printf("%+v\n", e)
		os.Exit(1)
	}
}
