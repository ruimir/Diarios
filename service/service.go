package service

import (
	"context"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"log"
)

type Service interface {
	IntegrarDiario(ctx context.Context, id int, origem string) (err error)
}

type service struct {
	pceDB     *sqlx.DB
	jwtSecret string
}

func (s service) IntegrarDiario(ctx context.Context, id int, origem string) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewService(pceCon string, jwtSecret string) Service {

	pcedb, OpenErr := sqlx.Open("godror", pceCon)
	if OpenErr != nil {
		log.Fatal("erro a abrir base de dados pce: ", OpenErr)
	}
	//Verificar se o segredo da chave JWT é longa
	if len(jwtSecret) < 50 {
		log.Fatal("segredo da jwt muito curto")
	}

	pingErr := pcedb.Ping()
	if pingErr != nil {
		log.Fatal("erro a conectar a base de dados pce: ", pingErr)
	}

	return &service{
		pceDB:     pcedb,
		jwtSecret: jwtSecret,
	}
}
