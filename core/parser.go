package core

import (
	"errors"
	"fmt"
)

type ParseOptions struct {
	Portal Portal
}

var (
	ErrMissingPortal     = errors.New("missing portal")
	ErrUnsupportedPortal = errors.New("unsupported portal")
)

type Parser interface {
	Portal() Portal
	Parse(html []byte, opts ParseOptions) (*Receipt, error)
}

var parserRegistry = map[Portal]Parser{}

func RegisterParser(p Parser) {
	parserRegistry[p.Portal()] = p
}

func ParseFromHTML(html []byte, opts ParseOptions) (*Receipt, error) {
	if opts.Portal == "" {
		return nil, ErrMissingPortal
	}

	p, ok := parserRegistry[opts.Portal]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedPortal, opts.Portal)
	}

	return p.Parse(html, opts)
}
