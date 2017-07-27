package dao

import (
	"github.com/IhorBondartsov/tournamentAPI/src/models"
)

type DAOInterface interface {
	Init() (DAOInterface, error)
	Close()

	CreatePlayer(id string, balance int) error
	CreateTournament(tournamentId, deposit int) error

	UpdatePlayerBalance(id string, balance int) error
	UpdateTeamTournament(tournamentId int, team models.Team) error

	GetPlayer(id string) *models.Player
	GetTurnament(id int) *models.Tournament
	GetTournamentWithPlayer(playerId string) *models.Tournament

	CheckPlayer(id string) bool
	CheckTournament(tournamentId int) bool

	DeleteAll() error
	DeletedTeamsFromTournament(tournamentId int) error
}
