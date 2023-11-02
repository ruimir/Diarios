package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/godror/godror"
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

const updError = "operacao é UPD, mas diario não existe"

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
	var operacao string
	var tipoDiario string

	errQuery := s.pceDB.QueryRowContext(ctx, "SELECT COD_ESPECIALIDADE,NUM_SEQUENCIAL,NUM_MECANOGRAFICO,NUM_EPISODIO,COD_MODULO,DIARIO,DATA_REGISTO,CONFIDENCIAL,OPERACAO,TIPODIARIO from CLI_MOVE_DIARIO where id=:id AND ORIGEM=:origem and TIPODIARIO in ('CS','ASC','DS')", id, origem).Scan(&codEspecialidade, &numSequencial, &numMecanografico, &numEpisodio, &codModulo, &diario, &dataRegistoStr, &confidencial, &operacao, &tipoDiario)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("diario não existe")
		} else {
			return errors.New("erro a obter diario")
		}
	}

	var problema string

	if operacao == "ADI" {
		problema = "Adenda"
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

	var numProcesso, nome string

	errQuery = s.pceDB.QueryRowContext(ctx, "select NUM_PROCESSO,nome from PCEDOENTES where NUM_SEQUENCIAL=:ns", numSequencial).Scan(&numProcesso, &nome)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("doente nao existe na tabela PCEDOENTES")
		} else {
			return errors.New("erro a obter número de processo e/ou nome do doente")
		}
	}

	if tipoDiario == "CS" {
		return processCSDiario(ctx, id, errQuery, s, numProcesso, menuPCE, operacao, createDiary, menuPCEXML, nome, dataRegisto, numMecanografico, numEpisodio, codModulo, nomeMedico, designacaoEspecialidade, confidencial, problema, diario, diarioBsimple, err, numSequencial)
	} else {
		return processASCouDSDiario()
	}

}

func processASCouDSDiario() error {
	return nil
}

func processCSDiario(ctx context.Context, id int, errQuery error, s service, numProcesso string, menuPCE string, operacao string, createDiary bool, menuPCEXML PCE, nome string, dataRegisto time.Time, numMecanografico string, numEpisodio int, codModulo string, nomeMedico string, designacaoEspecialidade string, confidencial string, problema string, diario string, diarioBsimple Diario_BSIMPLE, err error, numSequencial string) error {
	errQuery = s.pceDB.QueryRowContext(ctx, "select t.MENUPCE.extract('/').getClobVal()  from MENUPCE  t where NUMPROCESSO=:numprocesso", numProcesso).Scan(&menuPCE)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			//Diario nao existe

			if operacao == "UPD" {
				return errors.New(updError)
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

		if operacao == "UPD" {
			return errors.New(updError)
		}

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
				Problema:      problema,
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
				IdDiario: id,
			}},
		}
		menuPCEXML.LP.SOAPG.RSOAP = make([]RSOAP, 0)
		menuPCEXML.LP.SOAPG.RSOAP = append(menuPCEXML.LP.SOAPG.RSOAP, rsoap)
	} else {
		if operacao != "UPD" {
			//verificar se dia já existe
			var matchFound = false
			for i, rsoap := range menuPCEXML.LP.SOAPG.RSOAP {
				if rsoap.Data == dataRegisto.Format("02-01-2006") {
					//match found, append
					matchFound = true
					atoi, err := strconv.Atoi(rsoap.RSOAPL[len(rsoap.RSOAPL)-1].Titulo)
					if err != nil {
						atoi = 0
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
						Problema:      problema,
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
						IdDiario: id,
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
							Problema:      problema,
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
							IdDiario: id,
						}},
					}
					menuPCEXML.LP.SOAPG.RSOAP = append(menuPCEXML.LP.SOAPG.RSOAP, rsoap)
				}
			}
		} else {
			for i, rsoap := range menuPCEXML.LP.SOAPG.RSOAP {
				for i2, rsoapl := range rsoap.RSOAPL {
					if rsoapl.IdDiario == id {
						menuPCEXML.LP.SOAPG.RSOAP[i].RSOAPL[i2].RSO.Valor = diarioBsimple.Diario.Diario
						break
					}
				}
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
		//criar diario, nao pode ser UPD
		_, err := tx.Exec("INSERT INTO MENUPCE VALUES(:numProcesso, XMLTYPE(:menupce),:numSeq)", menuPCEXML.Processo, godror.Lob{IsClob: true, Reader: bytes.NewReader(output)}, numSequencial)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

	} else {
		//atualizar diario
		_, err := tx.Exec("UPDATE MENUPCE set MENUPCE=XMLTYPE(:menu) where numProcesso=:numProcesso", godror.Lob{IsClob: true, Reader: bytes.NewReader(output)}, menuPCEXML.Processo)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return errors.New("erro a atualizar diario")
		}

	}

	_, err = tx.Exec("update CLI_MOVE_DIARIO set DATA_INTEGRACAO=sysdate where ID=:id", id)
	if err != nil {
		_ = tx.Rollback()
		return err
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
