package tidio

import "github.com/gregoryv/rs"

// NewShell returns a new shell. Once an error occurs the shell no
// longer functions, similar to bash -e flag.
func NewShell(acc *rs.Syscall) *Shell {
	return &Shell{account: acc}
}

type Shell struct {
	account *rs.Syscall
	err     error
}

// Execf
func (me *Shell) Execf(format string, args ...interface{}) error {
	if me.err != nil {
		return me.err
	}
	me.err = me.account.Execf(format, args...)
	return me.err
}
