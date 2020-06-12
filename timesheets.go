package tidio

import (
	"encoding/gob"
	"io"
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
func (m *MemSheets) WriteState(w io.WriteCloser, err error) error {
	if err != nil {
		return err
	}
	defer w.Close()
	return gob.NewEncoder(w).Encode(m)
}

func (m *MemSheets) ReadState(r io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer r.Close()
	return gob.NewDecoder(r).Decode(m)
}
