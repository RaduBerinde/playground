package sql

import (
	"fmt"
	"testing"

	"github.com/RaduBerinde/playground/test_dep_prototype/server/testserver"
)

func TestA(t *testing.T) {
	fmt.Printf("TestA\n")
	testserver.TestSrvInstance.SqlSrv.(*SQLServer).Woof()
}
