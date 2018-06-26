package wc2018

import (
	"time"
	"fmt"
	"github.com/google/go-cmp/cmp"
)

type FifaId string

type Match struct {
	FifaId FifaId `json:"fifa_id"`
	Status string `json:"status"`
	Datetime time.Time `json:"datetime"`
	Time string `json:"time"`
	HomeTeam Team `json:"home_team"`
	AwayTeam Team `json:"away_team"`
	Winner string `json:"winner"`
	HomeTeamEvents Events `json:"home_team_events"`
	AwayTeamEvents Events `json:"away_team_events"`
}

func (m Match) IsInProgress(startEndThreshold time.Duration) bool {
	if m.Status == "in progress" {
		return true
	}

	dateTime := m.Datetime.In(time.Now().Location())

	if m.Status == "future" && dateTime.Sub(time.Now()) < startEndThreshold {
		return true
	}

	// 90 minutes + interval. Won't work perfectly with extra times.
	matchDuration, _ := time.ParseDuration("105m")
	if m.Status == "completed" && time.Now().Sub(dateTime.Add(matchDuration)) < startEndThreshold {
		return true
	}

	return false
}

func (m Match) Summary() string {
	if !m.IsInProgress(time.Second) {
		return "No one is playing right now."
	}

	return fmt.Sprintf("%s is playing against %s.\n%s %d - %d %s\nTime: %s",
		m.HomeTeam.Country, m.AwayTeam.Country, m.HomeTeam.Code, m.HomeTeam.Goals, m.AwayTeam.Goals, m.AwayTeam.Code, m.Time)
}

func (m Match) WhatHappenedSince(lastMatchData Match) (bool, Highlights) {
	var highlights Highlights

	// Match has started
	if cmp.Equal(lastMatchData, NoCurrentMatch) && !cmp.Equal(m, NoCurrentMatch) {
		highlights = append(highlights, MatchHasStarted{m})
	}

	for _, e := range m.HomeTeamEvents {
		if !lastMatchData.HomeTeamEvents.Contains(e) {
			highlights = append(highlights, eventToHighlight(e, m))
		}
	}
	for _, e := range m.AwayTeamEvents {
		if !lastMatchData.AwayTeamEvents.Contains(e) {
			highlights = append(highlights, eventToHighlight(e, m))
		}
	}

	// Half time events
	if lastMatchData.Time != "half-time" && m.Time == "half-time" {
		highlights = append(highlights, FirstHalfHasEnded{m})
	}
	if lastMatchData.Time == "half-time" && m.Time != "half-time" {
		highlights = append(highlights, SecondHalfHasStarted{m})
	}

	// Match has ended
	if (!cmp.Equal(lastMatchData, NoCurrentMatch) && lastMatchData.Status != "completed") &&
	   (m.Status == "completed" || cmp.Equal(m, NoCurrentMatch)) {
		highlights = append(highlights, MatchHasEnded{lastMatchData})
	}

	return len(highlights) > 0, highlights
}

type Team struct {
	Country string `json:"country"`
	Code string `json:"code"`
	Goals int `json:"goals"`
}

type Event struct {
	Id int `json:"id"`
	TypeOfEvent string `json:"type_of_event"`
	Player string `json:"player"`
	Time string `json:"time"`
}

type Events []Event

func (evs Events) Contains(e Event) bool {
	for _, ev := range evs {
		if e.Id == ev.Id {
			return true
		}
	}

	return false
}

type Highlight interface {
	ToString() string
}

type Highlights []Highlight

type MatchHasStarted struct {
	match Match
}

func (h MatchHasStarted) ToString() string {
	return fmt.Sprintf("%s - %s has started!", h.match.HomeTeam.Country, h.match.AwayTeam.Country)
}

type MatchHasEnded struct {
	match Match
}

func (h MatchHasEnded) ToString() string {
	return fmt.Sprintf("%s - %s has ended!\n%s %d - %d %s",
		h.match.HomeTeam.Country, h.match.AwayTeam.Country, h.match.HomeTeam.Code, h.match.HomeTeam.Goals, h.match.AwayTeam.Goals, h.match.AwayTeam.Code)
}

type FirstHalfHasEnded struct {
	match Match
}

func (h FirstHalfHasEnded) ToString() string {
	return fmt.Sprintf("First half of %s - %s has ended!\n%s %d - %d %s",
		h.match.HomeTeam.Country, h.match.AwayTeam.Country, h.match.HomeTeam.Code, h.match.HomeTeam.Goals, h.match.AwayTeam.Goals, h.match.AwayTeam.Code)
}

type SecondHalfHasStarted struct {
	match Match
}

func (h SecondHalfHasStarted) ToString() string {
	return fmt.Sprintf("Second half of %s - %s has started!", h.match.HomeTeam.Country, h.match.AwayTeam.Country)
}

type GoalWasScored struct {
	match Match
	player string
	time string
}

func (h GoalWasScored) ToString() string {
	return fmt.Sprintf("GOOOOAL! (%s) %s scored.\n%s %d - %d %s",
		h.time, h.player, h.match.HomeTeam.Code, h.match.HomeTeam.Goals, h.match.AwayTeam.Goals, h.match.AwayTeam.Code)
}

type OwnGoalWasScored struct {
	match Match
	player string
	time string
}

func (h OwnGoalWasScored) ToString() string {
	return fmt.Sprintf("Lol! (%s) %s scored an own goal...\n%s %d - %d %s",
		h.time, h.player, h.match.HomeTeam.Code, h.match.HomeTeam.Goals, h.match.AwayTeam.Goals, h.match.AwayTeam.Code)
}

type YellowCardWasIssued struct {
	player string
	time string
}

func (h YellowCardWasIssued) ToString() string {
	return fmt.Sprintf("Uh oh! Yellow card for %s (%s)", h.player, h.time)
}

type RedCardWasIssued struct {
	player string
	time string
}

func (h RedCardWasIssued) ToString() string {
	return fmt.Sprintf("Oh no! Red card for %s (%s). He's out.", h.player, h.time)
}

type PlayerEnteredAsSubstitution struct {
	player string
	time string
}

func (h PlayerEnteredAsSubstitution) ToString() string {
	return fmt.Sprintf("It's the turn of %s (%s).", h.player, h.time)
}

type PlayerWasSubstituted struct {
	player string
	time string
}

func (h PlayerWasSubstituted) ToString() string {
	return fmt.Sprintf("%s was substituted (%s).", h.player, h.time)
}

type UnrecognisedEvent struct {
	event Event
}

func (h UnrecognisedEvent) ToString() string {
	return fmt.Sprintf("Something happened but I didn't get what, exactly.\n%v", h.event)
}

func eventToHighlight(e Event, m Match) Highlight {
	switch e.TypeOfEvent {
	case "goal", "goal-penalty":
		return GoalWasScored{match: m, player: e.Player, time: e.Time}
	case "goal-own":
		return OwnGoalWasScored{match: m, player: e.Player, time: e.Time}
	case "yellow-card":
		return YellowCardWasIssued{player: e.Player, time: e.Time}
	case "red-card":
		return RedCardWasIssued{player: e.Player, time: e.Time}
	case "substitution-in":
		return PlayerEnteredAsSubstitution{player: e.Player, time: e.Time}
	case "substitution-out":
		return PlayerWasSubstituted{player: e.Player, time: e.Time}
	}

	return UnrecognisedEvent{e}
}
