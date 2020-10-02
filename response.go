package tidio

// Response defines a http response ready to be written.
type Response struct {
	view       interface{}
	statusCode int
}

// Fail sets the given status code on the response and returns the
// given error.
func (me *Response) Fail(code int, err error) error {
	me.statusCode = code
	return err
}

// End sets a 2xx code, optional view and returns nil. If code is not 2xx it panics.
//
func (me *Response) End(code int, view ...interface{}) error {
	if code < 200 || code > 299 {
		panic(code)
	}

	me.statusCode = code
	if len(view) > 0 {
		me.view = view[0]
	}
	return nil
}
