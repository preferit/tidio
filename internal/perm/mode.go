// Package perm provides permission control mechanism similar to
// unix filesystem.

package perm

import "strings"

// PermMode represents a bit in a rwxrwxrwx permission
type PermMode int

// User: rwx Group: rwx Other: rwx
const (
	// note the reverse order
	NoMode PermMode = 1 << iota
	OtherExec
	OtherWrite
	OtherRead

	GroupExec
	GroupWrite
	GroupRead

	UserExec
	UserWrite
	UserRead
)

func (m PermMode) String() string {
	var s strings.Builder
	s.WriteByte(mChar(m, UserRead, 'r'))
	s.WriteByte(mChar(m, UserWrite, 'w'))
	s.WriteByte(mChar(m, UserExec, 'x'))
	s.WriteByte(mChar(m, GroupRead, 'r'))
	s.WriteByte(mChar(m, GroupWrite, 'w'))
	s.WriteByte(mChar(m, GroupExec, 'x'))
	s.WriteByte(mChar(m, OtherRead, 'r'))
	s.WriteByte(mChar(m, OtherWrite, 'w'))
	s.WriteByte(mChar(m, OtherExec, 'x'))
	return s.String()
}

func mChar(m, mask PermMode, c byte) byte {
	if m&mask == mask {
		return c
	}
	return '-'
}
