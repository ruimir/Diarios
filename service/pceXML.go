package service

import (
	"encoding/xml"
	"log"
)

type PCE struct {
	XMLName  xml.Name `xml:"PCE"`
	Text     string   `xml:",chardata"`
	Titulo   string   `xml:"titulo,attr"`
	Nome     string   `xml:"nome,attr"`
	Processo string   `xml:"processo,attr"`
	Mostra   string   `xml:"mostra,attr"`
	BDI      struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		NOVO   struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Liga   string `xml:"liga,attr"`
		} `xml:"NOVO"`
	} `xml:"BDI"`
	LP struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		Activo string `xml:"activo,attr"`
		NOVO   struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Liga   string `xml:"liga,attr"`
		} `xml:"NOVO"`
		SOAPG struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Mostra string `xml:"mostra,attr"`
			Activo string `xml:"activo,attr"`
			NOVO   struct {
				Text   string `xml:",chardata"`
				Titulo string `xml:"titulo,attr"`
				Liga   string `xml:"liga,attr"`
			} `xml:"NOVO"`
			RSOAP []RSOAP `xml:"RSOAP"`
		} `xml:"SOAPG"`
		LPNT struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Mostra string `xml:"mostra,attr"`
		} `xml:"LPNT"`
		LPR struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Mostra string `xml:"mostra,attr"`
		} `xml:"LPR"`
	} `xml:"LP"`
	BN struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		NOVO   struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Liga   string `xml:"liga,attr"`
		} `xml:"NOVO"`
	} `xml:"BN"`
	DADOR struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		NOVO   []struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Liga   string `xml:"liga,attr"`
		} `xml:"NOVO"`
	} `xml:"DADOR"`
	NOTAALTA struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		NOVO   struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Liga   string `xml:"liga,attr"`
		} `xml:"NOVO"`
	} `xml:"NOTAALTA"`
	DOCS struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Mostra string `xml:"mostra,attr"`
		URG    struct {
			Text   string `xml:",chardata"`
			Titulo string `xml:"titulo,attr"`
			Mostra string `xml:"mostra,attr"`
		} `xml:"URG"`
	} `xml:"DOCS"`
}

type RSO struct {
	Text   string `xml:",chardata"`
	Titulo string `xml:"titulo,attr"`
	Valor  string `xml:"valor,attr"`
}

type RA struct {
	Text   string `xml:",chardata"`
	Titulo string `xml:"titulo,attr"`
	Valor  string `xml:"valor,attr"`
}

type RPI struct {
	Text   string `xml:",chardata"`
	Titulo string `xml:"titulo,attr"`
	Valor  string `xml:"valor,attr"`
}

type RSOAPL struct {
	Text          string `xml:",chardata"`
	Titulo        string `xml:"titulo,attr"`
	Mostra        string `xml:"mostra,attr"`
	Liga          string `xml:"liga,attr"`
	Conf          string `xml:"conf,attr"`
	Versao        string `xml:"versao,attr"`
	Data          string `xml:"data,attr"`
	Hora          string `xml:"hora,attr"`
	Autor         string `xml:"autor,attr"`
	Problema      string `xml:"problema,attr"`
	Episodio      string `xml:"episodio,attr"`
	Nome          string `xml:"nome,attr"`
	Especialidade string `xml:"especialidade,attr"`
	RSO           RSO    `xml:"RSO"`
	RA            RA     `xml:"RA"`
	RPI           RPI    `xml:"RPI"`
	IdDiario      int    `xml:"idDiario,attr"`
}

type RSOAP struct {
	Text          string   `xml:",chardata"`
	Titulo        string   `xml:"titulo,attr"`
	Mostra        string   `xml:"mostra,attr"`
	Liga          string   `xml:"liga,attr"`
	Versao        string   `xml:"versao,attr"`
	Data          string   `xml:"data,attr"`
	Hora          string   `xml:"hora,attr"`
	Autor         string   `xml:"autor,attr"`
	Problema      string   `xml:"problema,attr"`
	Episodio      string   `xml:"episodio,attr"`
	Nome          string   `xml:"nome,attr"`
	Especialidade string   `xml:"especialidade,attr"`
	RSOAPL        []RSOAPL `xml:"RSOAPL"`
}

type DOCUMENTOSAIDA struct {
	XMLName   xml.Name  `xml:"DOCUMENTOSAIDA"`
	Text      string    `xml:",chardata"`
	Servico   string    `xml:"servico,attr"`
	Nome      string    `xml:"nome,attr"`
	Processo  string    `xml:"processo,attr"`
	Episodio  string    `xml:"episodio,attr"`
	Publicado string    `xml:"publicado,attr"`
	Autor     string    `xml:"autor,attr"`
	RELATORIO RELATORIO `xml:"RELATORIO"`
}

type RELATORIO struct {
	Text           string `xml:",chardata"`
	Titulo         string `xml:"titulo,attr"`
	File           string `xml:"file,attr"`
	Mostra         string `xml:"mostra,attr"`
	Liga           string `xml:"liga,attr"`
	Versao         string `xml:"versao,attr"`
	Data           string `xml:"data,attr"`
	Hora           string `xml:"hora,attr"`
	Autor          string `xml:"autor,attr"`
	PROVENIENCIA10 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"PROVENIENCIA10"`
	PROVENIENCIA20 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"PROVENIENCIA20"`
	PROVENIENCIA30 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"PROVENIENCIA30"`
	PROVENIENCIA40 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"PROVENIENCIA40"`
	INTERNAMENTO10 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"INTERNAMENTO10"`
	REINTERNAMENTO10 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"REINTERNAMENTO10"`
	NADMSCI17 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"NADMSCI17"`
	NADMSCI18 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"NADMSCI18"`
	ALERTA01 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"ALERTA01"`
	ALERTA02 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"ALERTA02"`
	ALERTA03 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"ALERTA03"`
	ALERTA04 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"ALERTA04"`
	ALERTA05 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"ALERTA05"`
	ALERTA06 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"ALERTA06"`
	ALERTA07 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"ALERTA07"`
	ALERTA08 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"ALERTA08"`
	ENQUADRAMENTO1 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"ENQUADRAMENTO1"`
	SINDR01 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"SINDR01"`
	SINDR02 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"SINDR02"`
	SINDR03 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"SINDR03"`
	SINDROBS01 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"SINDROBS01"`
	INFECADM10 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"INFECADM10"`
	INFECADM20 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"INFECADM20"`
	INFECADM30 struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
	} `xml:"INFECADM30"`
	INFCONF struct {
		Text    string `xml:",chardata"`
		Titulo  string `xml:"titulo,attr"`
		Servico string `xml:"servico,attr"`
		Valor   string `xml:"valor,attr"`
	} `xml:"INF_CONF"`
}

type DADOSNOTAALTA struct {
	XMLName       xml.Name `xml:"DADOSNOTAALTA"`
	Text          string   `xml:",chardata"`
	Titulo        string   `xml:"titulo,attr"`
	Mostra        string   `xml:"mostra,attr"`
	Liga          string   `xml:"liga,attr"`
	Versao        string   `xml:"versao,attr"`
	Data          string   `xml:"data,attr"`
	Hora          string   `xml:"hora,attr"`
	Autor         string   `xml:"autor,attr"`
	Publica       string   `xml:"publica,attr"`
	Episodio      string   `xml:"episodio,attr"`
	Especialidade string   `xml:"especialidade,attr"`
	Interna       string   `xml:"interna,attr"`
	DTALTA        struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"DTALTA"`
	DESTINOT struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"DESTINOT"`
	DESTINOD struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"DESTINOD"`
	MTADM struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"MTADM"`
	RESUMO struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"RESUMO"`
	DIAGNO struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"DIAGNO"`
	TRATAM struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"TRATAM"`
	PROG struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"PROG"`
	PROPI struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"PROPI"`
	MONIT struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"MONIT"`
	ORIENT struct {
		Text   string `xml:",chardata"`
		Titulo string `xml:"titulo,attr"`
		Valor  string `xml:"valor,attr"`
		Texto  string `xml:"texto,attr"`
	} `xml:"ORIENT"`
}

func genericPCEXML() PCE {
	blob := `<PCE titulo="Processo Clínico Electrónico" nome="" processo="" mostra="0"> <BDI titulo="Base de Dados Integral" mostra="0"> <NOVO titulo="Nova BDI" liga="novabdi.aspx?x=BDI" /> </BDI> <LP titulo="Lista de Problemas" mostra="0" activo="-1"> <NOVO titulo="Novo Problema Activo" liga="novalp.aspx?x=LP" /> <SOAPG titulo="SOAPs Globais" mostra="0" activo="-1"> <NOVO titulo="Novo SOAP Global" liga="novasoap.aspx?x=RSOAPG" /> </SOAPG> <LPNT titulo="Problemas não Tratados" mostra="0"></LPNT> <LPR titulo="Problemas Resolvidos" mostra="0"></LPR> </LP> <BN titulo="Bloco de Notas Clínico" mostra="0"> <NOVO titulo="Novo Bloco Clínico" liga="novabn.aspx?x=BN" /> </BN> <DADOR titulo="Dador" mostra="0"> <NOVO titulo="Novo Dador Multiorgãos" liga="novadador.aspx?x=DADORM" /> <NOVO titulo="Novo Dador Corneas" liga="novadador.aspx?x=DADORC" /> </DADOR> <NOTAALTA titulo="Notas de Alta" mostra="0"> <NOVO titulo="Nova Nota de Alta" liga="nova2alta.aspx?x=NOTAALTA" /> </NOTAALTA> <DOCS titulo="Documentos" mostra="0"> <URG titulo="Notas de Alta Urgencia" mostra="0"></URG> </DOCS> </PCE>`
	var generic PCE
	if err := xml.Unmarshal([]byte(blob), &generic); err != nil {
		log.Fatal(err)
	}

	return generic
}

func genericDocumentosAIDA() DOCUMENTOSAIDA {
	blob := `<DOCUMENTOSAIDA servico="" nome="Nota de Admissão " processo="num_processo" episodio="episodio" publicado="id_medico" autor="nome do medico (NO) numero da ordem"><RELATORIO titulo="Nota de Admissão " file="defeiton_1" mostra="1" liga="" versao="1" data="sysdate" hora="hora_sistema" autor="id_medico"><PROVENIENCIA10 titulo="" valor="."/><PROVENIENCIA20 titulo="" valor=""/><PROVENIENCIA30 titulo="Outro Hospital" valor="" texto=""/><PROVENIENCIA40 titulo="Outra situação" valor="" texto=""/><INTERNAMENTO10 titulo="" valor=""/><REINTERNAMENTO10 titulo="" valor=""/><NADMSCI17 titulo="Estrangeiro" valor=""/><NADMSCI18 titulo="não fala português" valor=""/><ALERTA01 titulo="Surdo" valor=""/><ALERTA02 titulo="mudo" valor=""/><ALERTA03 titulo="cego" valor=""/><ALERTA04 titulo="Alergias" valor="" texto=""/><ALERTA05 titulo="Via aérea" valor="" texto=""/><ALERTA06 titulo="Pacemaker" valor=""/><ALERTA07 titulo="Deficiencia Fisica" valor="" texto=""/><ALERTA08 titulo="Deficiencia Psiquica" valor="" texto=""/><ENQUADRAMENTO1 titulo="" valor="colocar o texto aqui"/><SINDR01 titulo="Principal:" valor=""/><SINDR02 titulo="Secundário:" valor=""/><SINDR03 titulo="Outro:" valor=""/><SINDROBS01 titulo="" valor=""/><INFECADM10 titulo="" valor=""/><INFECADM20 titulo="" valor=""/><INFECADM30 titulo="Local da infecção: " valor=""/><INF_CONF titulo="Informação confidencial:" servico="" valor=""/></RELATORIO></DOCUMENTOSAIDA>`
	var generic DOCUMENTOSAIDA
	if err := xml.Unmarshal([]byte(blob), &generic); err != nil {
		log.Fatal(err)
	}

	return generic
}

func genericDadosNovaAlta() DADOSNOTAALTA {
	blob := `<DADOSNOTAALTA titulo="des_especialidade" mostra="1" liga="novaalta.aspx?x=DADOSNOTAALTA" versao="3" data="data_sistema" hora="hora_sistema" autor="nome_medico" publica="" episodio="episodio" especialidade="cod_especialidade" interna="0"><DTALTA titulo="Data da Alta: " valor="" texto="" /><DESTINOT titulo="" valor="" texto="" /><DESTINOD titulo="" valor="" texto="" /><MTADM titulo="" valor="" texto="" /><RESUMO titulo="" valor="colocar o texto aqui" texto="" /><DIAGNO titulo="" valor="" texto="" /><TRATAM titulo="" valor="" texto="" /><PROG titulo="" valor="" texto="" /><PROPI titulo="" valor="" texto="" /><MONIT titulo="" valor="" texto="" /><ORIENT titulo="" valor="" texto="" /></DADOSNOTAALTA>`
	var generic DADOSNOTAALTA
	if err := xml.Unmarshal([]byte(blob), &generic); err != nil {
		log.Fatal(err)
	}

	return generic
}
