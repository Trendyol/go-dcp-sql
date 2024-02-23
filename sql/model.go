package sql

type Model interface {
	Convert() string
}

type Raw struct {
	Query string
}

func (u *Raw) Convert() string {
	return u.Query
}
