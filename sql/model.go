package sql

type Model interface {
	Convert() *ExecArgs
}

type Raw struct {
	Query string
	Args  []interface{}
}

type ExecArgs struct {
	Query string
	Args  []interface{}
}

func (u *Raw) Convert() *ExecArgs {
	return &ExecArgs{
		Query: u.Query,
		Args:  u.Args,
	}
}
