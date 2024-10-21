package comshim

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	defer goleak.VerifyTestMain(m)
}
