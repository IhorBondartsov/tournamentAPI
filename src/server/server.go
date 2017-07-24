package server

import (
	"github.com/gorilla/mux"
	"fmt"
	"time"
	"github.com/gorilla/handlers"
	"net/http"
	"github.com/KharkivGophers/TeamplayerAPI/src/dao"
	"github.com/KharkivGophers/TeamplayerAPI/src/sys"
	"strconv"
	"github.com/KharkivGophers/TeamplayerAPI/src/models"
	"encoding/json"
	"reflect"
)

type HTTPServer struct {
	Port     int32
	Host     string
	StopProg chan struct{}
	DB       dao.DAOInterface
}

func NewHTTPServer(port int32, host string, stop chan struct{}) *HTTPServer {
	return &HTTPServer{
		Port:     port,
		Host:     host,
		StopProg: stop,

	}
}

func (server *HTTPServer) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("HTTPServer Failed")
			server.StopProg <- struct{}{}
		}
	}()

	go InitLogger()

	r := mux.NewRouter()
	r.HandleFunc("/fund", server.getFundHandler).Methods(http.MethodGet)
	r.HandleFunc("/take", server.getTakeHandler).Methods(http.MethodGet)
	r.HandleFunc("/announceTournament", server.getTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("/joinTournament", server.getJoinTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("/balance", server.getBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/reset", server.getRecetHandler).Methods(http.MethodGet)
	r.HandleFunc("/resultTournament", server.getResultTournamentHandler).Methods(http.MethodPost)

	port := fmt.Sprint(server.Port)

	srv := &http.Server{
		Handler:      r,
		Addr:         server.Host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//CORS provides Cross-Origin Resource Sharing middleware
	http.ListenAndServe(server.Host+":"+port, handlers.CORS()(r))

	go log.Fatal(srv.ListenAndServe())
}

func (server HTTPServer) getResultTournamentHandler(w http.ResponseWriter, r *http.Request) {
	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getBalanceHandler. Have not player", http.StatusBadRequest, nil)
		return
	}
	defer dataBase.Close()

	var winner models.ResultTurnament
	err = json.NewDecoder(r.Body).Decode(&winner)

	if err != nil {
		server.sendFailRespons(w, "getResultTournamentHandler.Bad json.", http.StatusBadRequest, err)
		return
	}

	tournament := dataBase.GetTournamentWithPlayer(winner.Winners[0].PlayerID)
	if tournament == nil {
		server.sendFailRespons(w, "getResultTournamentHandler.Cant found tournament with player", http.StatusBadRequest, nil)
		return
	}

	team := server.getTeamWhichWon(winner.Winners[0].PlayerID, tournament.TeamMembers)
	if reflect.DeepEqual(team, models.Team{}) {
		server.sendFailRespons(w, "getResultTournamentHandler.Cant found team which won", http.StatusBadRequest, nil)
		return
	}

	err = server.addBonus(winner.Winners[0].Prize, tournament.Deposit, team, &winner, dataBase)
	if err != nil {
		server.sendFailRespons(w, "getResultTournamentHandler. Cant add bonus", http.StatusBadRequest, err)
		return
	}

	winner.TournamentId = tournament.Id
	err = json.NewEncoder(w).Encode(winner)

	if err != nil {
		server.sendFailRespons(w, "getResultTournamentHandler. Cant encode Winner", http.StatusBadRequest, err)
		return
	}
	log.Infof("Result was send. Tournament id: %v. Winner team %v. Main player %v", tournament.Id, winner.Winners, tournament.TeamMembers[0].LeadingPlayer)
}

func (server HTTPServer) addBonus(prize, deposit int, team models.Team, winner *models.ResultTurnament, baseConn dao.DAOInterface) error {
	var win []models.Winner

	for key, val := range team.Proportion {
		balance := (float32(val) / float32(deposit)) * float32(prize)
		win = append(win, models.Winner{Prize: int(balance), PlayerID: key})
		baseConn.UpdatePlayerBalance(key, int(balance))
	}
	winner.Winners = win
	return nil
}

func (server HTTPServer) getTeamWhichWon(playerId string, team []models.Team) models.Team {
	for _, val := range team {
		for key, _ := range val.Proportion {
			if key == playerId {
				return val
			}
		}

	}
	return models.Team{}
}

func (server HTTPServer) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getBalanceHandler. Have not player", http.StatusBadRequest, nil)
		return
	}
	defer dataBase.Close()

	id := server.getValue(w, r, "playerId", "Invalid id")
	if id == "" {
		return
	}
	player := dataBase.GetPlayer(id)
	if player == nil {
		server.sendFailRespons(w, "getBalanceHandler. Have not player", http.StatusBadRequest, nil)
		return
	}
	err = json.NewEncoder(w).Encode(player)
	if err != nil {
		server.sendFailRespons(w, "getBalanceHandler. Can't encode json", http.StatusInternalServerError, err)
		return
	}
	log.Infof("getBalanceHandler. Player id: %v, balance = %v", player.Id, player.Balance)
}

func (server HTTPServer) getRecetHandler(w http.ResponseWriter, r *http.Request) {
	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getRecetHandler. Cant connection to DB", http.StatusInternalServerError, err)
		return
	}
	defer dataBase.Close()

	err = dataBase.DeleteAll()
	if err != nil {
		server.sendFailRespons(w, "getRecetHandler. Can't deleted db", http.StatusInternalServerError, err)
		return
	}

	log.Info("Base was deleted")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (server HTTPServer) getFundHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getFundHandler. Cant connection to DB", http.StatusInternalServerError, err)
		return
	}
	defer dataBase.Close()

	if err != nil {
		server.sendFailRespons(w, "getTournamentHandler. Can't create db client", http.StatusInternalServerError, err)
		return
	}

	id := server.getValue(w, r, "playerId", "Invalid id")
	if id == "" {
		return
	}

	balance := server.parseNumber(w, r, "points", "Invalid points")
	if balance < 0 {
		return
	}

	if dataBase.CheckPlayer(id) {
		err = dataBase.UpdatePlayerBalance(id, balance)
		log.Infof("getFundHandler. Player id = %v. The balance has been replenished", id)
	} else {
		err = dataBase.CreatePlayer(id, balance)
		log.Infof("getFundHandler. Player (id = %v) was created", id)
	}

	if err != nil {
		log.Infof("getFundHandler. Can't write to db. Error: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (server HTTPServer) getTakeHandler(w http.ResponseWriter, r *http.Request) {

	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getTakeHandler. Cant connection to DB", http.StatusInternalServerError, err)
		return
	}
	defer dataBase.Close()

	id := server.getValue(w, r, "playerId", "Invalid id")

	if id == "" {
		server.sendFailRespons(w, "getBalanceHandler. id is empty", http.StatusBadRequest, nil)
		return
	}

	balance := server.parseNumber(w, r, "points", "Invalid points")
	if balance < 0 {
		return
	}

	if dataBase.UpdatePlayerBalance(id, 0-balance) != nil {
		server.sendFailRespons(w, "getTakeHandler. Can't write to db", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))

}

func (server HTTPServer) getTournamentHandler(w http.ResponseWriter, r *http.Request) {

	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getTournamentHandler. Cant connection to DB", http.StatusInternalServerError, err)
		return
	}
	defer dataBase.Close()

	tournamentId := server.parseNumber(w, r, "tournamentId", "Invalid tournamentId")
	if tournamentId < 0 {
		return
	}

	if dataBase.CheckTournament(tournamentId) {
		server.sendFailRespons(w, "getTournamentHandler. We have this tournamentId in the base. Choose another id", http.StatusBadRequest, nil)
		return
	}

	deposit := server.parseNumber(w, r, "deposit", "Invalid deposit")
	if deposit < 0 {
		return
	}

	if dataBase.CreateTournament(tournamentId, deposit) != nil {
		log.Errorf("getTournamentHandler. Can't write to db. Cant create tournament")
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Infof("getTournamentHandler. Tournament (id = %v) was created", tournamentId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (server HTTPServer) getJoinTournamentHandler(w http.ResponseWriter, r *http.Request) {
	dataBase, err := server.DB.Init()
	if err != nil {
		server.sendFailRespons(w, "getJoinTournamentHandler. Cant connection to DB", http.StatusInternalServerError, err)
		return
	}
	defer dataBase.Close()

	tournamentId := server.parseNumber(w, r, "tournamentId", "Invalid tournamentId")
	if tournamentId < 0 {
		return
	}

	tournament := dataBase.GetTurnament(tournamentId)
	if tournament == nil {
		server.sendFailRespons(w, "getJoinTournamentHandler. We have not this tournamentId in the base. Choose another id", http.StatusBadRequest, nil)
		return
	}

	values := r.Form
	mainPlayerId := values["playerId"][0]

	team := server.createTeam(mainPlayerId, tournament.Deposit, values["backerId"], dataBase)

	if reflect.DeepEqual(models.Team{}, team) {
		server.sendFailRespons(w, "getJoinTournamentHandler. ", http.StatusBadRequest, nil)
		return
	}

	err = dataBase.UpdateTeamTournament(tournamentId, team)
	if err != nil {
		server.sendFailRespons(w, "getJoinTournamentHandler. Cant added team to the tournament", http.StatusInternalServerError, nil)
		return
	}

	log.Infof("getJoinTournamentHandler. Team was created. %v", team)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (server HTTPServer) parseNumber(w http.ResponseWriter, r *http.Request, key, errorMessage string) int {
	id, err := strconv.Atoi(r.FormValue(key))
	if !sys.ValidPositiveNumber(id) || err != nil {
		log.Errorf("parseID. Value: %v. %v", id, errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return -1
	}
	return id
}

func (server HTTPServer) getValue(w http.ResponseWriter, r *http.Request, key, errorMessage string) string {
	value := r.FormValue(key)
	if value == "" {
		log.Errorf("getValue. Value: %v. %v", value, errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return ""
	}
	return value
}

func (server HTTPServer) createTeam(mainPlayerID string, balance int, backer []string, baseConn dao.DAOInterface) models.Team {

	var team models.Team = models.Team{Proportion: make(map[string]int)}

	mainPlayer := baseConn.GetPlayer(mainPlayerID)
	team.LeadingPlayer = mainPlayer

	if team.LeadingPlayer.Balance >= balance {
		team.Proportion[mainPlayerID] = balance
		baseConn.UpdatePlayerBalance(mainPlayerID, 0-balance)
		return team
	}

	team.Proportion[mainPlayerID] = team.LeadingPlayer.Balance
	baseConn.UpdatePlayerBalance(mainPlayerID, 0-team.LeadingPlayer.Balance)
	balance -= team.LeadingPlayer.Balance

	for key, value := range backer {
		part := balance / (len(backer) - key)
		player := baseConn.GetPlayer(value)

		if player.Balance < part {
			part = player.Balance
		}
		baseConn.UpdatePlayerBalance(player.Id, 0-part)
		balance -= part
		team.Proportion[value] = part
	}

	if balance != 0 {
		server.returnBalance(team, baseConn)
		return models.Team{}
	}

	log.Infof("Team has be created. %v", team.Proportion)
	return team
}

func (server HTTPServer) returnBalance(team models.Team, baseConn dao.DAOInterface) {
	for key, val := range team.Proportion {
		baseConn.UpdatePlayerBalance(key, val)
	}
}

func (server HTTPServer) sendFailRespons(w http.ResponseWriter, errorMessage string, code int, err error) {
	log.Errorf("%v. Error: %v", errorMessage, err)
	http.Error(w, errorMessage, code)
}
