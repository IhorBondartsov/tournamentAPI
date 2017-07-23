package models

type Player struct {
	Id      string `json:"playerId"`
	Balance int `json:"balance"`
}

type Winner struct {
	PlayerID string`json:"playerID"`
	Prize    int`json:"prize"`
}

type ResultTurnament struct {
	TournamentId int `json:"tournamentId"`
	Winners []Winner `json:"winners"`
}

type Team struct {
	LeadingPlayer *Player
	Proportion    map[string]int //key id value points
}

type Tournament struct {
	Id          int
	TeamMembers []Team
	Deposit     int
}

func (tourn *Tournament) AddTeam(team Team) {
	tourn.TeamMembers = append(tourn.TeamMembers, team)
}
func (tourn *Tournament) DeletedTeam() {
	tourn.TeamMembers = []Team{}
}
