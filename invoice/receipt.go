// Package invoice defines domain models for Brazilian electronic invoices.
package invoice

import (
	"time"

	"github.com/glwbr/brisa/money"
)

type Portal string

const (
	PortalBA Portal = "BA"
	PortalCE Portal = "CE"
)

func (p Portal) String() string { return string(p) }

type Receipt struct {
	Key           string    `json:"key"`
	Portal        Portal    `json:"portal"`
	IssueDate     time.Time `json:"issue_date"`
	ReceiptNumber string    `json:"receipt_number,omitempty"`
	Series        string    `json:"series,omitempty"`

	Issuer   Issuer   `json:"issuer"`
	Consumer Consumer `json:"consumer"`
	Items    []Item   `json:"items"`

	Subtotal money.BRL `json:"subtotal"`
	Discount money.BRL `json:"discount"`
	Total    money.BRL `json:"total"`

	Payments []Payment `json:"payments,omitempty"`
	Taxes    Taxes     `json:"taxes"`

	RawHTML []byte `json:"-"`
}
