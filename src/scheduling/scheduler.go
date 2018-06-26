package scheduling

import (
	"time"
	"slack"
	"wc2018"
	"log"
)

const PollingDelayRatioAfterError = 1.

func NewScheduler(pi time.Duration, sb slack.Bot, m wc2018.Matches) Scheduler {
	return Scheduler{
		pollingInterval: pi,
		slackBot: sb,
		matches: m,
	}
}

type Scheduler struct {
	pollingInterval time.Duration
	slackBot slack.Bot
	matches wc2018.Matches
}

func (s Scheduler) Run() {
	s.slackBot.Say("Someone started me. I'll keep you posted about matches highlights.")

	previousIntervalMatches := make(map[wc2018.FifaId]wc2018.Match)
	firstPollingInterval := true
	pollingInterval := s.pollingInterval

	for {
		select {
		case <-time.After(pollingInterval):
			currentMatches, err := s.matches.GetCurrent()
			if err != nil {
				log.Printf("Something went wrong. Like Italy out of the tournament.\nError: %s", err)
				pollingInterval += time.Duration(pollingInterval.Seconds() * PollingDelayRatioAfterError) * time.Second
				continue
			}

			if firstPollingInterval {
				firstPollingInterval = false
				for _, cm := range currentMatches {
					s.slackBot.Say(cm.Summary())
				}
			} else {
				for _, cm := range currentMatches {
					previousIntervalMatch, found := previousIntervalMatches[cm.FifaId]
					if !found {
						previousIntervalMatch = wc2018.Match{}
					}

					somethingHappened, highlights := cm.WhatHappenedSince(previousIntervalMatch)
					if somethingHappened {
						for _, h := range highlights {
							s.slackBot.Say(h.ToString())
						}
					}

					previousIntervalMatches[cm.FifaId] = cm
				}
			}

			pollingInterval = s.pollingInterval
		}
	}
}
