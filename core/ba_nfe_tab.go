package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/glwbr/brisa/money"
	"golang.org/x/net/html"
)

// TODO: extract the common logic, such as section/sectionIndex, parseBRL, normalizeText, normalizeDigits, etc. to a separate file
var ErrBANFeTabNotFound = errors.New("ba nfe tab not found")

const (
	baSectionDados        = "Dados da NFC-e"
	baSectionEmitente     = "Emitente"
	baSectionDestinatario = "Destinatário"
)

// ParseBAFromAbasNFe parses the "Abas NFe" layout from SEFAZ-BA into a Receipt.
// htmlBytes should be the HTML of the NFC-e details page.
func ParseBAFromAbasNFe(htmlBytes []byte) (*Receipt, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	fmt.Println(string(htmlBytes))

	nfe := doc.Find("#NFe")
	if nfe.Length() == 0 {
		return nil, ErrBANFeTabNotFound
	}

	root := nfe
	sections := buildBASectionIndex(root)
	dados := sections[baSectionDados]
	emitente := sections[baSectionEmitente]
	destinatario := sections[baSectionDestinatario]

	totalStr := sections.firstValue("Valor Total da Nota Fiscal")
	if totalStr == "" {
		totalStr = sections.firstValue("Valor Total")
	}
	total, _ := parseBRL(totalStr)
	subtotal := total
	discount := money.BRL(0)

	r := &Receipt{
		Key:           normalizeDigits(doc.Find("#lbl_chave_acesso").Text()),
		Portal:        PortalBA,
		Series:        strings.TrimSpace(dados["Série"]),
		ReceiptNumber: strings.TrimSpace(dados["Número"]),

		Issuer: Issuer{
			Name: strings.TrimSpace(emitente["Nome / Razão Social"]),
			CNPJ: normalizeDigits(emitente["CNPJ"]),
			Address: Address{
				State: strings.TrimSpace(emitente["UF"]),
			},
			StateRegID: strings.TrimSpace(
				normalizeDigits(emitente["Inscrição Estadual"]),
			),
		},

		Consumer: Consumer{
			Document: normalizeDigits(destinatario["CPF"]),
			Name:     strings.TrimSpace(destinatario["Nome"]),
		},

		Items:    []Item{},
		Subtotal: subtotal,
		Discount: discount,
		Total:    total,

		RawHTML: htmlBytes,
	}

	if issueDate := strings.TrimSpace(dados["Data de Emissão"]); issueDate != "" {
		if ts, err := parseBADate(issueDate); err == nil {
			r.IssueDate = ts
		}
	}

	// Items are not listed in the NFe tab; they can be attached by parsing the
	// Produtos/Serviços tab separately and merging later.
	r.Items = []Item{}

	return r, nil
}

type Sections map[string]map[string]string

func (s Sections) firstValue(label string) string {
	label = normalizeText(label)
	for _, section := range s {
		if val, ok := section[label]; ok && strings.TrimSpace(val) != "" {
			return strings.TrimSpace(val)
		}
	}
	return ""
}

func buildBASectionIndex(nfe *goquery.Selection) Sections {
	sections := Sections{}
	textCache := map[*html.Node]string{}

	currentSection := ""
	nfe.Find("table").Each(func(_ int, table *goquery.Selection) {
		if title := extractSectionTitle(table, textCache); title != "" {
			currentSection = title
			if _, ok := sections[currentSection]; !ok {
				sections[currentSection] = map[string]string{}
			}
			return
		}

		if currentSection == "" {
			return
		}

		for label, value := range collectLabelValues(table, textCache) {
			if value == "" {
				continue
			}
			if _, exists := sections[currentSection][label]; !exists {
				sections[currentSection][label] = value
			}
		}
	})

	return sections
}

func extractSectionTitle(table *goquery.Selection, cache map[*html.Node]string) string {
	title := table.Find("td.table-titulo-aba, td.table-titulo-aba-interna").First()
	return cachedText(title, cache)
}

func collectLabelValues(sel *goquery.Selection, cache map[*html.Node]string) map[string]string {
	values := map[string]string{}

	var walk func(*goquery.Selection)
	walk = func(node *goquery.Selection) {
		node.Children().Each(func(_ int, child *goquery.Selection) {
			if goquery.NodeName(child) == "label" {
				label := cachedText(child, cache)
				valSel := child.Next()
				for valSel.Length() > 0 && goquery.NodeName(valSel) != "span" {
					valSel = valSel.Next()
				}
				if value := cachedText(valSel, cache); label != "" && value != "" {
					if _, ok := values[label]; !ok {
						values[label] = value
					}
				}
			}
			walk(child)
		})
	}

	walk(sel)
	return values
}

func cachedText(sel *goquery.Selection, cache map[*html.Node]string) string {
	if sel.Length() == 0 {
		return ""
	}

	node := sel.Get(0)
	if node == nil {
		return ""
	}

	if cached, ok := cache[node]; ok {
		return cached
	}

	text := normalizeText(sel.Text())
	cache[node] = text
	return text
}

func parseBADate(raw string) (time.Time, error) {
	raw = strings.ReplaceAll(raw, "\u00A0", " ")
	raw = strings.TrimSpace(raw)

	layouts := []string{
		"02/01/2006 15:04:05-07:00",
		"02/01/2006 15:04:05",
	}

	for _, layout := range layouts {
		if ts, err := time.Parse(layout, raw); err == nil {
			return ts, nil
		}
	}

	return time.Time{}, errors.New("invalid BA date format")
}

func normalizeDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}

// parseBRL is a tiny adapter around your money.Parse to keep this file tidy.
func parseBRL(s string) (money.BRL, error) {
	return money.Parse(s)
}
