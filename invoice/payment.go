package invoice

import "github.com/glwbr/brisa/money"

type PaymentMethod string

const (
	PaymentCash        PaymentMethod = "cash"
	PaymentCreditCard  PaymentMethod = "credit_card"
	PaymentDebitCard   PaymentMethod = "debit_card"
	PaymentPix         PaymentMethod = "pix"
	PaymentBankSlip    PaymentMethod = "boleto"
	PaymentStoreCredit PaymentMethod = "store_credit"
	PaymentOther       PaymentMethod = "other"
)

type Payment struct {
	Method       PaymentMethod `json:"method"`
	Amount       money.BRL     `json:"amount"`
	Installments int           `json:"installments,omitempty"`
}
