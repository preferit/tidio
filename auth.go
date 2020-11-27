package tidio

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/rs"
)

func NewBasicAuth(c *Credentials) *BasicAuth {
	return &BasicAuth{cred: c}
}

type BasicAuth struct {
	cred *Credentials
}

func (me *BasicAuth) Set(v interface{}) error {
	if me.cred == nil { // anonymous
		return nil
	}
	switch v := v.(type) {
	case *http.Request:
		plain := []byte(me.cred.account + ":" + me.cred.secret)
		b := base64.StdEncoding.EncodeToString(plain)
		v.Header.Set("Authorization", "Basic "+b)
		return nil
	default:
		return ant.SetFailed(v, me)
	}
}

func authenticate(sys *rs.System, r *http.Request, trace fox.Logger) (*rs.Account, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return rs.Anonymous, nil
	}
	trace.Log(h)

	name, secret, ok := r.BasicAuth()
	if !ok {
		return rs.Anonymous, fmt.Errorf("authentication failed")
	}

	asRoot := rs.Root.Use(sys)
	cmd := rs.NewCmd("/bin/secure", "-c", "-a", name, "-s", secret)
	trace.Log(cmd)
	if err := asRoot.Run(cmd); err != nil {
		return rs.Anonymous, err
	}
	var acc rs.Account
	err := asRoot.LoadAccount(&acc, name)
	return &acc, err
}
