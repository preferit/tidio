package tidio

import (
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/rs"
)

func NewRouter(sys *rs.System) *Router {
	nop := fox.NewSyncLog(ioutil.Discard)

	return &Router{
		Logger: nop,
		warn:   nop.Log,
		sys:    sys,
	}
}

type Router struct {
	fox.Logger
	warn func(...interface{})

	sys *rs.System
}

func (me *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		sys: me.sys,
	}
	err := resp.Build(me.sys, r)
	if err != nil {
		resp.WriteError(w, err)
	}
	resp.Send(w)
}
