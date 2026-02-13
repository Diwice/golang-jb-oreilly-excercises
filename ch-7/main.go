package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Ranker interface {
	Ranking() []string
}

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams []Team
	Wins  map[string]int
}

func (l League) teamExists(teamName string) bool {
	for i := range l.Teams {
		if l.Teams[i].Name == teamName {
			return true
		}
	}
	return false
}

func (l *League) MatchResult(teamOneName, teamTwoName string, teamOneScore, teamTwoScore int) error {
	if !(l.teamExists(teamOneName) && l.teamExists(teamTwoName)) {
		return fmt.Errorf("Either of teams do not exist")
	}

	if teamOneScore > teamTwoScore {
		l.Wins[teamOneName]++
	} else if teamTwoScore > teamOneScore {
		l.Wins[teamTwoName]++
	}

	return nil
}

func (l League) Ranking() []string {
	newSlice := make([]string, 0, len(l.Wins))
	for i := range l.Wins {
		newSlice = append(newSlice, i)
	}

	sort.SliceStable(newSlice, func(i, j int) bool {
		return l.Wins[newSlice[i]] < l.Wins[newSlice[j]]
	})

	return newSlice
}

func makeLeague() *League {
	return &League{
		Teams: []Team{},
		Wins:  map[string]int{},
	}
}

func RankPrinter(rank Ranker, w io.Writer) error {
	rankings := rank.Ranking()
	for i := range rankings {
		if _, err := io.WriteString(w, rankings[i]+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	newLeague := makeLeague()
	fmt.Println(newLeague)
	fmt.Println(newLeague.MatchResult("Uno", "Dos", 2, 1))
	newLeague.Teams = []Team{Team{Name: "Uno", Players: []string{"Somebody"}}, Team{Name: "Dos", Players: []string{"Once"}}, Team{Name: "Tres", Players: []string{"Told"}}}
	newLeague.MatchResult("Uno", "Dos", 2, 1)
	newLeague.MatchResult("Dos", "Tres", 1, 1)
	newLeague.MatchResult("Tres", "Uno", 2, 1)
	newLeague.MatchResult("Uno", "Tres", 2, 1)
	fmt.Println(newLeague)
	fmt.Println(newLeague.Ranking())
	RankPrinter(newLeague, os.Stdout)
}
