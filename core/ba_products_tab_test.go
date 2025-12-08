package core

import (
	"path/filepath"
	"testing"
)

func TestParseBAFromAbasProdutosServicosItems(t *testing.T) {
	html := mustRead(t, filepath.Join("..", "testdata", "assai_abas_produtos_servicos.html"))

	items, err := ParseBAFromAbasProdutosServicos(html)
	if err != nil {
		t.Fatalf("ParseBAFromAbasProdutosServicos returned error: %v", err)
	}

	if len(items) == 0 {
		t.Fatalf("items length = 0; want parsed products")
	}

	hasCode := false
	hasNCM := false
	prevLine := 0
	hasTaxes := false

	for _, item := range items {
		if item.LineNumber <= 0 {
			t.Fatalf("item line number <= 0: %+v", item)
		}
		if item.LineNumber < prevLine {
			t.Fatalf("item line numbers not ascending: %d then %d", prevLine, item.LineNumber)
		}
		prevLine = item.LineNumber

		if item.Description == "" {
			t.Fatalf("item description empty for line %d", item.LineNumber)
		}
		if item.Quantity <= 0 {
			t.Fatalf("item quantity <= 0 for line %d", item.LineNumber)
		}
		if item.Unit == "" {
			t.Fatalf("item unit empty for line %d", item.LineNumber)
		}
		if item.Total == 0 {
			t.Fatalf("item total is zero for line %d", item.LineNumber)
		}
		if item.UnitPrice == 0 {
			t.Fatalf("item unit price is zero for line %d", item.LineNumber)
		}
		if item.Taxes != nil {
			hasTaxes = true
		}
		if item.Code != "" {
			hasCode = true
		}
		if item.NCM != "" {
			hasNCM = true
		}
	}

	if !hasCode {
		t.Fatalf("expected at least one item with product code")
	}
	if !hasNCM {
		t.Fatalf("expected at least one item with NCM code")
	}
	if !hasTaxes {
		t.Fatalf("expected at least one item with tax data")
	}
}
