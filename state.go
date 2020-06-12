package tidio

import (
	"io"

	"github.com/gregoryv/box"
)

type Timesheet struct {
	Filename string
	Owner    string
	io.ReadCloser
	Content string
}

func (s *Timesheet) Equal(b *Timesheet) bool {
	return s.Filename == b.Filename
}

// ----------------------------------------

type MemSheets struct {
	Sheets []*Timesheet
}

func (m MemSheets) New() *MemSheets {
	e := &m
	e.Sheets = make([]*Timesheet, 0)
	return e
}

func (m *MemSheets) AddTimesheet(s *Timesheet) error {
	m.Sheets = append(m.Sheets, s)
	return nil
}

// here so we can synchronize all read/write operations
func (m *MemSheets) Save(store *box.Store, filename string) error {
	return store.SaveAs(m, filename)
}

func (m *MemSheets) Load(store *box.Store, filename string) error {
	return store.Load(m, filename)
}
