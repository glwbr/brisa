package core

type BAParser struct{}

func init() {
	RegisterParser(BAParser{})
}

func (p BAParser) Portal() Portal {
	return PortalBA
}

func (p BAParser) Parse(html []byte, opts ParseOptions) (*Receipt, error) {
	return ParseBAFromAbasNFe(html)
}
