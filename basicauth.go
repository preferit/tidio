package tidio

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// ParseBasicAuth parses a basic header
func ParseBasicAuth(h string) (*BasicAuth, error) {
	prefix := "Basic "
	if len(h) > 7 && h[:len(prefix)] != prefix {
		return nil, fmt.Errorf("ParseBasicAuth: prefix missmatch")
	}
	plain, err := base64.StdEncoding.DecodeString(h[len(prefix):])
	if err != nil {
		return nil, fmt.Errorf("ParseBasicAuth: %w", err)
	}
	parts := bytes.Split(plain, []byte(":"))
	if len(parts) != 2 {
		return nil, fmt.Errorf("ParseBasicAuth: invalid token")
	}
	return &BasicAuth{
		AccountName: string(parts[0]),
		Secret:      string(parts[1]),
	}, nil
}

type BasicAuth struct {
	AccountName string
	Secret      string
}

// Token returns a base64 encoded AccountName:Secret
func (me *BasicAuth) Token() string {
	return base64.StdEncoding.EncodeToString([]byte(me.String()))
}

// String returns a colon separated AccountName:Secret
func (me *BasicAuth) String() string {
	return fmt.Sprintf("%s:%s", me.AccountName, me.Secret)
}
