package tidio

import (
	"encoding/base64"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestParseBasicAuth(t *testing.T) {
	raw := "john:secret"
	token := base64.StdEncoding.EncodeToString([]byte(raw))
	ok, bad := asserter.NewMixed(t)
	ok(ParseBasicAuth("Basic: " + token))
	bad(ParseBasicAuth("Basic: jibberish"))
	bad(ParseBasicAuth("Basic: text"))
	bad(ParseBasicAuth("Bearer: xxx"))
}
