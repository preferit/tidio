package tidio

import (
	"encoding/base64"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestParseBasicAuth(t *testing.T) {
	raw := []byte("john:secret")
	token := base64.StdEncoding.EncodeToString(raw)
	ok, _ := asserter.NewMixed(t)
	ok(ParseBasicAuth("Basic " + token))
}

func TestParseBasicAuth_bad(t *testing.T) {
	_, bad := asserter.NewMixed(t)
	bad(ParseBasicAuth("Basic jibberish"))
	bad(ParseBasicAuth("Basic text"))
	bad(ParseBasicAuth("Bearer xxx"))
}
