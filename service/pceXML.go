package service

import "encoding/xml"

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
			RSOAP struct {
				Text          string `xml:",chardata"`
				Titulo        string `xml:"titulo,attr"`
				Mostra        string `xml:"mostra,attr"`
				Liga          string `xml:"liga,attr"`
				Versao        string `xml:"versao,attr"`
				Data          string `xml:"data,attr"`
				Hora          string `xml:"hora,attr"`
				Autor         string `xml:"autor,attr"`
				Problema      string `xml:"problema,attr"`
				Episodio      string `xml:"episodio,attr"`
				Nome          string `xml:"nome,attr"`
				Especialidade string `xml:"especialidade,attr"`
				RSOAPL        struct {
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
					RSO           struct {
						Text   string `xml:",chardata"`
						Titulo string `xml:"titulo,attr"`
						Valor  string `xml:"valor,attr"`
					} `xml:"RSO"`
					RA struct {
						Text   string `xml:",chardata"`
						Titulo string `xml:"titulo,attr"`
						Valor  string `xml:"valor,attr"`
					} `xml:"RA"`
					RPI struct {
						Text   string `xml:",chardata"`
						Titulo string `xml:"titulo,attr"`
						Valor  string `xml:"valor,attr"`
					} `xml:"RPI"`
				} `xml:"RSOAPL"`
			} `xml:"RSOAP"`
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

func genericPCEXML() PCE {
	return PCE{}
}
