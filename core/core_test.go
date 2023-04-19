package core_test

import (
	"testing"

	"github.com/templari/shire-client/core"
)

func TestCoreLogin(t *testing.T) {
	core := core.MakeCore("http://localhost:3011", nil)
	core.Login(1, "12346")
}
