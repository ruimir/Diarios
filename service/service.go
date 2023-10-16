package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"time"
)

type Service interface {
	IntegrarDiario(ctx context.Context, id int, origem string) (err error)
}

type service struct {
	pceDB     *sqlx.DB
	jwtSecret string
}

type Diario_BSIMPLE struct {
	Type           string `json:"type"`
	Nsequencial    string `json:"nsequencial"`
	Nome           string `json:"nome"`
	Datanascimento string `json:"datanascimento"`
	Episodio       struct {
		Numepisodio      string `json:"numepisodio"`
		Modulo           string `json:"modulo"`
		CodEspecialidade string `json:"codEspecialidade"`
		DesEspecialidade string `json:"desEspecialidade"`
	} `json:"episodio"`
	Diario struct {
		Diario string `json:"diario"`
	} `json:"diario"`
}

func (s service) IntegrarDiario(ctx context.Context, id int, origem string) (err error) {

	if origem != "BSIMPLE" {
		return errors.New("origem não suportada, apenas BSIMPLE")
	}

	var codEspecialidade string
	var numSequencial string
	var numMecanografico string
	var numEpisodio int
	var codModulo string
	var diario string
	var dataRegistoStr string
	var confidencial string

	errQuery := s.pceDB.QueryRowContext(ctx, "SELECT COD_ESPECIALIDADE,NUM_SEQUENCIAL,NUM_MECANOGRAFICO,NUM_EPISODIO,COD_MODULO,DIARIO,DATA_REGISTO,CONFIDENCIAL from CLI_MOVE_DIARIO where id=:id AND ORIGEM=:origem", id, origem).Scan(&codEspecialidade, &numSequencial, &numMecanografico, &numEpisodio, &codModulo, &diario, &dataRegistoStr, &confidencial)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("diario não existe")
		} else {
			return errors.New("erro a obter diario")
		}
	}

	var diarioBsimple Diario_BSIMPLE

	err = json.Unmarshal([]byte(diario), &diarioBsimple)
	if err != nil {
		return errors.New("erro a fazer unmarshal do diario BSIMPLE")
	}

	dataRegisto, err := time.Parse("2006-01-02T15:04:05Z07:00", dataRegistoStr)
	if err != nil {
		return errors.New("erro a converter data de registo")
	}

	if confidencial == "C" {
		confidencial = "1"
	} else {
		confidencial = "0"
	}

	var designacaoEspecialidade string

	errQuery = s.pceDB.QueryRowContext(ctx, "select DESIGNACAO from sil.especialidades where CODIGO=:cod", codEspecialidade).Scan(&designacaoEspecialidade)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("especialidade definida no diario não existe")
		} else {
			return errors.New("erro a obter especialidade definida no diario")
		}
	}

	var nomeMedico string
	var numOrdem string

	errQuery = s.pceDB.QueryRowContext(ctx, "select UTILIZADOR,IDINT from UTILIZADORES where IDUTILIZADOR=:id", numMecanografico).Scan(&nomeMedico, &numOrdem)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("médico definido no diario não existe")
		} else {
			return errors.New("erro a obter médico definido no diario")
		}
	}

	var menuPCE string
	var createDiary = false

	var menuPCEXML PCE

	errQuery = s.pceDB.QueryRowContext(ctx, "select MENUPCE from MENUPCE where PCE_NUM_SEQ=:numSeq", numSequencial).Scan(&menuPCE)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			//Diario nao existe
			var numProcesso, nome string

			errQuery = s.pceDB.QueryRowContext(ctx, "select NUM_PROCESSO,nome from PCEDOENTES where NUM_SEQUENCIAL=:ns", numSequencial).Scan(&numProcesso, &nome)
			if errQuery != nil {
				if errors.Is(errQuery, sql.ErrNoRows) {
					return errors.New("doente nao existe na tabela PCEDOENTES")
				} else {
					return errors.New("erro a obter número de processo e/ou nome do doente")
				}
			}

			createDiary = true
			menuPCEXML = genericPCEXML()
			menuPCEXML.Nome = nome
			menuPCEXML.Processo = numProcesso

		} else {
			return errors.New("erro a obter menuPCE")
		}
	}

	if !createDiary {
		if err := xml.Unmarshal([]byte(menuPCE), &menuPCEXML); err != nil {
			return errors.New("erro a fazer unmarshal do menuPCE")
		}
	}

	if len(menuPCEXML.LP.SOAPG.RSOAP) == 0 {
		//criar um diario apenas
		var rsoap = RSOAP{
			Text:          "RSOAP",
			Titulo:        dataRegisto.Format("02-01-2006") + " - SOAP Diário",
			Mostra:        "1",
			Liga:          "novasoap.aspx?x=RSOAPG",
			Versao:        "2",
			Data:          dataRegisto.Format("02-01-2006"),
			Hora:          dataRegisto.Format("15:04:05"),
			Autor:         numMecanografico,
			Problema:      "",
			Episodio:      strconv.Itoa(numEpisodio) + codModulo,
			Nome:          nomeMedico,
			Especialidade: designacaoEspecialidade,
			RSOAPL: []RSOAPL{{
				Text:          "",
				Titulo:        "",
				Mostra:        "1",
				Liga:          "",
				Conf:          confidencial,
				Versao:        "1",
				Data:          dataRegisto.Format("02-01-2006"),
				Hora:          dataRegisto.Format("15:04:05"),
				Autor:         numMecanografico,
				Problema:      "",
				Episodio:      strconv.Itoa(numEpisodio) + codModulo,
				Nome:          nomeMedico,
				Especialidade: designacaoEspecialidade,
				RSO: RSO{
					Text:   "",
					Titulo: "Dados Semiológicos",
					Valor:  diario,
				},
				RA: RA{
					Text:   "",
					Titulo: "Análise/Avaliação",
					Valor:  "",
				},
				RPI: RPI{
					Text:   "",
					Titulo: "Decisão/Plano",
					Valor:  "",
				},
			}},
		}
		menuPCEXML.LP.SOAPG.RSOAP = make([]RSOAP, 0)
		menuPCEXML.LP.SOAPG.RSOAP = append(menuPCEXML.LP.SOAPG.RSOAP, rsoap)
	} else {
		//verificar se dia já existe
		var matchFound = false
		for i, rsoap := range menuPCEXML.LP.SOAPG.RSOAP {
			if rsoap.Data == dataRegisto.Format("02-01-2006") {
				//match found, append
				matchFound = true
				atoi, err := strconv.Atoi(rsoap.RSOAPL[len(rsoap.RSOAPL)-1].Titulo)
				if err != nil {
					return errors.New("erro a converter titulo do diario para inteiro")
				}
				rsoapl := RSOAPL{
					Text:          "RSOAP",
					Titulo:        strconv.Itoa(atoi+1) + " - SOAP Diário - " + dataRegisto.Format("02-01-2006"),
					Mostra:        "1",
					Liga:          "novasoap.aspx?x=RSOAPG",
					Versao:        "2",
					Data:          dataRegisto.Format("02-01-2006"),
					Hora:          dataRegisto.Format("15:04:05"),
					Autor:         numMecanografico,
					Problema:      "",
					Episodio:      strconv.Itoa(numEpisodio) + codModulo,
					Nome:          nomeMedico,
					Especialidade: designacaoEspecialidade,
					RSO: RSO{
						Text:   "",
						Titulo: "Dados Semiológicos",
						Valor:  diarioBsimple.Diario.Diario,
					},
					RA: RA{
						Text:   "",
						Titulo: "Análise/Avaliação",
						Valor:  "",
					},
					RPI: RPI{
						Text:   "",
						Titulo: "Decisão/Plano",
						Valor:  "",
					},
				}

				menuPCEXML.LP.SOAPG.RSOAP[i].RSOAPL = append(menuPCEXML.LP.SOAPG.RSOAP[i].RSOAPL, rsoapl)

				break
			}
			if !matchFound {
				var rsoap = RSOAP{
					Text:          "RSOAP",
					Titulo:        dataRegisto.Format("02-01-2006") + " - SOAP Diário",
					Mostra:        "1",
					Liga:          "novasoap.aspx?x=RSOAPG",
					Versao:        "2",
					Data:          dataRegisto.Format("02-01-2006"),
					Hora:          dataRegisto.Format("15:04:05"),
					Autor:         numMecanografico,
					Problema:      "",
					Episodio:      strconv.Itoa(numEpisodio) + codModulo,
					Nome:          nomeMedico,
					Especialidade: designacaoEspecialidade,
					RSOAPL: []RSOAPL{{
						Text:          "",
						Titulo:        "",
						Mostra:        "1",
						Liga:          "",
						Conf:          confidencial,
						Versao:        "1",
						Data:          dataRegisto.Format("02-01-2006"),
						Hora:          dataRegisto.Format("15:04:05"),
						Autor:         numMecanografico,
						Problema:      "",
						Episodio:      strconv.Itoa(numEpisodio) + codModulo,
						Nome:          nomeMedico,
						Especialidade: designacaoEspecialidade,
						RSO: RSO{
							Text:   "",
							Titulo: "Dados Semiológicos",
							Valor:  diarioBsimple.Diario.Diario,
						},
						RA: RA{
							Text:   "",
							Titulo: "Análise/Avaliação",
							Valor:  "",
						},
						RPI: RPI{
							Text:   "",
							Titulo: "Decisão/Plano",
							Valor:  "",
						},
					}},
				}
				menuPCEXML.LP.SOAPG.RSOAP = append(menuPCEXML.LP.SOAPG.RSOAP, rsoap)
			}
		}
	}

	output, err := xml.Marshal(menuPCEXML)
	if err != nil {
		return errors.New("erro a fazer marshal do menuPCE")
	}

	// Begin the transaction
	tx, err := s.pceDB.Begin()
	if err != nil {
		return errors.New("erro a iniciar transação")
	}

	if createDiary {
		//criar diario
		_, err := tx.Exec("INSERT INTO MENUPCE VALUES(:numProcesso, :menupce,:numSeq)", menuPCEXML.Processo, string(output), numSequencial)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

	} else {
		//atualizar diario
		_, err := tx.Exec("UPDATE MENUPCE set MENUPCE=:menu where numProcesso=:numProcesso", string(output), menuPCEXML.Processo)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return errors.New("erro a atualizar diario")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("erro a fazer commit da transação")

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
