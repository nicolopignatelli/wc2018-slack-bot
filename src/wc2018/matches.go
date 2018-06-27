package wc2018

import (
	"net/http"
	"encoding/json"
		"io/ioutil"
	"time"
)

var NoCurrentMatches []Match

func NewMatches(c time.Duration) Matches {
	return Matches{
		currentMatchThreshold: c,
	}
}

type Matches struct {
	currentMatchThreshold time.Duration
}

func (ms Matches) GetCurrent() ([]Match, error) {
	response, err := http.Get("http://worldcup.sfg.io/matches/today")
	if err != nil {
		return []Match{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Match{}, err
	}

	var matches []Match

	err = json.Unmarshal(body, &matches)
	if err != nil {
		return []Match{}, err
	}

	var currentMatches = NoCurrentMatches

	for _, m := range matches {
		if m.IsInProgress(ms.currentMatchThreshold) {
			currentMatches = append(currentMatches, m)
		}
	}

	return currentMatches, nil
}
