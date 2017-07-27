FROM alpine:edge AS build
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.8.3-r0 gcc=6.4.0-r4 g++=6.4.0-r4
WORKDIR /github.com/IhorBondartsov/tournamentAPI/src/
ENV GOPATH /github.com/IhorBondartsov/tournamentAPI/src/
ADD src /github.com/IhorBondartsov/tournamentAPI/src/cmd/
RUN go get cmd
    RUN CGO_ENABLED=1 GOOS=linux go install -a cmd

FROM alpine:edge
WORKDIR /github.com/IhorBondartsov/tournamentAPI/src/
RUN cd /github.com/IhorBondartsov/tournamentAPI/src/
COPY --from=build /github.com/IhorBondartsov/tournamentAPI/src/cmd /github.com/IhorBondartsov/tournamentAPI/src/cmd
CMD ["/github.com/IhorBondartsov/tournamentAPI/src/cmd"]