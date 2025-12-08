package core

import "github.com/glwbr/brisa/money"

type Taxes struct {
	IPIPercent    float64 `json:"ipi_percent,omitempty"`
	PISPercent    float64 `json:"pis_percent,omitempty"`
	ICMSPercent   float64 `json:"icms_percent,omitempty"`
	COFINSPercent float64 `json:"cofins_percent,omitempty"`

	Amount money.BRL `json:"amount"`
}

type TaxTotals struct {
	IPI    money.BRL `json:"ipi,omitempty"`
	PIS    money.BRL `json:"pis,omitempty"`
	ICMS   money.BRL `json:"icms,omitempty"`
	COFINS money.BRL `json:"cofins,omitempty"`
	Other  money.BRL `json:"other,omitempty"`
}
