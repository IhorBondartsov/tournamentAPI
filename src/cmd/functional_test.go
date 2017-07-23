package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"time"
	"fmt"
	"net/http"
	"github.com/KharkivGophers/TeamplayerAPI/src/models"
)

var timeForSleep time.Duration = 1000 * time.Millisecond

func TestFound(t *testing.T) {

	var httpClient = &http.Client{}

	Convey("Send correct JSON. Should be return all ok ", t, func() {

		exepted := models.Player{Id:"P1", Balance:300}
		time.Sleep(timeForSleep)
		//res, _ :=
		httpClient.Get("http://" + host + ":" + fmt.Sprint(port) + "/fund?playerId=P1&points=300")
		//bodyBytes, _ := ioutil.ReadAll(res.Body)
		//bodyString := string(bodyBytes)
		actual := db.GetPlayer("P1")
		So(actual, ShouldEqual, exepted)
	})
}
