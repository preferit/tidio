package tidio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/preferit/tidio/internal"
)

func NewTimesheet() *Timesheet {
	return &Timesheet{}
}

type Timesheet struct {
	Path string
	io.ReadCloser
	Content string
}

func (s *Timesheet) Equal(b *Timesheet) bool {
	return s.Path == b.Path
}

// ----------------------------------------

type Timesheets interface {
	Stateful
	FilePersistent
	AddTimesheet(*Timesheet) error
	FindTimesheet(*Timesheet) error
	Map(SheetMapfunc) error
}

// ----------------------------------------

func NewMemSheets() *MemSheets {
	return &MemSheets{
		Sheets:      make([]*Timesheet, 0),
		Source:      internal.None("MemSheets"),
		Destination: internal.None("MemSheets"),
	}
}

type MemSheets struct {
	Source
	Destination
	Sheets []*Timesheet
}

func (m *MemSheets) Load() error {
	r, err := m.Source.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(&m.Sheets)
}

func (m *MemSheets) Save() error {
	w, err := m.Destination.Create()
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(&m.Sheets)
}

func (me *MemSheets) PersistToFile(filename string) {
	me.Source = FileSource(filename)
	me.Destination = FileDestination(filename)
}

func (m *MemSheets) AddTimesheet(s *Timesheet) error {
	m.Sheets = append(m.Sheets, s)
	return nil
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
