package tidio

import "github.com/gregoryv/box"

func NewState() *State {
	return &State{
		Timesheets: make([]*Timesheet, 0),
		Accounts:   make([]*Account, 0),
	}
}

type State struct {
	Timesheets []*Timesheet
	Accounts   []*Account
}

func (d *State) Add(entity interface{}) error {
	switch entity := entity.(type) {
	case *Timesheet:
		d.Timesheets = append(d.Timesheets, entity)
	}
	return nil
}

// here so we can synchronize all read/write operations
func (d *State) Save(store *box.Store, filename string) error {
	return store.SaveAs(d, filename)
}

func (d *State) Load(store *box.Store, filename string) error {
	return store.Load(d, filename)
}
