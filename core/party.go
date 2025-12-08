package core

type Issuer struct {
	Name        string `json:"name"`
	CNPJ        string `json:"cnpj"`
	TradeName   string `json:"trade_name,omitempty"`
	StateRegID  string `json:"state_reg_id,omitempty"`
	MunicipalID string `json:"municipal_id,omitempty"`

	Address Address `json:"address"`
}

type Consumer struct {
	Document string `json:"document,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	ZipCode    string `json:"zip_code,omitempty"`
}
