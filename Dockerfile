FROM golang

ADD .           /go/src/github.com/IhorBondartsov/tournamentAPI
ADD src         /github.com/IhorBondartsov/tournamentAPI/src
ADD src/cmd     /github.com/IhorBondartsov/tournamentAPI/src/cmd
ADD src/dao     /github.com/IhorBondartsov/tournamentAPI/src/dao
ADD src/dao/DB  /github.com/IhorBondartsov/tournamentAPI/src/dao/DB
ADD src/models  /github.com/IhorBondartsov/tournamentAPI/src/models
ADD src/server  /github.com/IhorBondartsov/tournamentAPI/src/server
ADD src/sys     /github.com/IhorBondartsov/tournamentAPI/src/sys

RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/Sirupsen/logrus

RUN go install github.com/IhorBondartsov/tournamentAPI/src/cmd/

ENTRYPOINT /go/bin/cmd


EXPOSE 8080

