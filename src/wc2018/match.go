package wc2018

import (
	"time"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"strconv"
)

var NoMatchData = Match{}

type FifaId string

type Match struct {
	FifaId FifaId `json:"fifa_id"`
	Status string `json:"status"`
	Datetime time.Time `json:"datetime"`
	Time string `json:"time"`
	HomeTeam Team `json:"home_team"`
	AwayTeam Team `json:"away_team"`
	HomeTeamEvents Events `json:"home_team_events"`
	AwayTeamEvents Events `json:"away_team_events"`
	HomeTeamStatistics Statistics `json:"home_team_statistics"`
	AwayTeamStatistics Statistics `json:"away_team_statistics"`
	Weather Weather `json:"weather"`
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
	return fmt.Sprintf("%s is playing against %s.\n%s %d - %d %s\nTime: %s",
		m.HomeTeam.Country, m.AwayTeam.Country, m.HomeTeam.Code, m.HomeTeam.Goals, m.AwayTeam.Goals, m.AwayTeam.Code, m.Time)
}

func (m Match) WhatHappenedSince(lastMatchData Match) (bool, Highlights) {
	var highlights Highlights

	// Match has started
	if cmp.Equal(lastMatchData, NoMatchData) && !cmp.Equal(m, NoMatchData) {
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

	// Score changed
	if lastMatchData.HomeTeam.Goals != m.HomeTeam.Goals ||
		lastMatchData.AwayTeam.Goals != m.AwayTeam.Goals {
		highlights = append(highlights, ScoreHasChanged{m})
	}

	// Half time events
	if lastMatchData.Time != "half-time" && m.Time == "half-time" {
		highlights = append(highlights, FirstHalfHasEnded{m})
	}
	if lastMatchData.Time == "half-time" && m.Time != "half-time" {
		highlights = append(highlights, SecondHalfHasStarted{m})
	}

	// Match has ended
	if (!cmp.Equal(lastMatchData, NoMatchData) && lastMatchData.Status != "completed") &&
	   (m.Status == "completed" || cmp.Equal(m, NoMatchData)) {
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

type Statistics struct {
	StartingEleven Players `json:"starting_eleven"`
}

type Player struct {
	Name string `json:"name"`
	ShirtNumber int `json:"shirt_number"`
	Captain bool `json:"captain"`
}

func (p Player) ToString() string {
	s := strconv.Itoa(p.ShirtNumber) + " " + p.Name
	if p.Captain {
		s += " (C)"
	}

	return s
}

type Players []Player

func (ps Players) ToString() string {
	var psStrings []string
	for _, p := range ps {
		psStrings = append(psStrings, p.ToString())
	}

	return strings.Join(psStrings, ",")
}

type Weather struct {
	Humidity string `json:"humidity"`
	TempCelsius string `json:"temp_celsius"`
	TempFarenheit string `json:"temp_farenheit"`
	WindSpeed string `json:"wind_speed"`
	Description string `json:"description"`
}

type Highlight interface {
	ToString() string
}

type Highlights []Highlight

type MatchHasStarted struct {
	match Match
}

func (h MatchHasStarted) ToString() string {
	return fmt.Sprintf(
		"%s - %s has started!\n" +
		"| Weather | %s, %s°C/%s°F, Wind %s, Humidity %s\n" +
		"| Starting eleven %s | %s\n" +
		"| Starting eleven %s | %s",
		h.match.HomeTeam.Country, h.match.AwayTeam.Country,
		h.match.Weather.Description, h.match.Weather.TempCelsius, h.match.Weather.TempFarenheit, h.match.Weather.WindSpeed, h.match.Weather.Humidity,
		h.match.HomeTeam.Country, h.match.HomeTeamStatistics.StartingEleven.ToString(),
		h.match.AwayTeam.Country, h.match.AwayTeamStatistics.StartingEleven.ToString(),
	)
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

type ScoreHasChanged struct {
	match Match
}

func (h ScoreHasChanged) ToString() string {
	return fmt.Sprintf("%s %d - %d %s",
		h.match.HomeTeam.Code, h.match.HomeTeam.Goals, h.match.AwayTeam.Goals, h.match.AwayTeam.Code)
}

type GoalWasScored struct {
	player string
	time string
}

func (h GoalWasScored) ToString() string {
	return fmt.Sprintf("⚽ GOOOOAL! (%s) %s scored. ⚽", h.time, h.player)
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
	match Match
}

func (h YellowCardWasIssued) ToString() string {
	return fmt.Sprintf("| %s - %s | Uh oh! Yellow card for %s (%s)", h.match.HomeTeam.Code, h.match.AwayTeam.Code, h.player, h.time)
}

type RedCardWasIssued struct {
	player string
	time string
	match Match
}

func (h RedCardWasIssued) ToString() string {
	return fmt.Sprintf("| %s - %s | Oh no! Red card for %s (%s). He's out.", h.match.HomeTeam.Code, h.match.AwayTeam.Code, h.player, h.time)
}

type PlayerEnteredAsSubstitution struct {
	player string
	time string
	match Match
}

func (h PlayerEnteredAsSubstitution) ToString() string {
	return fmt.Sprintf("| %s - %s | It's the turn of %s (%s).", h.match.HomeTeam.Code, h.match.AwayTeam.Code, h.player, h.time)
}

type PlayerWasSubstituted struct {
	player string
	time string
	match Match
}

func (h PlayerWasSubstituted) ToString() string {
	return fmt.Sprintf("| %s - %s | %s was substituted (%s).", h.match.HomeTeam.Code, h.match.AwayTeam.Code, h.player, h.time)
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
		return GoalWasScored{player: e.Player, time: e.Time}
	case "goal-own":
		return OwnGoalWasScored{match: m, player: e.Player, time: e.Time}
	case "yellow-card":
		return YellowCardWasIssued{player: e.Player, time: e.Time, match: m}
	case "red-card", "yellow-card-second":
		return RedCardWasIssued{player: e.Player, time: e.Time, match: m}
	case "substitution-in":
		return PlayerEnteredAsSubstitution{player: e.Player, time: e.Time, match: m}
	case "substitution-out":
		return PlayerWasSubstituted{player: e.Player, time: e.Time, match: m}
	}

	return UnrecognisedEvent{e}
}
