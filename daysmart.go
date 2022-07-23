package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DaysmartDate time.Time

type LeagueResponse struct {
	Leagues []LeagueInfo `json:"data"`
}

type LeagueInfo struct {
	Id         string           `json:"id"`
	Attributes LeagueAttributes `json:"attributes"`
}

type LeagueAttributes struct {
	Name string `json:"name"`
}

type LeagueTeamsResponse struct {
	Data LeagueTeamData `json:"data"`
}

type LeagueTeamData struct {
	Id            string                  `json:"id"`
	Relationships LeagueTeamRelationships `json:"relationships"`
}

type LeagueTeamRelationships struct {
	Teams LeagueTeamRelationData `json:"teams"`
}

type LeagueTeamRelationData struct {
	Data []LeagueTeamRelationshipInfo `json:"data"`
}

type LeagueTeamRelationshipInfo struct {
	Id string `json:"id"`
}

type TeamResponse struct {
	Teams []TeamInfo `json:"data"`
}

type TeamInfo struct {
	Id         string         `json:"id"`
	Attributes TeamAttributes `json:"attributes"`
}

type TeamAttributes struct {
	Name string `json:"name"`
}

type EventResponse struct {
	Events   []EventInfo `json:"data"`
	Included []EventType `json:"included`
}

type EventInfo struct {
	Id            string             `json:"id"`
	Attributes    EventAttributes    `json:"attributes"`
	Relationships EventRelationships `json:"relationships"`
}

type EventType struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes EventTypeAttributes `json:"attributes"`
}

type EventTypeAttributes struct {
	Name string `json:"name"`
}

type EventAttributes struct {
	EventStart    DaysmartDate `json:"start"`
	HomeTeamId    int          `json:"hteam_id"`
	VisitorTeamId int          `json:"vteam_id"`
}

type EventRelationships struct {
	Resource EventResource `json:"resource`
}

type EventResource struct {
	Data EventResourceData `json:"data"`
}

type EventResourceData struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type League struct {
	Id    string
	Name  string
	Teams []Team
}

type Team struct {
	Id   string
	Name string
}

type ScheduleEvent struct {
	HomeTeam    string
	VisitorTeam string
	Location    string
	EventTime   DaysmartDate
}

func (j *DaysmartDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") + "Z"
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*j = DaysmartDate(t)
	return nil
}

func (j DaysmartDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j DaysmartDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func getAllLeagues(company string) []League {
	resp, err := http.Get("https://apps.daysmartrecreation.com/dash/jsonapi/api/v1/leagues?include=facility%2Csport%2Cseason%2CprogramType%2Cproduct%2ChouseProduct%2CregistrationProduct%2CskillLevel%2CageRange%2CprereqLevel&filter[start_date__lte]=2022-07-21&filter[facility_id]=1&filter[season.end_date__gte]=2022-07-21&filter[programType.billing_type]=team&filter[programType.is_active]=true&page[size]=9000&sort=facility_id%2Cname&fields[league]=id%2Cname%2Cstart_date%2Cfacility%2Csport%2Cseason&fields[facility]=name&fields[sport]=name&fields[season]=name&company=" + company)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	var leagueResp LeagueResponse
	json.Unmarshal(body, &leagueResp)
	var leagues []League

	for _, league := range leagueResp.Leagues {
		leagues = append(leagues, League{Id: league.Id, Name: league.Attributes.Name})
	}

	return leagues
}

func getAllTeams(leagueId string, companyName string) []Team {
	resp, err := http.Get("https://apps.daysmartrecreation.com/dash/jsonapi/api/v1/leagues/" + leagueId + "?cache[save]=false&include=sport%2Cteams.homeEvents.statEvents.stat%2Cteams.visitingEvents.statEvents.stat%2CprogramType%2Cfacility%2Cseason%2Cproduct%2ChouseProduct%2CregistrationProduct%2CskillLevel%2CageRange%2CprereqLevel&company=" + companyName)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	var leagueResp LeagueTeamsResponse
	json.Unmarshal(body, &leagueResp)

	var teamStrList string

	for _, team := range leagueResp.Data.Relationships.Teams.Data {
		teamStrList = teamStrList + team.Id + "%2C"
	}

	resp, err = http.Get("https://apps.daysmartrecreation.com/dash/jsonapi/api/v1/teams?cache[save]=false&filter[id__in]=" + teamStrList + "&filter[inactive]=false&include=statEvents.stat%2Cfacility%2Cleague%2Cseason%2Cproduct%2CskillLevel%2Csport%2CageRange%2Cemployee%2CprogramType%2CcustomForm&company=" + companyName)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	var teamResp TeamResponse
	json.Unmarshal(body, &teamResp)
	var teams []Team

	for _, team := range teamResp.Teams {
		teams = append(teams, Team{Name: team.Attributes.Name, Id: team.Id})
	}

	return teams
}

func getTeamIdByName(teamName string, teams []Team) string {
	for _, team := range teams {
		if team.Name == teamName {
			return team.Id
		}
	}
	return ""
}

func getTeamNameById(teamId int, teams []Team) string {
	id := strconv.Itoa(teamId)
	for _, team := range teams {
		if team.Id == id {
			return team.Name
		}
	}
	return ""
}

func getAllEvents(leagueId string, companyName string) EventResponse {
	resp, err := http.Get("https://apps.daysmartrecreation.com/dash/jsonapi/api/v1/events?cache[save]=false&page[size]=10000&sort=start&include=resource.facility%2ChomeTeam.league%2CvisitingTeam.league%2CeventType&filter[homeTeam.league_id]=" + leagueId + "&filter[publish]=true&company=" + companyName)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	var eventResp EventResponse
	json.Unmarshal(body, &eventResp)
	return eventResp
}

func getScheduleForTeam(teamId string, leagueId string, companyName string) []ScheduleEvent {
	teams := getAllTeams(leagueId, companyName)
	eventResp := getAllEvents(leagueId, companyName)
	events := eventResp.Events
	iTeamId, _ := strconv.Atoi(teamId)
	var teamsEvents []ScheduleEvent

	for _, event := range events {
		if event.Attributes.HomeTeamId == iTeamId || event.Attributes.VisitorTeamId == iTeamId {
			teamsEvents = append(teamsEvents, ScheduleEvent{
				HomeTeam:    getTeamNameById(event.Attributes.HomeTeamId, teams),
				VisitorTeam: getTeamNameById(event.Attributes.VisitorTeamId, teams),
				Location:    getLocationFromEventResp(eventResp, event),
				EventTime:   event.Attributes.EventStart,
			})
		}
	}

	return teamsEvents
}

func getLocationFromEventResp(eventResp EventResponse, event EventInfo) string {
	facility := ""
	resource := ""

	for _, ev := range eventResp.Included {
		if ev.Type == "facility" {
			facility = ev.Attributes.Name
		} else if ev.Type == "resource" && ev.Id == event.Relationships.Resource.Data.Id {
			resource = ev.Attributes.Name
		}
	}

	return facility + ": " + resource
}

func getAllLeaguesAndTeams() []League {
	leagues := getAllLeagues("blackhawks")
	var fullLeagues []League
	for _, league := range leagues {
		fullLeagues = append(fullLeagues, League{
			Id:    league.Id,
			Name:  league.Name,
			Teams: getAllTeams(league.Id, "blackhawks")})
	}
	return fullLeagues
}
