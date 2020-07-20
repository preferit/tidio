package tidio

import (
	"fmt"
	"strings"
)

func NewAuth() *Auth {
	return &Auth{}
}

type Auth struct {
	token string
}

// Parse value for bearer token
func (me *Auth) Parse(v string) error {
	prefix := "Bearer "
	if strings.Index(v, prefix) == -1 {
		return fmt.Errorf("Parse: no %s found", prefix)
	}
	me.token = v[len(prefix):]
	return nil
}
