package main

import (
	"strings"
)

var rawMatches string = `Mexico - Canada: 0 - 5
Spain - Brazil: 10 - 2
Germany - France: 2 - 2
Uruguay - Italy: 6 - 6
Argentina - Australia: 3 - 1`

var matches [][]string

func preProcessData() {
	matchesArray := strings.Split(rawMatches, "\n")
	for _, val := range matchesArray {
		val := strings.Split(val, ":")
		match, score := strings.Split(val[0], "-"), strings.Split(val[1], "-")
		matches = append(matches, []string{strings.Trim(match[0], " "), strings.Trim(match[1], " "), strings.Trim(score[0], " "), strings.Trim(score[1], " ")})
	}
}
