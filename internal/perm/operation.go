package perm

type Operation uint8

const (
	OpRead Operation = iota
	OpWrite
	OpExec
)

func (o Operation) String() string {
	switch o {
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpExec:
		return "exec"
	}
	panic("unknown operation")
}

func (o Operation) Short() string {
	switch o {
	case OpRead:
		return "r"
	case OpWrite:
		return "w"
	case OpExec:
		return "x"
	}
	panic("unknown operation")
}
