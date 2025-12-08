package core

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseBAFromAbasNFe(t *testing.T) {
	html := mustRead(t, filepath.Join("..", "testdata", "assai_abas_nfe.html"))

	r, err := ParseBAFromAbasNFe(html)
	if err != nil {
		t.Fatalf("ParseBAFromAbasNFe returned error: %v", err)
	}

	if got := r.Key; len(got) != 44 {
		t.Fatalf("receipt key length = %d; want 44", len(got))
	}

	if got := r.Issuer.Name; got == "" {
		t.Fatalf("issuer name is empty; want non-empty")
	}

	if got := r.Issuer.CNPJ; len(got) != 14 {
		t.Fatalf("issuer cnpj length = %d; want 14 digits", len(got))
	}

	if got := r.Issuer.Address.State; got == "" {
		t.Fatalf("issuer state is empty; want non-empty")
	}

	if got := r.Consumer.Document; got != "" && len(got) != 11 && len(got) != 14 {
		t.Fatalf("consumer document length = %d; want 11 or 14 digits (cpf/cnpj)", len(got))
	}

	if got := r.Series; got == "" {
		t.Fatalf("series is empty; want non-empty")
	}

	if got := r.ReceiptNumber; got == "" {
		t.Fatalf("receipt number is empty; want non-empty")
	}

	if got := r.IssueDate; got.IsZero() {
		t.Fatalf("issue date is zero; want parsed timestamp")
	}

	if got := r.Total; got == 0 {
		t.Fatalf("total is zero; want parsed value")
	}

	if got := r.Subtotal; got == 0 {
		t.Fatalf("subtotal is zero; want parsed value")
	}

	if got := len(r.Items); got != 0 {
		t.Fatalf("items length = %d; want 0 (NFe tab has no products)", got)
	}
}

func TestBuildBASectionIndexRequiredLabels(t *testing.T) {
	sections := mustBuildSections(t)

	expected := map[string][]struct {
		label         string
		mustHaveValue bool
	}{
		"Dados da NFC-e": {
			{label: "Modelo", mustHaveValue: true},
			{label: "Série", mustHaveValue: true},
			{label: "Número", mustHaveValue: true},
			{label: "Data de Emissão", mustHaveValue: true},
			{label: "Valor Total da Nota Fiscal", mustHaveValue: true},
		},
		"Emitente": {
			{label: "CNPJ", mustHaveValue: true},
			{label: "Nome / Razão Social", mustHaveValue: true},
			{label: "Inscrição Estadual", mustHaveValue: true},
			{label: "UF", mustHaveValue: true},
		},
		"Destinatário": {
			{label: "CPF", mustHaveValue: false}, // may be blank
		},
		"Emissão": {
			{label: "Processo", mustHaveValue: true},
			{label: "Versão do Processo", mustHaveValue: true},
			{label: "Tipo de Emissão", mustHaveValue: true},
			{label: "Finalidade", mustHaveValue: true},
			{label: "Natureza da Operação", mustHaveValue: true},
			{label: "Indicador de Intermediador/Marketplace", mustHaveValue: true},
			{label: "Tipo da Operação", mustHaveValue: true},
			{label: "Digest Value da NF-e", mustHaveValue: true},
		},
	}

	for sectionName, labels := range expected {
		section, ok := sections[sectionName]
		if !ok {
			t.Fatalf("missing section %q in parsed markup", sectionName)
		}
		for _, expectedLabel := range labels {
			label := normalizeText(expectedLabel.label)
			got, ok := section[label]
			if !ok {
				t.Fatalf("section %q missing label %q", sectionName, label)
			}
			if expectedLabel.mustHaveValue && got == "" {
				t.Fatalf("section %q label %q is empty; want non-empty value", sectionName, label)
			}
		}
	}
}

func TestParseFromHTMLUsesRegisteredParser(t *testing.T) {
	html := mustRead(t, filepath.Join("..", "testdata", "assai_abas_nfe.html"))

	r, err := ParseFromHTML(html, ParseOptions{Portal: PortalBA})
	if err != nil {
		t.Fatalf("ParseFromHTML returned error: %v", err)
	}

	if r.Portal != PortalBA {
		t.Fatalf("parsed portal = %s; want %s", r.Portal, PortalBA)
	}
}

func TestParseBAFromAbasNFeErrorsWithoutTab(t *testing.T) {
	_, err := ParseBAFromAbasNFe([]byte("<html><body><div id=\"Other\"></div></body></html>"))
	if err == nil {
		t.Fatalf("expected error when NFe tab is missing")
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func mustBuildSections(t *testing.T) Sections {
	t.Helper()
	html := mustRead(t, filepath.Join("..", "testdata", "assai_abas_nfe.html"))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}
	nfe := doc.Find("#NFe")
	if nfe.Length() == 0 {
		t.Fatalf("missing NFe tab in test fixture")
	}
	return buildBASectionIndex(nfe)
}
