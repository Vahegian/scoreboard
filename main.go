package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BoardEntry struct {
	Home   string
	Away   string
	ScoreH uint64
	ScoreA uint64
	// StartTimeUnix int64
	Id uint64
}

type ScoreBoard struct {
	matchId uint64
	board   map[uint64]BoardEntry
}

func (s *ScoreBoard) StartGame(homeTeam string, awayTeam string) uint64 {
	if homeTeam == "" || awayTeam == "" {
		fmt.Println("Can't Start a Match With Unknown Team")
		return 0
	}

	if homeTeam == awayTeam {
		fmt.Println("Team can't play against itself")
		return 0
	}

	for _, v := range s.board { // Using the loop to keep it simple, in production will use caching to improve performance
		if strings.EqualFold(v.Home, homeTeam) || strings.EqualFold(v.Home, awayTeam) ||
			strings.EqualFold(v.Away, homeTeam) || strings.EqualFold(v.Away, awayTeam) {
			fmt.Println("Team Is Already Playing")
			return 0
		}
	}

	s.matchId++
	s.board[s.matchId] = BoardEntry{Home: homeTeam, Away: awayTeam, ScoreH: 0, ScoreA: 0, Id: s.matchId}
	return s.matchId
}

func (s *ScoreBoard) FinishGame(id uint64) {
	delete(s.board, id)
}

func (s *ScoreBoard) UpdateScore(id uint64, homeScore uint64, awayScore uint64) {
	if _, ok := s.board[id]; ok {
		entry := s.board[id]
		entry.ScoreH = homeScore
		entry.ScoreA = awayScore
		s.board[id] = entry
	} else {
		fmt.Println("Match Not Found.")
	}
}

func (s *ScoreBoard) GetSummary() []BoardEntry {
	scoreBoardSlice := make([]BoardEntry, 0, len(s.board))
	for _, v := range s.board {
		scoreBoardSlice = append(scoreBoardSlice, v)
	}

	sort.SliceStable(scoreBoardSlice, func(i, j int) bool {
		tsi := scoreBoardSlice[i].ScoreH + scoreBoardSlice[i].ScoreA
		tsj := scoreBoardSlice[j].ScoreH + scoreBoardSlice[j].ScoreA
		if tsi == tsj {
			return scoreBoardSlice[i].Id > scoreBoardSlice[j].Id
		}
		return tsi > tsj
	})

	return scoreBoardSlice

}

func playDemo(sb *ScoreBoard) {
	for _, v := range matches {
		id := sb.StartGame(v[0], v[1])
		go func(id uint64, scores []string) {
			time.Sleep(time.Second * time.Duration(id))
			hs, err := strconv.ParseUint(scores[2], 10, 64)
			if err != nil {
				fmt.Println("Error parsing string to uint")
				return
			}
			as, err := strconv.ParseUint(scores[3], 10, 64)
			if err != nil {
				fmt.Println("Error parsing string to uint")
				return
			}
			sb.UpdateScore(id, hs, as)
		}(id, v)
		time.Sleep(time.Second)
		go func(id uint64) {
			time.Sleep(time.Second * time.Duration(15+id))
			// fmt.Printf("\nMatch with id '%d' Finished\n", id)
			sb.FinishGame(id)
		}(id)
	}
}

func displaySummary(sb *ScoreBoard) {
	summary := sb.GetSummary()
	if len(summary) == 0 {
		// fmt.Println("No matches found")
		return
	}

	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("Live Matches")
	for _, v := range summary {
		fmt.Printf("%d:\t%s - %s\t\t%d - %d\n", v.Id, v.Home, v.Away, v.ScoreH, v.ScoreA)
	}

}

func main() {
	var sb *ScoreBoard = &ScoreBoard{board: make(map[uint64]BoardEntry)}

	preProcessData()
	go playDemo(sb)
	for {
		time.Sleep(time.Second)
		displaySummary(sb)
		if len(sb.board) == 0 {
			return
		}
	}
}
