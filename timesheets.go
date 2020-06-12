package tidio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Timesheet struct {
	Filename string
	Owner    string
	io.ReadCloser
	Content string
}

func (t Timesheet) New() *Timesheet {
	return &t
}

func (s *Timesheet) Equal(b *Timesheet) bool {
	return s.Filename == b.Filename
}

// ----------------------------------------

type Timesheets interface {
	Stateful
	AddTimesheet(*Timesheet) error
	FindTimesheet(*Timesheet) error
	Map(SheetMapfunc) error
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

func (m *MemSheets) FindTimesheet(sheet *Timesheet) error {
	for _, s := range m.Sheets {
		if s.Equal(sheet) {
			*sheet = *s
			sheet.ReadCloser = ioutil.NopCloser(strings.NewReader(sheet.Content))
			return nil
		}
	}
	return fmt.Errorf("timesheet not found")
}

func (m *MemSheets) Map(fn SheetMapfunc) error {
	var next bool
	for _, s := range m.Sheets {
		next = true // by default we continue
		err := fn(&next, s)
		if !next || err != nil {
			return err
		}
	}
	return nil
}

type SheetMapfunc func(*bool, *Timesheet) error
