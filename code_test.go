package passwordless_test

import (
	"testing"

	"github.com/alextanhongpin/passwordless"
)

func TestCodeConstructor(t *testing.T) {
	code := passwordless.NewCode()
	if err := code.Validate(); err != nil {
		t.Fatal(err)
	}

	if code.Code == "" {
		t.Fatal(".Code cannot be empty")
	}
}
