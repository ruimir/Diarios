package service

import (
	"context"
	"database/sql"
	"errors"
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

	var codEspecialidade string

	errQuery := s.pceDB.QueryRowContext(ctx, "SELECT COD_ESPECIALIDADE from CLI_MOVE_DIARIO where id=:id AND ORIGEM=:origem", id, origem).Scan(&codEspecialidade)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("diario não existe")
		}
	}

	var designacaoEspecialidade string

	errQuery = s.pceDB.QueryRowContext(ctx, "select DESIGNACAO from sil.especialidades where CODIGO=:cod", codEspecialidade).Scan(&designacaoEspecialidade)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("especialidade definida no diario não existe")
		}
	}

	var nomeMedico string
	var numOrdem string

	errQuery = s.pceDB.QueryRowContext(ctx, "select UTILIZADOR,IDINT from UTILIZADORES where IDUTILIZADOR=:id", codEspecialidade).Scan(&nomeMedico, &numOrdem)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("médico definido no diario não existe")
		}
	}

	var menuPCE string
	var createDiary = false

	var menuPCEXML PCE

	errQuery = s.pceDB.QueryRowContext(ctx, "select MENUPCE from MENUPCE where PCE_NUM_SEQ=:numSeq", codEspecialidade).Scan(&menuPCE)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			//Diario nao existe
			createDiary = true
			menuPCEXML = genericPCEXML()
		}
	}

	if createDiary {
	}

	return nil
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
