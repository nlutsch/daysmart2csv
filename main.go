package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	//runConsoleMode()
	runWebAppMode()
}

func getLeagues(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Cache-Control", "no-cache")
	company := req.URL.Query().Get("company")
	var leagues = getAllLeagues(company)

	jsonString, err := json.Marshal(leagues)
	if err != nil {
		fmt.Println(err)
	}

	resp.Write(jsonString)
}

func getTeams(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Cache-Control", "no-cache")
	leagueId := req.URL.Query().Get("leagueId")
	company := req.URL.Query().Get("company")

	var teams = getAllTeams(leagueId, company)

	jsonString, err := json.Marshal(teams)
	if err != nil {
		fmt.Println(err)
	}

	resp.Write(jsonString)
}

func getSchedule(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Cache-Control", "no-cache")
	leagueId := req.URL.Query().Get("leagueId")
	company := req.URL.Query().Get("company")
	teamId := req.URL.Query().Get("teamId")

	var events = getScheduleForTeam(teamId, leagueId, company)

	jsonString, err := json.Marshal(events)
	if err != nil {
		fmt.Println(err)
	}

	resp.Write(jsonString)
}

func runConsoleMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter a Company: ")
	fmt.Println("blackhawks")
	fmt.Println("johnnys")

	company, _ := reader.ReadString('\n')
	company = strings.Replace(company, "\r\n", "", -1)

	leagues := getAllLeagues(company)
	fmt.Println("Select a League: ")
	for _, league := range leagues {
		fmt.Println(league.Id + ": " + league.Name)
	}

	league_id, _ := reader.ReadString('\n')
	league_id = strings.Replace(league_id, "\r\n", "", -1)

	teams := getAllTeams(league_id, company)
	fmt.Println("Select a Team: ")
	for _, team := range teams {
		fmt.Println(team.Id + ": " + team.Name)
	}

	team_id, _ := reader.ReadString('\n')
	team_id = strings.Replace(team_id, "\r\n", "", -1)

	schedule := getScheduleForTeam(team_id, league_id, company)
	for _, game := range schedule {
		fmt.Println(game.VisitorTeam + " @ " + game.HomeTeam + " at " + game.EventTime.Format(time.RFC822))
	}
}

func getCurrentExecutingPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func runWebAppMode() {
	currentPath := getCurrentExecutingPath() + "/public" // For loading files from file directory
	fs := http.FileServer(http.Dir(currentPath))
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		fs.ServeHTTP(resp, req)
	})
	http.HandleFunc("/getleagues", getLeagues)
	http.HandleFunc("/getteams", getTeams)
	http.HandleFunc("/getschedule", getSchedule)

	// Start HTTP Web Server
	http.ListenAndServe(":8080", nil)
}
