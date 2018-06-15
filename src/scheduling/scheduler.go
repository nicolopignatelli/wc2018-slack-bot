package scheduling

import (
	"time"
	"slack"
	"wc2018"
)

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

	previousIntervalMatch := wc2018.Match{}
	firstPollingInterval := true

	for {
		select {
		case <-time.After(s.pollingInterval):
			currentMatch, err := s.matches.GetCurrent()
			if err != nil {
				s.slackBot.Say("Something went wrong. Like Italy out of the tournament.")
				continue
			}

			if firstPollingInterval {
				firstPollingInterval = false
				s.slackBot.Say(currentMatch.Summary())
			} else {
				somethingHappened, highlights := currentMatch.WhatHappenedSince(previousIntervalMatch); if somethingHappened {
					for _, h := range highlights {
						s.slackBot.Say(h.ToString())
					}
				}
			}

			previousIntervalMatch = currentMatch
		}
	}
}
