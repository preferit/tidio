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
	OtherX
	OtherW
	OtherR

	GroupX
	GroupW
	GroupR

	UserX
	UserW
	UserR
)

// combinations
const (
	OtherRW  = (OtherR | OtherW)
	OtherRWX = (OtherR | OtherW | OtherX)
	OtherR_X = (OtherR | OtherX)

	GroupRW  = (GroupR | GroupW)
	GroupRWX = (GroupR | GroupW | GroupX)
	GroupR_X = (GroupR | GroupX)

	UserRW  = (UserR | UserW)
	UserRWX = (UserR | UserW | UserX)
	UserR_X = (UserR | UserX)
)

func (m PermMode) String() string {
	var s strings.Builder
	s.WriteByte(mChar(m, UserR, 'r'))
	s.WriteByte(mChar(m, UserW, 'w'))
	s.WriteByte(mChar(m, UserX, 'x'))
	s.WriteByte(mChar(m, GroupR, 'r'))
	s.WriteByte(mChar(m, GroupW, 'w'))
	s.WriteByte(mChar(m, GroupX, 'x'))
	s.WriteByte(mChar(m, OtherR, 'r'))
	s.WriteByte(mChar(m, OtherW, 'w'))
	s.WriteByte(mChar(m, OtherX, 'x'))
	return s.String()
}

func mChar(m, mask PermMode, c byte) byte {
	if m&mask == mask {
		return c
	}
	return '-'
}
