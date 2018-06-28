package wc2018

import (
	"net/http"
	"encoding/json"
		"io/ioutil"
	"time"
)

var NoCurrentMatches map[FifaId]Match

func NewMatches(c time.Duration) Matches {
	return Matches{
		currentMatchThreshold: c,
	}
}

type Matches struct {
	currentMatchThreshold time.Duration
}

func (ms Matches) GetCurrent() (map[FifaId]Match, error) {
	response, err := http.Get("http://worldcup.sfg.io/matches/today")
	if err != nil {
		return map[FifaId]Match{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return map[FifaId]Match{}, err
	}

	var matches map[FifaId]Match

	err = json.Unmarshal(body, &matches)
	if err != nil {
		return map[FifaId]Match{}, err
	}

	var currentMatches = NoCurrentMatches

	for _, m := range matches {
		if m.IsInProgress(ms.currentMatchThreshold) {
			currentMatches[m.FifaId] = m
		}
	}

	return currentMatches, nil
}
