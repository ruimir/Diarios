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

	var numTransferencia int
	errQuery = s.pceDB.QueryRowContext(ctx, "select NUM_TRANSFERENCIA from pceinternados a1 where INT_EPISODIO = :int_episodio and COD_ESPECIALIDADE_PREV = :cod_especialidade and DTA_SAIDA is null", numEpisodio, codEspecialidade).Scan(&numTransferencia)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("paciente nao esta internado")
		} else {
			return errors.New("erro a identificar se paciente esta internado")
		}
	}

	if tipoDiario == "CS" {
		return processCSDiario(ctx, id, errQuery, s, numProcesso, menuPCE, operacao, createDiary, menuPCEXML, nome, dataRegisto, numMecanografico, numEpisodio, codModulo, nomeMedico, designacaoEspecialidade, confidencial, problema, diario, diarioBsimple, err, numSequencial)
	} else {
		return processASCouDSDiario(ctx, s, tipoDiario, codEspecialidade, numEpisodio, numProcesso, nomeMedico, numOrdem, diario, designacaoEspecialidade)
	}

}

func processASCouDSDiario(ctx context.Context, s service, tipoDiario string, codEspecialidade string, numEpisodio int, numProcesso string, nomeMedico string, numOrdem string, diario string, designacaoEspecialidade string) error {

	var dataCriacao string
	errQuery := s.pceDB.QueryRowContext(ctx, "select  to_char(nvl((select (select max(to_date(to_char(a2.dta_saida, 'ddmmyyyy') || to_char(to_date(a2.hora_saida, 'sssss'), 'hh24mi'), 'ddmmyyyyhh24mi')) from PCEINTERNADOS a2 where a1.INT_EPISODIO = a2.INT_EPISODIO and a2.NUM_TRANSFERENCIA < a1.NUM_TRANSFERENCIA) dta_admissao from pceinternados a1 where INT_EPISODIO = :int_episodio and COD_ESPECIALIDADE_PREV = :cod_especialidade_prev and DTA_SAIDA is null), (select to_date(to_char(dta_internamento, 'ddmmyyyy') || to_char(to_date(hora_internamento, 'sssss'), 'hh24mi'), 'ddmmyyyyhh24mi') dta_admissao from pceadmissoes where INT_EPISODIO = :int_episodio2)),'YYYY-MM-DD HH24:MI:SS') from dual", numEpisodio, codEspecialidade, numEpisodio).Scan(&dataCriacao)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return errors.New("data de criacao nao existe")
		} else {
			return errors.New("erro a obter data de criacao")
		}
	}

	//vou fazer um select a tabla docadmissao com o numEpisodio e a data de cricacaoi =data do select anterior

	var existeDocAdmissao = true
	var nProcesso string

	errQuery = s.pceDB.QueryRowContext(ctx, "select nProcesso from DOC_ADMISSAO where episodio=:numEpisodio and DATACRIACAO=to_date(:data,'YYYY-MM-DD HH24:MI:SS')", numEpisodio, dataCriacao).Scan(&nProcesso)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			existeDocAdmissao = false
		} else {
			return errors.New("erro a fazer query a tabela docadmissao")
		}
	}

	var dadosNovaAlta = genericDadosNovaAlta()
	var documentosAIDA = genericDocumentosAIDA()

	//preencher documentosAIDA
	documentosAIDA.Processo = numProcesso
	documentosAIDA.Episodio = strconv.Itoa(numEpisodio)
	documentosAIDA.Publicado = numOrdem
	documentosAIDA.Autor = nomeMedico + " (NO)" + numOrdem
	documentosAIDA.RELATORIO.Hora = ""
	documentosAIDA.RELATORIO.Data = ""
	documentosAIDA.RELATORIO.Autor = numOrdem
	documentosAIDA.RELATORIO.ENQUADRAMENTO1.Valor = diario

	//preencher dadosNovaAlta
	dadosNovaAlta.Data = ""
	dadosNovaAlta.Hora = ""
	dadosNovaAlta.Autor = nomeMedico + " (NO)" + numOrdem
	dadosNovaAlta.Episodio = strconv.Itoa(numEpisodio)
	dadosNovaAlta.Especialidade = codEspecialidade
	dadosNovaAlta.RESUMO.Valor = diario

	//preparar XMLs
	altaXML, err := xml.Marshal(dadosNovaAlta)
	if err != nil {
		return errors.New("erro a fazer marshal dos dadosNovaAlta")
	}

	documentosXML, err := xml.Marshal(documentosAIDA)
	if err != nil {
		return errors.New("erro a fazer marshal dos documentosAIDA")
	}

	// Begin the transaction
	tx, err := s.pceDB.Begin()
	if err != nil {
		return errors.New("erro a iniciar transação")
	}

	if existeDocAdmissao {
		//se o resultado for null vamos correr os dois inserts

		idMedico, err := strconv.Atoi(numOrdem)
		if err != nil {
			return err
		}

		_, err = tx.Exec("insert into doc_admissao (nprocesso, EPISODIO, modulo, nsequencial, nome, datan, sexo, servico, esp, desp, datacriacao, dataalteracao, relatorio, filex, estado, menuxml, adita, utilcria, utilpub, utnome, utilres, utnomeres, ORDEM, extra_inf) select a2.num_processo, a1.int_episodio, 'INT', a2.num_sequencial, a2.nome, a2.dta_nascimento, a2.sexo, (select idserv from servico where SSONHO = :cod_esp and ROWNUM=1), :cod_esp2, :des_esp, (select nvl((select (select max(to_date(to_char(a2.dta_saida, 'ddmmyyyy') || to_char(to_date(a2.hora_saida, 'sssss'), 'hh24miss'), 'ddmmyyyyhh24miss')) from PCEINTERNADOS a2 where a1.INT_EPISODIO = a2.INT_EPISODIO and a2.NUM_TRANSFERENCIA < a1.NUM_TRANSFERENCIA) dta_admissao from pceinternados a1 where INT_EPISODIO = :numEpisodio2 and COD_ESPECIALIDADE_PREV = :cod_esp3 and DTA_SAIDA is null), (select to_date(to_char(dta_internamento, 'ddmmyyyy') || to_char(to_date(hora_internamento, 'sssss'), 'hh24miss'), 'ddmmyyyyhh24miss') dta_admissao from pceadmissoes where INT_EPISODIO = :numEpisodio)) from dual), sysdate, 'Relatório de Admissão (Rascunho)', 'defeiton_1', 0, null, null, :id_medico, null, :nome_medico, null, null, 1, null from pceadmissoes a1, pcedoentes a2 where a1.num_sequencial = a2.num_sequencial and int_episodio = :episodio2", codEspecialidade, codEspecialidade, designacaoEspecialidade, numEpisodio, codEspecialidade, numEpisodio, idMedico, nomeMedico, numEpisodio)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.Exec("insert into doc_alta (NPROCESSO, EPISODIO, MODULO, NSEQUENCIAL, NOME, DATAN, SEXO, SERVICO, ESP, DESP, DATACRIACAO, DATAALTERACAO, RELATORIO, FILEX, ESTADO, TIPO, MENUXML, ADITA, UTILCRIA, UTILPUB, UTNOME, UTILRES, UTNOMERES, ORDEM, EXTRA_INF, FECHOU) select a2.num_processo, a1.int_episodio, 'INT', a2.num_sequencial, a2.nome, a2.dta_nascimento, a2.sexo, (select idserv from servico where SSONHO = :cod_esp and ROWNUM=1), :cod_esp2, :des_esp, (select nvl((select (select max(to_date(to_char(a2.dta_saida, 'ddmmyyyy') || to_char(to_date(a2.hora_saida, 'sssss'), 'hh24miss'), 'ddmmyyyyhh24miss')) from PCEINTERNADOS a2 where a1.INT_EPISODIO = a2.INT_EPISODIO and a2.NUM_TRANSFERENCIA < a1.NUM_TRANSFERENCIA) dta_admissao from pceinternados a1 where INT_EPISODIO = :numEpisodio and COD_ESPECIALIDADE_PREV = :cod_esp3 and DTA_SAIDA is null), (select to_date(to_char(dta_internamento, 'ddmmyyyy') || to_char(to_date(hora_internamento, 'sssss'), 'hh24miss'), 'ddmmyyyyhh24miss') dta_admissao from pceadmissoes where INT_EPISODIO = :numEpisodio2)) from dual), sysdate, 'Relatório de Alta (Rascunho)', 'defeiton_1', 9, 0, null, null, :id_medico, null, :nome_medico, null, null, 1, null, 0 from pceadmissoes a1, pcedoentes a2 where a1.num_sequencial = a2.num_sequencial and int_episodio = :episodio2", codEspecialidade, codEspecialidade, designacaoEspecialidade, numEpisodio, codEspecialidade, numEpisodio, idMedico, nomeMedico, numEpisodio)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	switch tipoDiario {
	case "ASC":
		{
			//se for asc vamos fazer o update na tabela doc_admissao
			//no update ao doc_admissao alterar a coluna dataalteracao com o valor sysdate
			_, err = tx.Exec("update DOC_ADMISSAO set DATAALTERACAO=sysdate, MENUXML=XMLTYPE(:docadmissao) where EPISODIO= :numEpisodio and to_date(:data,'YYYY-MM-DD HH24:MI:SS')", godror.Lob{IsClob: true, Reader: bytes.NewReader(documentosXML)}, numEpisodio, dataCriacao)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	case "DS":
		{
			//se for ds vamos fazer o update na tabela doc_alta
			//no update ao doc_alta alterar a coluna dataalteracao com o valor sysdate e passar a coluna estado para 0 (zero)
			_, err = tx.Exec("update DOC_ALTA set DATAALTERACAO= sysdate, ESTADO = 0, MENUXML=XMLTYPE(:docalta) where EPISODIO = :numEpisodio and DATACRIACAO = to_date(:data,'YYYY-MM-DD HH24:MI:SS')", godror.Lob{IsClob: true, Reader: bytes.NewReader(altaXML)}, numEpisodio, dataCriacao)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	default:
		{
			_ = tx.Rollback()
			return errors.New("tipo de diario nao suportado")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("erro a fazer commit da transação")

	}

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
			Text:          "",
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
						Text:          "",
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
						Text:          "",
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
