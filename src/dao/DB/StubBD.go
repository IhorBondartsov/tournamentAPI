package DB

import (
	"github.com/KharkivGophers/TeamplayerAPI/src/models"
	"github.com/KharkivGophers/TeamplayerAPI/src/dao"
)

type StubDB struct {
	Players    map[string]int // key its id, value its balance
	Tournament map[int]*models.Tournament
}

func (base *StubDB) Init() (dao.DAOInterface, error) {
	return base, nil
}
func (base *StubDB) Close()       {}

func (base *StubDB) CreatePlayer(id string, balance int) error {
	if base.Players == nil {
		base.Players = make(map[string]int)
	}
	base.Players[id] = balance
	return nil
}

func (base *StubDB) UpdatePlayerBalance(id string, balance int) error {
	base.Players[id] = base.Players[id] + balance
	return nil
}

func (base *StubDB) GetPlayer(id string) *models.Player {
	if val, ok := base.Players[id]; ok {
		return &models.Player{Id: id, Balance: val}
	}
	return nil
}

func (base *StubDB) GetTurnament(id int) *models.Tournament {
	if val, ok := base.Tournament[id]; ok {
		return val
	}
	return nil
}

func (base *StubDB) CreateTournament(tournamentId, deposit int) error {
	if base.Tournament == nil {
		base.Tournament = make(map[int]*models.Tournament)
	}
	base.Tournament[tournamentId] = &models.Tournament{Id: tournamentId, Deposit: deposit}
	return nil
}

func (base *StubDB) CheckPlayer(id string) bool {
	_, ok := base.Players[id]
	return ok
}

func (base *StubDB) CheckTournament(tournamentId int) bool {
	_, ok := base.Tournament[tournamentId]
	return ok
}

func (base *StubDB) UpdateTeamTournament(tournamentId int, team models.Team) error {
	base.Tournament[tournamentId].AddTeam(team)
	return nil
}
func (base *StubDB) DeleteAll() error {
	base.Tournament = make(map[int]*models.Tournament)
	base.Players = make(map[string]int)
	return nil
}
func (base *StubDB) DeletedTeamsFromTournament(tournamentId int) error {
	base.Tournament[tournamentId].DeletedTeam()
	return nil
}
func (base *StubDB) GetTournamentWithPlayer(playerId string) *models.Tournament {
	for _, val := range base.Tournament {
		if base.haveTeam(val.TeamMembers, playerId) {
			return val
		}
	}
	return nil
}

func (base *StubDB) haveTeam(teams []models.Team, playerId string) bool {
	for _, val := range teams {
		if val.LeadingPlayer.Id == playerId {
			return true
		}
	}
	return false
}
