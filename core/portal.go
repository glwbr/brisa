package core

type Portal string

const (
	PortalBA Portal = "BA"
	PortalCE Portal = "CE"
)

func (p Portal) String() string {
	return string(p)
}
