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

func genericPCEXML() PCE {
	blob := `<PCE titulo="Processo Clínico Electrónico" nome="" processo="" mostra="0"> <BDI titulo="Base de Dados Integral" mostra="0"> <NOVO titulo="Nova BDI" liga="novabdi.aspx?x=BDI" /> </BDI> <LP titulo="Lista de Problemas" mostra="0" activo="-1"> <NOVO titulo="Novo Problema Activo" liga="novalp.aspx?x=LP" /> <SOAPG titulo="SOAPs Globais" mostra="0" activo="-1"> <NOVO titulo="Novo SOAP Global" liga="novasoap.aspx?x=RSOAPG" /> </SOAPG> <LPNT titulo="Problemas não Tratados" mostra="0"></LPNT> <LPR titulo="Problemas Resolvidos" mostra="0"></LPR> </LP> <BN titulo="Bloco de Notas Clínico" mostra="0"> <NOVO titulo="Novo Bloco Clínico" liga="novabn.aspx?x=BN" /> </BN> <DADOR titulo="Dador" mostra="0"> <NOVO titulo="Novo Dador Multiorgãos" liga="novadador.aspx?x=DADORM" /> <NOVO titulo="Novo Dador Corneas" liga="novadador.aspx?x=DADORC" /> </DADOR> <NOTAALTA titulo="Notas de Alta" mostra="0"> <NOVO titulo="Nova Nota de Alta" liga="nova2alta.aspx?x=NOTAALTA" /> </NOTAALTA> <DOCS titulo="Documentos" mostra="0"> <URG titulo="Notas de Alta Urgencia" mostra="0"></URG> </DOCS> </PCE>`
	var generic PCE
	if err := xml.Unmarshal([]byte(blob), &generic); err != nil {
		log.Fatal(err)
	}

	return generic
}
