package tidio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gregoryv/nugo"
	"github.com/preferit/tidio/internal"
)

func NewResource() *Resource {
	return &Resource{}
}

type Resource struct {
	nugo.Seal
	Path string
	io.ReadCloser
	Content string
}

func (s *Resource) Equal(b *Resource) bool {
	return s.Path == b.Path
}

// ----------------------------------------

type Resources interface {
	Stateful
	FilePersistent
	AddTimesheet(*Resource) error
	FindTimesheet(*Resource) error
	Map(SheetMapfunc) error
}

// ----------------------------------------

func NewMemResources() *MemResources {
	return &MemResources{
		Sheets:      make([]*Resource, 0),
		Source:      internal.None("MemSheets"),
		Destination: internal.None("MemSheets"),
	}
}

type MemResources struct {
	Source
	Destination
	Sheets []*Resource
}

func (m *MemResources) Load() error {
	r, err := m.Source.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(&m.Sheets)
}

func (m *MemResources) Save() error {
	w, err := m.Destination.Create()
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(&m.Sheets)
}

func (me *MemResources) PersistToFile(filename string) {
	me.Source = FileSource(filename)
	me.Destination = FileDestination(filename)
}

func (m *MemResources) AddTimesheet(s *Resource) error {
	m.Sheets = append(m.Sheets, s)
	return nil
}

func (m *MemResources) FindTimesheet(sheet *Resource) error {
	for _, s := range m.Sheets {
		if s.Equal(sheet) {
			*sheet = *s
			sheet.ReadCloser = ioutil.NopCloser(strings.NewReader(sheet.Content))
			return nil
		}
	}
	return fmt.Errorf("timesheet not found")
}

func (m *MemResources) Map(fn SheetMapfunc) error {
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

type SheetMapfunc func(*bool, *Resource) error
