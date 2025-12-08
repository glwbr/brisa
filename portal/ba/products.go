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

var ErrProductsTabNotFound = errors.New("products tab not found")

func ParseProductsTab(htmlBytes []byte) ([]invoice.Item, error) {
	doc, err := scraper.ParseHTML(htmlBytes)
	if err != nil {
		return nil, err
	}

	prod := doc.Find("#Prod")
	if prod.Length() == 0 {
		return nil, ErrProductsTabNotFound
	}

	var items []invoice.Item
	prod.Find("td.table_produtos").Each(func(_ int, td *goquery.Selection) {
		cache := map[*html.Node]string{}

		summary := td.Find("table.toggle").First()
		if summary.Length() == 0 {
			return
		}
		summaryVals := scraper.CollectLabelValues(summary, cache)

		detail := td.Find("table.toggable").First()
		detailVals := map[string]string{}
		if detail.Length() > 0 {
			detailVals = scraper.CollectLabelValues(detail, cache)
		}

		item := invoice.Item{
			LineNumber:  parse.Int(summaryVals["Número"]),
			Description: summaryVals["Descrição"],
			Quantity:    parse.Quantity(parse.FirstNonEmpty(summaryVals["Qtd."], detailVals["Quantidade Comercial"])),
			Unit:        invoice.ParseUnit(strings.ToUpper(parse.FirstNonEmpty(summaryVals["Unidade Comercial"], detailVals["Unidade Comercial"]))),
			Total:       parseMoneyOrZero(parse.FirstNonEmpty(summaryVals["Valor (R$)"], detailVals["Valor Total"])),
			Code:        detailVals["Código do Produto"],
			NCM:         detailVals["Código NCM"],
			CEST:        detailVals["Código CEST"],
			CFOP:        detailVals["CFOP"],
		}

		if gtin := parse.FirstNonEmpty(detailVals["Código EAN Comercial"], detailVals["Código EAN Tributável"]); gtin != "" && !strings.EqualFold(gtin, "SEM GTIN") {
			item.GTIN = parse.Digits(gtin)
		}

		if up := parse.FirstNonEmpty(detailVals["Valor unitário de comercialização"], detailVals["Valor unitário de tributação"]); up != "" {
			if price, err := money.Parse(up); err == nil && price != 0 {
				item.UnitPrice = price
			}
		}

		if item.UnitPrice == 0 && item.Quantity > 0 && item.Total != 0 {
			item.UnitPrice = money.FromFloat(item.Total.Float64() / item.Quantity)
		}

		item.Taxes = parseTaxes(detail, detailVals)
		items = append(items, item)
	})

	return items, nil
}

func parseMoneyOrZero(s string) money.BRL {
	v, _ := money.Parse(s)
	return v
}

func parseTaxes(detail *goquery.Selection, vals map[string]string) *invoice.Taxes {
	taxAmount := parseMoneyOrZero(vals["Valor Aproximado dos Tributos"])

	icmsPercent := parse.Percent(parse.FirstNonEmpty(
		vals["Alíquota do ICMS Normal"],
		vals["Alíquota do ICMS"],
	))
	icmsValue := parseMoneyOrZero(parse.FirstNonEmpty(
		vals["Valor do ICMS Normal"],
		vals["Valor do ICMS ST"],
	))

	pisPercent, pisValue := parseTaxSection(detail, "PIS")
	cofinsPercent, cofinsValue := parseTaxSection(detail, "COFINS")

	total := taxAmount.Add(icmsValue).Add(pisValue).Add(cofinsValue)

	if icmsPercent == 0 && pisPercent == 0 && cofinsPercent == 0 && total == 0 {
		return nil
	}

	return &invoice.Taxes{
		ICMSPercent:   icmsPercent,
		PISPercent:    pisPercent,
		COFINSPercent: cofinsPercent,
		Amount:        total,
	}
}

func parseTaxSection(detail *goquery.Selection, titleSubstr string) (float64, money.BRL) {
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

	if goquery.NodeName(next) == "div" {
		next = next.Find("table").First()
	}
	if next.Length() == 0 {
		return 0, 0
	}

	cache := map[*html.Node]string{}
	vals := scraper.CollectLabelValues(next, cache)

	return parse.Percent(vals["Alíquota"]), parseMoneyOrZero(vals["Valor"])
}
