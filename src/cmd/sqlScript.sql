drop database tournament;
create database tournament;
use tournament;

create table players (
id_player VARCHAR(40) not null,
balance int not null,
primary key (id_player));

create table tournaments(
id_tournament int not null,
deposit int not null,
primary key (id_tournament));


create table team(
id_team int not null,
primary key (id_team));

create table runTournament
(id_tournament int not null,
id_team int not null,
foreign key (id_team) references team(id_team),
foreign key (id_tournament) references tournaments(id_tournament));

create table teams
(id_player VARCHAR(40) not null,
id_team int not null,
part int not null,
foreign key (id_team) references team(id_team),
foreign key (id_player) references players(id_player));