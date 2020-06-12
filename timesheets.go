package tidio

import (
	"encoding/json"
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

type Timesheets interface {
	WriteState(io.WriteCloser, error) error
	ReadState(io.ReadCloser, error) error
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
	return json.NewEncoder(w).Encode(m)
}

func (m *MemSheets) ReadState(r io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(m)
}
