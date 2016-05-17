package sql

import (
	"fmt"
	"testing"

	"github.com/RaduBerinde/playground/test_dep_prototype/server/testingshim"
)

func TestA(t *testing.T) {
	fmt.Printf("TestA\n")
	testingshim.NewTestServer().SQLSrv().(*SQLServer).Woof()
}
