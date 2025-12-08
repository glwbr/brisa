package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/glwbr/brisa/money"
	"golang.org/x/net/html"
)

var ErrBAProdutosTabNotFound = errors.New("ba produtos/servicos tab not found")

// ParseBAFromAbasProdutosServicos parses the Produtos/Serviços tab into a slice of Items.
// It does not populate header data; merge the result into a Receipt parsed from the NFe tab.
func ParseBAFromAbasProdutosServicos(htmlBytes []byte) ([]Item, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	prod := doc.Find("#Prod")
	if prod.Length() == 0 {
		return nil, ErrBAProdutosTabNotFound
	}

	return parseBAItemsFromProd(doc), nil
}

func parseBAItemsFromProd(doc *goquery.Document) []Item {
	prod := doc.Find("#Prod")
	if prod.Length() == 0 {
		return nil
	}

	items := []Item{}

	prod.Find("td.table_produtos").Each(func(_ int, td *goquery.Selection) {
		cache := map[*html.Node]string{}

		summary := td.Find("table.toggle").First()
		if summary.Length() == 0 {
			return
		}
		summaryVals := collectLabelValues(summary, cache)

		detail := td.Find("table.toggable").First()
		detailVals := map[string]string{}
		if detail.Length() > 0 {
			detailVals = collectLabelValues(detail, cache)
		}

		item := Item{
			LineNumber:  parseLineNumber(summaryVals["Número"]),
			Description: summaryVals["Descrição"],
			Quantity: parseQuantity(
				firstNonEmpty(summaryVals["Qtd."], detailVals["Quantidade Comercial"], detailVals["Quantidade Tributável"]),
			),
			Unit: normalizeUnit(
				firstNonEmpty(summaryVals["Unidade Comercial"], detailVals["Unidade Comercial"], detailVals["Unidade Tributável"]),
			),
			Total: parseMoneyOrZero(firstNonEmpty(summaryVals["Valor (R$)"], detailVals["Valor Total"])),
			Code:  detailVals["Código do Produto"],
			NCM:   detailVals["Código NCM"],
			CEST:  detailVals["Código CEST"],
			CFOP:  detailVals["CFOP"],
		}

		if gtin := firstNonEmpty(detailVals["Código EAN Comercial"], detailVals["Código EAN Tributável"]); gtin != "" && !strings.EqualFold(gtin, "SEM GTIN") {
			item.GTIN = normalizeDigits(gtin)
		}

		if unitPriceStr := firstNonEmpty(detailVals["Valor unitário de comercialização"], detailVals["Valor unitário de tributação"]); unitPriceStr != "" {
			if up, err := parseBRL(unitPriceStr); err == nil && up != 0 {
				item.UnitPrice = up
			}
		}

		if item.UnitPrice == 0 && item.Quantity > 0 && item.Total != 0 {
			item.UnitPrice = money.FromFloat(item.Total.Float64() / item.Quantity)
		}

		item.Taxes = parseBATaxes(detail, detailVals)

		items = append(items, item)
	})

	return items
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func parseMoneyOrZero(s string) money.BRL {
	v, err := parseBRL(s)
	if err != nil {
		return 0
	}
	return v
}

func parseQuantity(s string) float64 {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseLineNumber(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func normalizeUnit(s string) Unit {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "KG":
		return UnitKilogram
	case "G":
		return UnitGram
	case "L":
		return UnitLiter
	case "UN", "UND", "UNID":
		return UnitUnit
	case "M":
		return UnitMeter
	default:
		return Unit(s)
	}
}

func parsePercent(s string) float64 {
	s = strings.ReplaceAll(s, "%", "")
	return parseQuantity(s)
}

func parseBATaxes(detail *goquery.Selection, detailVals map[string]string) *Taxes {
	taxAmount := parseMoneyOrZero(detailVals["Valor Aproximado dos Tributos"])

	icmsPercent := parsePercent(firstNonEmpty(
		detailVals["Alíquota do ICMS Normal"],
		detailVals["Alíquota do ICMS"],
		detailVals["Alíquota do ICMS ST"],
	))
	icmsValue := parseMoneyOrZero(firstNonEmpty(
		detailVals["Valor do ICMS Normal"],
		detailVals["Valor do ICMS ST retido"],
		detailVals["Valor do ICMS ST"],
	))

	pisPercent, pisValue := parseTaxSection(detail, "PIS")
	cofinsPercent, cofinsValue := parseTaxSection(detail, "COFINS")

	totalValue := taxAmount
	totalValue = totalValue.Add(icmsValue).Add(pisValue).Add(cofinsValue)

	if icmsPercent == 0 && pisPercent == 0 && cofinsPercent == 0 && totalValue == 0 {
		return nil
	}

	return &Taxes{
		ICMSPercent:   icmsPercent,
		PISPercent:    pisPercent,
		COFINSPercent: cofinsPercent,
		Amount:        totalValue,
	}
}

func parseTaxSection(detail *goquery.Selection, titleSubstr string) (percent float64, amount money.BRL) {
	title := detail.Find("td.table-titulo-aba-interna").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(strings.ToUpper(s.Text()), strings.ToUpper(titleSubstr))
	}).First()
	if title.Length() == 0 {
		return 0, 0
	}

	titleTable := title.ParentsFiltered("table").First()
	if titleTable.Length() == 0 {
		return 0, 0
	}

	next := titleTable.Next()
	for next.Length() > 0 && goquery.NodeName(next) != "table" && goquery.NodeName(next) != "div" {
		next = next.Next()
	}
	if next.Length() == 0 {
		return 0, 0
	}

	// When wrapped in a div.toggable, the useful table is inside it.
	if goquery.NodeName(next) == "div" {
		next = next.Find("table").First()
	}

	if next.Length() == 0 {
		return 0, 0
	}

	cache := map[*html.Node]string{}
	vals := collectLabelValues(next, cache)

	percent = parsePercent(firstNonEmpty(vals["Alíquota"], vals["Alíquota do ICMS Normal"]))
	amount = parseMoneyOrZero(vals["Valor"])

	return percent, amount
}
