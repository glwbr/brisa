package core

import "github.com/glwbr/brisa/money"

type Unit string

const (
	UnitKilogram Unit = "KG"
	UnitGram     Unit = "G"
	UnitLiter    Unit = "L"
	UnitUnit     Unit = "UN"
	UnitMeter    Unit = "M"
)

type Item struct {
	LineNumber  int    `json:"line_number"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description"`
	Details     string `json:"details,omitempty"`

	Quantity  float64   `json:"quantity"`
	Unit      Unit      `json:"unit"`
	UnitPrice money.BRL `json:"unit_price"`
	Total     money.BRL `json:"total"`

	NCM  string `json:"ncm,omitempty"`
	GTIN string `json:"gtin,omitempty"`
	CFOP string `json:"cfop,omitempty"`
	CEST string `json:"cest,omitempty"`

	Taxes *Taxes `json:"taxes,omitempty"`
}
