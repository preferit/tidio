package tidio

import "github.com/gregoryv/box"

type Data struct {
	Timesheets []*Timesheet
}

func (d *Data) Add(entity interface{}) error {
	switch entity := entity.(type) {
	case *Timesheet:
		d.Timesheets = append(d.Timesheets, entity)
	}
	return nil
}

// here so we can synchronize all read/write operations
func (d *Data) Save(store *box.Store, filename string) error {
	return store.SaveAs(d, filename)
}

func (d *Data) Load(store *box.Store, filename string) error {
	return store.Load(d, filename)
}
