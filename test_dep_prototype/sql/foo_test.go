package sql

import (
	"fmt"
	"testing"

	"github.com/RaduBerinde/playground/test_dep_prototype/server/testingshim"
)

func TestFoo(t *testing.T) {
	fmt.Printf("TestFoo\n")
	testingshim.NewTestServer().SQLSrv().(*SQLServer).Woof()
}
