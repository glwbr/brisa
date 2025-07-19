package invoice

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/glwbr/brisa/pkg/money"
)

// Invoice represents a complete electronic invoice
type Invoice struct {
	AccessKey string    `json:"access_key"`
	IssueDate time.Time `json:"issue_date"`
	Issuer    Issuer    `json:"issuer"`
	Customer  Customer  `json:"customer"`
	Items     []Item    `json:"items"`
	Total     money.BRL `json:"total"`
	Payments  []Payment `json:"payments"`
	Taxes     Taxes     `json:"taxes"`
}

// Issuer represents the invoice issuer (store/company)
type Issuer struct {
	Name        string  `json:"name"`
	CNPJ        string  `json:"cnpj"`
	TradeName   string  `json:"trade_name,omitempty"`
	Address     Address `json:"address"`
	StateRegID  string  `json:"state_reg_id,omitempty"` // Inscrição Estadual
	MunicipalID string  `json:"municipal_id,omitempty"` // Inscrição Municipal
}

// Customer represents the invoice customer
type Customer struct {
	Document string   `json:"document,omitempty"`
	Name     string   `json:"name,omitempty"`
	Address  *Address `json:"address,omitempty"`
}

// Item represents a product or service in the invoice
type Item struct {
	Code        string    `json:"code"`
	Description string    `json:"description"`
	NCM         string    `json:"ncm,omitempty"` // Nomenclatura Comum do Mercosul
	Unit        Unit      `json:"unit"`
	Details     string    `json:"details"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   money.BRL `json:"unit_price"`
	Taxes       *Taxes    `json:"taxes,omitempty"`
	Total       money.BRL `json:"total"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
}

// Taxes holds tax percentages and their total amount
type Taxes struct {
	IPI    float64   `json:"ipi,omitempty"` // percent
	PIS    float64   `json:"pis"`           // percent
	ICMS   float64   `json:"icms"`          // percent
	COFINS float64   `json:"cofins"`        // percent
	Amount money.BRL `json:"amount"`
}

// Payment holds payment method and value
type Payment struct {
	Method       PaymentMethod `json:"method"`
	Installments int           `json:"installments,omitempty"` // Optional; typically 1 for most invoices
	Amount       money.BRL     `json:"amount"`
}

type (
	PaymentMethod string
	Unit          string
)

// Constants for known unit types
const (
	UnitKilogram Unit = "Kg"
	UnitGram     Unit = "g"
	UnitLiter    Unit = "L"
	UnitUnit     Unit = "Un"
	UnitMeter    Unit = "m"
)

// Constants for known payment methods
const (
	PaymentCash        PaymentMethod = "cash"
	PaymentCreditCard  PaymentMethod = "credit_card"
	PaymentDebitCard   PaymentMethod = "debit_card"
	PaymentPix         PaymentMethod = "pix"
	PaymentBankSlip    PaymentMethod = "boleto"
	PaymentStoreCredit PaymentMethod = "store_credit"
)

// String returns a JSON string representation of the invoice
func (i *Invoice) String() string {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling invoice: %v", err)
	}
	return string(data)
}

// SaveToFile saves the invoice to a file in JSON format
func (i *Invoice) SaveToFile(filepath string) error {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling invoice: %w", err)
	}

	// Create file with restrictive permissions (only owner can read/write)
	err = os.WriteFile(filepath, data, 0o600)
	if err != nil {
		return fmt.Errorf("error writing invoice to file: %w", err)
	}

	return nil
}
