package ba

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/glwbr/brisa/invoice"
	"github.com/glwbr/brisa/money"
	"github.com/glwbr/brisa/parse"
	"github.com/glwbr/brisa/scraper"
	"golang.org/x/net/html"
)

var ErrNFeTabNotFound = errors.New("nfe tab not found")

const (
	sectionDados        = "Dados da NFC-e"
	sectionEmitente     = "Emitente"
	sectionDestinatario = "Destinatário"
)

func ParseNFeTab(htmlBytes []byte) (*invoice.Receipt, error) {
	doc, err := scraper.ParseHTML(htmlBytes)
	if err != nil {
		return nil, err
	}

	nfe := doc.Find("#NFe")
	if nfe.Length() == 0 {
		return nil, ErrNFeTabNotFound
	}

	sections := buildSectionIndex(nfe)
	dados := sections[sectionDados]
	emitente := sections[sectionEmitente]
	destinatario := sections[sectionDestinatario]

	totalStr := sections.firstValue("Valor Total da Nota Fiscal")
	if totalStr == "" {
		totalStr = sections.firstValue("Valor Total")
	}
	total, _ := money.Parse(totalStr)

	r := &invoice.Receipt{
		Key:           parse.Digits(doc.Text("#lbl_chave_acesso")),
		Portal:        invoice.PortalBA,
		Series:        strings.TrimSpace(dados["Série"]),
		ReceiptNumber: strings.TrimSpace(dados["Número"]),
		Issuer: invoice.Issuer{
			Name:       strings.TrimSpace(emitente["Nome / Razão Social"]),
			CNPJ:       parse.Digits(emitente["CNPJ"]),
			StateRegID: parse.Digits(emitente["Inscrição Estadual"]),
			Address:    invoice.Address{State: strings.TrimSpace(emitente["UF"])},
		},
		Consumer: invoice.Consumer{
			Document: parse.Digits(destinatario["CPF"]),
			Name:     strings.TrimSpace(destinatario["Nome"]),
		},
		Subtotal: total,
		Total:    total,
		RawHTML:  htmlBytes,
	}

	if issueDate := strings.TrimSpace(dados["Data de Emissão"]); issueDate != "" {
		if ts, err := parse.BrazilianDate(issueDate); err == nil {
			r.IssueDate = ts
		}
	}

	return r, nil
}

type sections map[string]map[string]string

func (s sections) firstValue(label string) string {
	label = parse.Text(label)
	for _, section := range s {
		if val, ok := section[label]; ok && strings.TrimSpace(val) != "" {
			return strings.TrimSpace(val)
		}
	}
	return ""
}

func buildSectionIndex(nfe *goquery.Selection) sections {
	s := sections{}
	cache := map[*html.Node]string{}
	current := ""

	nfe.Find("table").Each(func(_ int, table *goquery.Selection) {
		if title := extractSectionTitle(table, cache); title != "" {
			current = title
			if _, ok := s[current]; !ok {
				s[current] = map[string]string{}
			}
			return
		}
		if current == "" {
			return
		}
		for label, value := range scraper.CollectLabelValues(table, cache) {
			if value != "" {
				if _, exists := s[current][label]; !exists {
					s[current][label] = value
				}
			}
		}
	})
	return s
}

func extractSectionTitle(table *goquery.Selection, cache map[*html.Node]string) string {
	title := table.Find("td.table-titulo-aba, td.table-titulo-aba-interna").First()
	return scraper.CachedText(title, cache)
}
