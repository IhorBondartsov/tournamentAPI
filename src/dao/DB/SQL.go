package DB

import (
	"github.com/IhorBondartsov/tournamentAPIv/tournamentAPI/src/models"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/IhorBondartsov/tournamentAPIv/tournamentAPI/src/dao"
	log "github.com/Sirupsen/logrus"
)

const (
	nameDB             = "tournament"
	tablePlayers       = "players"
	tableTournaments   = "tournaments"
	tableTeam          = "team"
	tableTeams         = "teams"
	tableRunTournament = "runTournament"
)

type SQL struct {
	DB       *sql.DB
	User     string
	Password string
	URI      string
	TypeConn string
}

func (base *SQL) Init() (dao.DAOInterface, error) {
	db, err := sql.Open("mysql", base.User+":"+base.Password+"@"+base.TypeConn+"("+base.URI+")/"+nameDB)
	if err != nil {
		return nil, err
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	base.DB = db
	return &SQL{DB: db}, nil
}

func (base *SQL) Close() {
	base.DB.Close()
}

func (base *SQL) CreatePlayer(id string, balance int) error {
	rows, err := base.DB.Query("INSERT INTO "+tablePlayers+" VALUES (?,?);", id, balance)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (base *SQL) CreateTournament(tournamentId, deposit int) error {
	rows, err := base.DB.Query("INSERT INTO "+tableTournaments+" VALUES (?,?);", tournamentId, deposit)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()
	return nil
}

func (base *SQL) UpdatePlayerBalance(id string, balance int) error {
	rows, err := base.DB.Query("UPDATE "+tablePlayers+" SET balance = balance + ? WHERE id_player = ?;", balance, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (base *SQL) GetPlayer(id string) *models.Player {
	var player models.Player
	rows, err := base.DB.Query("SELECT * FROM "+tablePlayers+" WHERE id_player = ?;", id)

	if err != nil {
		log.Error(err)
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&player.Id, &player.Balance)
		return &player
	}
	return nil
}
func (base *SQL) GetTurnament(id int) *models.Tournament {
	var tournament models.Tournament
	rows, err := base.DB.Query("SELECT * FROM "+tableTournaments+" WHERE id_tournament = ?;", id)

	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&tournament.Id, &tournament.Deposit)
		return &tournament
	}
	return nil
}

func (base *SQL) CheckPlayer(id string) bool {
	rows, err := base.DB.Query("SELECT balance FROM "+tablePlayers+" WHERE id_player = ?", id)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (base *SQL) CheckTournament(tournamentId int) bool {
	rows, err := base.DB.Query("SELECT balance FROM "+tableTournaments+" WHERE id_tournament = ?", tournamentId)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (base *SQL) DeleteAll() error {
	changeSafeUpdates0 := "SET SQL_SAFE_UPDATES = 0;"
	deleteTournaments := "DELETE FROM " + tableTournaments
	deletePlayers := "DELETE FROM " + tablePlayers
	deleteRunTournament := "DELETE FROM " + tableRunTournament
	deleteTeam := "DELETE FROM " + tableTeam
	deleteTeams := "DELETE FROM " + tableTeams
	changeSafeUpdates1 := "SET SQL_SAFE_UPDATES = 1;"

	commands := []string{changeSafeUpdates0, deleteTournaments, deletePlayers, deleteRunTournament, deleteTeam, deleteTeams, changeSafeUpdates1}

	for _, val := range commands {
		_, err := base.DB.Query(val)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (base *SQL) DeletedTeamsFromTournament(tournamentId int) error {

	return nil
}

func (base *SQL) GetTournamentWithPlayer(playerId string) *models.Tournament {
	idTeam, err := base.getTeamID(playerId)

	if err != nil {
		log.Error(err)
		return nil
	}
	idTournament, err := base.getTournamentID(idTeam)
	if err != nil {
		log.Error(err)
		return nil
	}

	tournament := base.getTournament(idTournament)
	if tournament == nil {
		return nil
	}

	rows, err := base.DB.Query("SELECT id_team FROM " + tableRunTournament + " WHERE id_tournament = '" + idTournament + "'")
	if err != nil {
		return nil
	}

	for rows.Next() {

		rows.Scan(&idTeam)
		team := base.getTeam(idTeam)
		tournament.TeamMembers = append(tournament.TeamMembers, *team)
	}

	return tournament
}

func (base *SQL) UpdateTeamTournament(tournamentId int, team models.Team) error {
	var idTeam int
	rows, err := base.DB.Query("SELECT MAX(id_team) AS id_team FROM " + tableTeam)
	if err != nil {
		log.Error(err)
		return err
	}

	if rows.Next() {
		rows.Scan(&idTeam)
		idTeam += 1
	}

	_, err = base.DB.Query("insert into team value(?);", idTeam)
	if err != nil {
		log.Error(err)
		return err
	}

	for key, val := range team.Proportion {
		_, err = base.DB.Query("insert into teams value(?, ?, ?);", key, idTeam, val)
		if err != nil {
			log.Error(err)
			return err
		}

	}
	_, err = base.DB.Query("insert into "+tableRunTournament+" value(?, ?);", tournamentId, idTeam)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (base *SQL) getTournamentID(idTeam string) (string, error) {
	var idTournament string

	rows, err := base.DB.Query("SELECT id_tournament FROM " + tableRunTournament + " WHERE id_team = '" + idTeam + "'")

	if err != nil {
		log.Error(err)
		return idTournament, err
	}

	if rows.Next() {
		rows.Scan(&idTournament)
	}
	return idTournament, nil
}

func (base *SQL) getTeamID(idplayer string) (string, error) {
	var idTeam string
	rows, err := base.DB.Query("SELECT id_team FROM " + tableTeams + " WHERE id_player ='" + idplayer + "'")
	if err != nil {
		log.Error(err)
		return idTeam, err
	}

	if rows.Next() {
		rows.Scan(&idTeam)
	}

	return idTeam, nil
}

func (base *SQL) getTournament(idTournament string) (*models.Tournament) {
	tournament := models.Tournament{}
	rows, err := base.DB.Query("SELECT * FROM " + tableTournaments + " WHERE id_tournament = '" + idTournament + "'")
	if err != nil {
		return nil
	}

	if rows.Next() {
		rows.Scan(&tournament.Id, &tournament.Deposit)
	}
	return &tournament
}

func (base *SQL) getTeam(idTeam string) *models.Team {
	team := models.Team{Proportion: make(map[string]int)}

	var idPlayer string
	var part int
	var some int

	rows, err := base.DB.Query("SELECT * FROM " + tableTeams + " WHERE id_team = '" + idTeam + "'")

	if err != nil {
		return nil
	}

	for x := 0; rows.Next(); x++ {
		rows.Scan(&idPlayer, &part, &some)
		team.Proportion[idPlayer] = some
		if x == 0 {
			team.LeadingPlayer = &models.Player{Id: idPlayer, Balance: part}
		}
	}
	return &team
}
