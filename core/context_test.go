package core

import (
	"testing"
)

func TestContext(t *testing.T) {
	_ = Context()
	Shutdown()
}
